package customer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

type RAGRequest struct {
	// general params
	ConversationId string
	Input          string
	CheckQuality   bool

	// models
	SummaryModelId string
	ChatModelId    string

	// ids for scoped content
	FolderIds      []string
	DocumentIds    []string
	WebsiteIds     []string
	WebsitePageIds []string
}

type RAGResponse struct {
	ConverationId string                 `json:"converationId"`
	Documents     []*queries.Document    `json:"documents"`
	WebsitePages  []*queries.WebsitePage `json:"websitePages"`
	Response      string                 `json:"response"`
}

func (c *Customer) RAG(
	ctx context.Context,
	db queries.DBTX,
	args *RAGRequest,
) (*RAGResponse, error) {
	logger := c.logger.With("function", "RAG")
	logger.InfoContext(ctx, "Beginning document retrieval pathway")
	dmodel := queries.New(db)

	// get the required structures
	logger.DebugContext(ctx, "Getting required objects ...")
	docstore, err := c.GetDocstore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the docstore: %s", err)
	}
	embs := c.GetEmbeddings(ctx)

	// track all token usage across this request
	usageRecords := make(chan *tokens.TokenRecord, 100)

	/// get the conversation
	logger.InfoContext(ctx, "Getting conversation ...")
	conv, err := llm.AutoConversation(ctx, logger, db, c.ID, args.ConversationId, prompts.RAG_COMPLETE_SYSTEM_PROMPT, "Information Chat", "rag")
	if err != nil {
		return nil, fmt.Errorf("failed to parse the conversation: %s", err)
	}

	/// get the chat llm
	logger.InfoContext(ctx, "Getting the chat llm ...")
	var chatLLMId pgtype.UUID
	chatLLMId.Scan(args.ChatModelId)
	chatLLM, err := llm.GetLLM(ctx, db, c.ID, chatLLMId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the chat llm: %s", err)
	}

	/// generate a simpler query based on the users request

	// get the llm
	vectorQueryModel, err := dmodel.GetInteralLLM(ctx, "Vector Query Generator")
	if err != nil {
		return nil, fmt.Errorf("failed to get internal llm: %s", err)
	}

	// create the completion based on the user's query
	logger.InfoContext(ctx, "Simplifying the user's query ...")
	vectorQueryLLM := &llm.LLM{Llm: vectorQueryModel}
	simpleQueryResponse, err := llm.SingleCompletion(
		ctx, vectorQueryLLM, logger, c.ID,
		prompts.RAG_SIMPLE_QUERY_SYSTEM_PROMPT,
		usageRecords,
		&llm.CompletionArgs{Input: args.Input},
	)
	if err != nil || simpleQueryResponse == "" {
		return nil, fmt.Errorf("error performing the simplifier completion: %s", err)
	}

	// parse the query list
	simpleQueries := strings.Split(simpleQueryResponse, ",")

	/// query the vector store (potentially multiple times based on the returned query list)
	/// and collect into a list of documents and website pages

	// create structures for routines
	var wg sync.WaitGroup
	k := 2 // number of objects to query
	docs := make(chan *vectorstore.DocumentResponse, k*len(simpleQueries))
	pages := make(chan *vectorstore.WebsitePageResponse, k*len(simpleQueries))
	errCh := make(chan error, 2*len(simpleQueries))

	queryInput := vectorstore.QueryInput{
		CustomerId:     c.ID,
		Docstore:       docstore,
		Embeddings:     embs,
		DB:             db,
		K:              2,
		Logger:         logger,
		IncludeContent: true,
	}

	// parse the args for the associated ids
	websiteIds := make([]uuid.UUID, len(args.WebsiteIds))
	for _, item := range args.WebsiteIds {
		parsed, err := uuid.Parse(item)
		if err == nil {
			websiteIds = append(websiteIds, parsed)
		}
	}
	websitePageIds := make([]uuid.UUID, len(args.WebsitePageIds))
	for _, item := range args.WebsitePageIds {
		parsed, err := uuid.Parse(item)
		if err == nil {
			websitePageIds = append(websitePageIds, parsed)
		}
	}
	folderIds := make([]uuid.UUID, len(args.FolderIds))
	for _, item := range args.FolderIds {
		parsed, err := uuid.Parse(item)
		if err == nil {
			folderIds = append(folderIds, parsed)
		}
	}
	documentIds := make([]uuid.UUID, len(args.DocumentIds))
	for _, item := range args.DocumentIds {
		parsed, err := uuid.Parse(item)
		if err == nil {
			documentIds = append(documentIds, parsed)
		}
	}

	logger.InfoContext(ctx, "Querying the users information for each simple query ...")
	for _, item := range simpleQueries {
		l := logger.With("query", item)
		l.InfoContext(ctx, "Running query ...")

		// get the docs
		wg.Add(1)
		go func(q string, input vectorstore.QueryInput) {
			defer wg.Done()
			l.InfoContext(ctx, "Querying documents ...")

			input.Query = q
			response, err := vectorstore.QueryDocuments(ctx, &vectorstore.QueryDocstoreInput{
				QueryInput:  &input,
				FolderIds:   folderIds,
				DocumentIds: documentIds,
			})
			if err != nil {
				l.ErrorContext(ctx, "failed to query the documents", "error", err)
				errCh <- err
				return
			}

			for _, doc := range response {
				docs <- doc
			}
		}(item, queryInput)

		// get the pages
		wg.Add(1)
		go func(q string, input vectorstore.QueryInput) {
			defer wg.Done()

			l.InfoContext(ctx, "Querying website pages ...")

			input.Query = q
			response, err := vectorstore.QueryWebsitePages(ctx, &vectorstore.QueryWebsitePagesInput{
				QueryInput:     &input,
				WebsiteIds:     websiteIds,
				WebsitePageIds: websitePageIds,
			})
			if err != nil {
				l.ErrorContext(ctx, "failed to query the website pages", "error", err)
				errCh <- err
			}

			for _, page := range response {
				pages <- page
			}
		}(item, queryInput)
	}

	// collect the items
	logger.DebugContext(ctx, "Waiting for go-routines to finish ...")
	wg.Wait()
	logger.DebugContext(ctx, "Collecting channels ...")
	close(docs)
	close(pages)
	close(errCh)

	// check for errors
	logger.DebugContext(ctx, "Checking for errors ...")
	for err := range errCh {
		return nil, err
	}

	/// spawn go-routines for each document and website page

	// get the summary model
	var summaryModelId pgtype.UUID
	summaryModelId.Scan(args.SummaryModelId)
	summaryLLM, err := llm.GetLLM(ctx, db, c.ID, summaryModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to get a suitable model for summarizaion: %s", err)
	}

	// get the ranker llm
	rankerModel, err := dmodel.GetInteralLLM(ctx, "Content Ranker")
	if err != nil {
		return nil, fmt.Errorf("failed to get internal llm: %s", err)
	}
	rankerLLM := &llm.LLM{Llm: rankerModel}

	// channel for critical errors
	errCh = make(chan error)

	// objects to update
	updateDocs := make(chan *vectorstore.DocumentResponse, k*len(simpleQueries))
	updatePages := make(chan *vectorstore.WebsitePageResponse, k*len(simpleQueries))

	// objects that passed validation
	validatedDocs := make(chan *vectorstore.DocumentResponse, k*len(simpleQueries))
	validatedPages := make(chan *vectorstore.WebsitePageResponse, k*len(simpleQueries))

	// parse all documents
	for i := range docs {
		l := logger.With("doc", i.ID.String())
		wg.Add(1)
		go func(doc *vectorstore.DocumentResponse) {
			defer wg.Done()
			l.InfoContext(ctx, "Parsing document ...")

			/// generate a summary of the content if needed
			genSummary := false
			if doc.Summary == "" {
				l.InfoContext(ctx, "Document has no summary")
				genSummary = true
			} else {
				l.InfoContext(ctx, "Document already has a summary")
				if doc.SummarySha256 != doc.Sha256 {
					l.InfoContext(ctx, "Summary sha256 does not match")
					genSummary = true
				}
			}

			// generate the summary of the document
			if genSummary {
				// create a completion for this document
				summary, err := llm.SingleCompletion(
					ctx, summaryLLM, l, c.ID,
					prompts.RAG_SIMPLE_QUERY_SYSTEM_PROMPT,
					usageRecords,
					&llm.CompletionArgs{Input: doc.Content},
				)
				if err != nil {
					// do not fail for summaries like this
					l.WarnContext(ctx, "Failed to create the summary", "error", err)
					return
				} else {
					// add summary and update
					doc.Summary = summary
					updateDocs <- doc
				}
			}

			/// validate the summary for relevance
			if args.CheckQuality {
				l.InfoContext(ctx, "Ranking the summary ...")
				rankerResponse, err := llm.SingleCompletionJson[prompts.RagRankerSchema](
					ctx, rankerLLM, l, c.ID,
					prompts.RAG_RANKER_SYSTEM_PROMPT,
					usageRecords,
					&llm.CompletionArgs{
						Input:      doc.Summary,
						Json:       true,
						JsonSchema: prompts.RAG_RANKER_SCHEMA,
					},
				)
				if err != nil {
					errCh <- fmt.Errorf("error performing the ranking completion: %s", err)
				}

				l.DebugContext(ctx, "Successfully ranked summary", "relevance", rankerResponse.Relevance, "quality", rankerResponse.Quality)

				// evaluate the quality
				if rankerResponse.Relevance > 40 && rankerResponse.Quality > 70 {
					l.InfoContext(ctx, "Document passed performance evaluation")
					validatedDocs <- doc
				}
			} else {
				l.InfoContext(ctx, "Skipping quality check")
				validatedDocs <- doc
			}
		}(i)
	}

	// parse all pages
	for i := range pages {
		l := logger.With("page", i.ID.String())
		wg.Add(1)
		go func(page *vectorstore.WebsitePageResponse) {
			defer wg.Done()
			l.InfoContext(ctx, "Parsing website page ...")

			/// generate a summary of the content if needed
			genSummary := false
			if page.Summary == "" {
				l.InfoContext(ctx, "Page has no summary")
				genSummary = true
			} else if page.SummarySha256 != page.Sha256 {
				l.InfoContext(ctx, "Summary sha256 does not match page sha256")
				genSummary = true
			} else {
				l.InfoContext(ctx, "Using cached summary")
				genSummary = false
			}

			// generate the summary of the document
			if genSummary {
				// create a completion for this document
				summary, err := llm.SingleCompletion(
					ctx, summaryLLM, l, c.ID,
					prompts.RAG_SIMPLE_QUERY_SYSTEM_PROMPT,
					usageRecords,
					&llm.CompletionArgs{Input: page.Content},
				)
				if err != nil {
					// do not fail for summaries like this
					l.WarnContext(ctx, "Failed to create the summary for the page", "error", err)
					return
				} else {
					// add summary and update
					page.Summary = summary
					updatePages <- page
				}
			}

			if args.CheckQuality {
				/// validate the summary for relevance
				l.InfoContext(ctx, "Ranking the page summary ...")
				rankerResponse, err := llm.SingleCompletionJson[prompts.RagRankerSchema](
					ctx, rankerLLM, l, c.ID,
					prompts.RAG_RANKER_SYSTEM_PROMPT,
					usageRecords,
					&llm.CompletionArgs{
						Input:      page.Summary,
						Json:       true,
						JsonSchema: prompts.RAG_RANKER_SCHEMA,
					},
				)
				if err != nil {
					errCh <- fmt.Errorf("error performing the ranking completion on the page: %s", err)
				}

				l.DebugContext(ctx, "Successfully ranked summary", "relevance", rankerResponse.Relevance, "quality", rankerResponse.Quality)

				// evaluate the quality
				if rankerResponse.Relevance > 40 && rankerResponse.Quality > 70 {
					l.InfoContext(ctx, "Page passed performance evaluation")
					validatedPages <- page
				}
			} else {
				l.InfoContext(ctx, "Skipping quality check")
				validatedPages <- page
			}
		}(i)
	}

	/// TODO -- parse the web for results to aid this portion of the conversation

	/// collect all go routines and channels
	wg.Wait()
	close(updateDocs)
	close(updatePages)
	close(validatedDocs)
	close(validatedPages)
	close(usageRecords)
	close(errCh)

	// check for critical errors
	for err := range errCh {
		return nil, err
	}

	// create lists of the validated objects
	validatedDocsList := make([]*vectorstore.DocumentResponse, 0)
	for doc := range validatedDocs {
		validatedDocsList = append(validatedDocsList, doc)
	}
	validatedPagesList := make([]*vectorstore.WebsitePageResponse, 0)
	for page := range validatedPages {
		validatedPagesList = append(validatedPagesList, page)
	}

	/// compose the query for the correct model
	documentSummaries := ""
	for _, doc := range validatedDocsList {
		documentSummaries += fmt.Sprintf("\nDocument: %s", doc.Summary)
	}
	pageSummaries := ""
	for _, page := range validatedPagesList {
		pageSummaries += fmt.Sprintf("\nInternal Page: %s", page.Summary)
	}
	query := fmt.Sprintf("User Query: %s\nDocuments: %s\nInternal Pages: %s", args.Input, documentSummaries, pageSummaries)

	/// send the request and the

	/// send the conversation request and the reporting/updating functions concurrently to save time
	var convResponse string
	errCh = make(chan error)

	// send the conversation request
	wg.Add(1)
	go func() {
		defer wg.Done()
		r, err := conv.Completion(ctx, db, chatLLM, query)
		if err != nil {
			errCh <- fmt.Errorf("failed to complete the conversation: %s", err)
			return
		}
		convResponse = r

		/// analyze the response for hallucinations
		logger.InfoContext(ctx, "TODO -- analyze for halucinations")
	}()

	// run code to finish after the request has been sent
	wg.Add(1)
	go func() {
		defer wg.Done()
		// report all token usage
		logger.InfoContext(ctx, "Reporting usage ...")
		totalRecords := make([]*tokens.TokenRecord, 0)
		for item := range usageRecords {
			totalRecords = append(totalRecords, item)
		}
		if err := utils.ReportUsage(ctx, logger, db, c.ID, totalRecords, nil); err != nil {
			errCh <- fmt.Errorf("failed to report usage: %s", err)
		}

		// update all documents
		logger.InfoContext(ctx, "Updating documents ...")
		for doc := range updateDocs {
			_, err := dmodel.UpdateDocumentSummary(ctx, &queries.UpdateDocumentSummaryParams{
				ID:            doc.ID,
				Summary:       doc.Summary,
				SummarySha256: doc.Sha256,
			})
			if err != nil {
				logger.ErrorContext(ctx, "failed to update document", "error", err, "doc.ID", doc.ID)
			}
		}
		logger.InfoContext(ctx, "Successfully updated documents")

		logger.InfoContext(ctx, "Updating website pages ...")
		for page := range updatePages {
			_, err := dmodel.UpdateWebsitePageSummary(ctx, &queries.UpdateWebsitePageSummaryParams{
				ID:            page.ID,
				Summary:       page.Summary,
				SummarySha256: page.Sha256,
			})
			if err != nil {
				logger.ErrorContext(ctx, "failed to update website page", "error", err, "page.ID", page.ID)
			}
		}
		logger.InfoContext(ctx, "Successfully updated website pages")
	}()

	wg.Wait()
	close(errCh)

	// check for errors
	for err := range errCh {
		return nil, err
	}

	// create the return lists
	returnDocs := make([]*queries.Document, len(validatedDocsList))
	returnPages := make([]*queries.WebsitePage, len(validatedPagesList))
	for _, item := range validatedDocsList {
		returnDocs = append(returnDocs, item.Document.Document)
	}
	for _, item := range validatedPagesList {
		returnPages = append(returnPages, item.WebsitePage)
	}

	/// return to the user
	return &RAGResponse{
		ConverationId: conv.ID.String(),
		Documents:     returnDocs,
		WebsitePages:  returnPages,
		Response:      convResponse,
	}, nil
}
