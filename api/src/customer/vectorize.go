package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/datastore"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/webparse"
)

func (c *Customer) VectorizeDatastore(
	ctx context.Context,
	pool *pgxpool.Pool,
	job *queries.VectorizeJob,
) error {
	logger := c.logger.With("vectorizeJobId", job.ID.String())
	// keep track of the duration
	startTime := time.Now().UTC()

	// get the embeddings
	emb := llm.GetEmbeddings(logger, c.Customer)

	// track token usage throughout the program
	usageRecords := make([]*tokens.UsageRecord, 0)

	// create the model object
	dmodel := queries.New(pool)

	// set the job status
	logger.InfoContext(ctx, "Processing documents ...")
	if _, err := dmodel.UpdateVectorizeJobStatus(ctx, &queries.UpdateVectorizeJobStatusParams{
		ID:      job.ID,
		Status:  queries.VectorizeStatusInProgress,
		Message: "Processing documents",
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to update the job status", err)
	}

	// get all documents
	docs, err := dmodel.GetDocumentsByCustomer(ctx, c.ID)
	if err != nil {
		return fmt.Errorf("error getting all the documents: %s", err)
	}

	// process the documents
	for _, doc := range docs {
		// create a transaction
		tx, err := pool.Begin(ctx)
		if err != nil {
			return slogger.Error(ctx, logger, "failed to start a transaction", err)
		}

		errMsg := ""
		usageRecord, err := c.handleDocumentVectorization(ctx, tx, logger, emb, doc)
		if err == nil {
			if err := tx.Commit(ctx); err != nil {
				return slogger.Error(ctx, logger, "failed to commit the transaction", err)
			}
			if usageRecord != nil {
				usageRecords = append(usageRecords, usageRecord)
			}
		} else {
			errMsg = err.Error()
		}

		// update the state of the job item
		if _, err := dmodel.CreateVectorizeItem(ctx, &queries.CreateVectorizeItemParams{
			JobID:    job.ID,
			ObjectID: doc.ID,
			Error:    errMsg,
		}); err != nil {
			return slogger.Error(ctx, logger, "failed to create the vectorize item object", err)
		}
	}

	// track usage
	if err := utils.ReportUsage(ctx, logger, pool, c.ID, usageRecords, nil); err != nil {
		return slogger.Error(ctx, logger, "failed to report the usage", err)
	}

	logger.InfoContext(ctx, "Processing websites ...")
	if _, err := dmodel.UpdateVectorizeJobStatus(ctx, &queries.UpdateVectorizeJobStatusParams{
		ID:      job.ID,
		Status:  queries.VectorizeStatusInProgress,
		Message: "Processing websites ...",
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to update the job status", err)
	}

	// get the websites
	sites, err := dmodel.GetWebsitesByCustomer(ctx, c.ID)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to get the websites", err)
	}

	// parse all sites
	for _, site := range sites {
		response, err := c.handleWesbiteVectorization(ctx, pool, job, site)
		if err != nil {
			return slogger.Error(ctx, logger, "failed to process the site", err)
		}
		// add usage records
		usageRecords = append(usageRecords, response...)
	}

	// purge the datastore for all records 10 seconds older than the time this program took to run
	endTime := time.Now().UTC()
	diff := endTime.Sub(startTime) + 10*time.Second
	result := time.Now().UTC().Add(-diff).Format("2006-01-02 15:04:05")
	if err = c.PurgeDatastore(ctx, pool, &purgeDatastoreRequest{
		Timestamp: &result,
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to purge the datastore", err)
	}

	logger.InfoContext(ctx, "Reporting usage ...")
	if _, err := dmodel.UpdateVectorizeJobStatus(ctx, &queries.UpdateVectorizeJobStatusParams{
		ID:      job.ID,
		Status:  queries.VectorizeStatusInProgress,
		Message: "Reporting usage ...",
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to update the job status", err)
	}

	if err := utils.ReportUsage(ctx, logger, pool, c.ID, usageRecords, nil); err != nil {
		return slogger.Error(ctx, logger, "failed to report the usage", err)
	}

	logger.InfoContext(ctx, "Successfully vectorized customer store")
	return nil
}

func (c *Customer) handleDocumentVectorization(
	ctx context.Context,
	db queries.DBTX,
	l *slog.Logger,
	emb gollm.Embeddings,
	item *queries.Document,
) (*tokens.UsageRecord, error) {
	logger := l.With("docID", item.ID, "filename", item.Filename)
	dmodel := queries.New(db)

	doc, err := datastore.NewDocumentFromDocument(ctx, l, item)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to parse the database doc", err)
	}

	// get the cleaned data from the document and a parser
	logger.InfoContext(ctx, "Fetching document from datastore ...")
	cleaned, err := doc.GetCleaned(ctx)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "there was an issue getting the cleaned document", err)
	}

	// TODO -- there is potential to get this to work
	// see if the document changed
	newSha256, err := doc.GetSha256()
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the sha256 of the document", err)
	}
	if doc.VectorSha256 == newSha256 {
		l.InfoContext(ctx, "This document has not changed")
		if err := dmodel.UpdateDocumentVectorSig(ctx, &queries.UpdateDocumentVectorSigParams{
			ID:           doc.ID,
			VectorSha256: newSha256,
		}); err != nil {
			return nil, slogger.Error(ctx, logger, "failed to touch document", err)
		}
		return nil, nil
	} else {
		l.InfoContext(ctx, "The signatures do not match", "old", doc.Sha256, "newSha256", newSha256)
	}

	// delete the old vectors
	if err := dmodel.DeleteDocumentVectors(ctx, doc.ID); err != nil {
		return nil, slogger.Error(ctx, logger, "failed to delete old vectors", err)
	}

	// embed the content
	logger.InfoContext(ctx, "Embedding the document ...")
	res, err := emb.Embed(ctx, cleaned.String())
	if err != nil {
		return nil, slogger.Error(ctx, logger, "error embedding the content", err)
	}

	// insert the vectors into the database
	logger.InfoContext(ctx, "Inserting all documents into the database")
	for index, vec := range res.Embeddings {
		logger.DebugContext(ctx, "Processing index", "index", index)
		// create raw vector object
		vecId, err := dmodel.CreateVector(ctx, &queries.CreateVectorParams{
			Raw:        vec.Raw,
			Embeddings: &vec.Embedding,
			CustomerID: c.ID,
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to insert the vector object", err)
		}

		// create a reference to the vector onto the document
		_, err = dmodel.CreateDocumentVector(ctx, &queries.CreateDocumentVectorParams{
			DocumentID:    doc.ID,
			VectorStoreID: vecId,
			CustomerID:    c.ID,
			Index:         int32(index),
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed creating document vector relationship", err)
		}
		logger.InfoContext(ctx, "Finished.")
	}

	// set updated and sha field for intelligent parsing in the future
	if err := dmodel.UpdateDocumentVectorSig(ctx, &queries.UpdateDocumentVectorSigParams{
		ID:           doc.ID,
		VectorSha256: newSha256,
	}); err != nil {
		return nil, slogger.Error(ctx, logger, "failed to touch document", err)
	}

	logger.InfoContext(ctx, "Successfully processed document")
	return res.Usage, nil
}

func (c *Customer) handleWesbiteVectorization(
	ctx context.Context,
	pool *pgxpool.Pool,
	job *queries.VectorizeJob,
	site *queries.Website,
) ([]*tokens.UsageRecord, error) {
	logger := c.logger.With("site.ID", site.ID.String(), "site.Domain", site.Domain)
	logger.InfoContext(ctx, "Parsing site ...")

	// track token usage throughout the program
	usageRecords := make([]*tokens.UsageRecord, 0)

	// get the pages
	dmodel := queries.New(pool)
	pages, err := dmodel.GetWebsitePagesBySite(ctx, site.ID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue fetching the sites: %v", err)
	}

	logger.InfoContext(ctx, "Creating embeddings for each page ...")

	// create an embeddings object
	emb := llm.GetEmbeddings(logger, c.Customer)

	for _, page := range pages {
		// create a transaction
		tx, err := pool.Begin(ctx)
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to start a transaction", err)
		}

		errMsg := ""
		usageRecord, err := c.handleWebsitePageVectorization(ctx, tx, logger, emb, page)
		if err == nil {
			// commit the transction
			if err := tx.Commit(ctx); err != nil {
				return nil, slogger.Error(ctx, logger, "failed to commit the transaction", err)
			}
			if usageRecord != nil {
				usageRecords = append(usageRecords, usageRecord)
			}
		} else {
			errMsg = err.Error()
		}

		// update the state of the job item
		if _, err := dmodel.CreateVectorizeItem(ctx, &queries.CreateVectorizeItemParams{
			JobID:    job.ID,
			ObjectID: page.ID,
			Error:    errMsg,
		}); err != nil {
			return nil, slogger.Error(ctx, logger, "failed to create the vectorize item object", err)
		}
	}

	// report usage
	if err := utils.ReportUsage(ctx, logger, pool, c.ID, usageRecords, nil); err != nil {
		return nil, slogger.Error(ctx, logger, "failed to report usage", err)
	}

	logger.InfoContext(ctx, "Processed all pages")
	return usageRecords, nil
}

func (c *Customer) handleWebsitePageVectorization(
	ctx context.Context,
	db queries.DBTX,
	l *slog.Logger,
	emb gollm.Embeddings,
	page *queries.WebsitePage,
) (*tokens.UsageRecord, error) {
	logger := l.With("page", page.Url)
	dmodel := queries.New(db)

	logger.InfoContext(ctx, "Scraping the page ...")

	// scrape the webpage
	scrapeResponse, err := webparse.ScrapeSingle(ctx, logger, page)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to scrape the page", err)
	}

	// create a signature to compare the old vs new
	newSha256 := utils.GenerateFingerprint([]byte(scrapeResponse.Content))
	if page.Sha256 == newSha256 {
		logger.InfoContext(ctx, "this page did not change")
		// update the page header
		if _, err := dmodel.UpdateWebsitePageSignature(ctx, &queries.UpdateWebsitePageSignatureParams{
			ID:     page.ID,
			Sha256: newSha256,
		}); err != nil {
			return nil, slogger.Error(ctx, logger, "failed to update the page signature", err)
		}
		return nil, nil
	} else {
		logger.InfoContext(ctx, "The signatures do not match", "oldSHA256", page.Sha256, "newSHA256", newSha256)
	}

	logger.InfoContext(ctx, "Vecorizing the content ...")

	// update the page header
	if _, err := dmodel.UpdateWebsitePageSignature(ctx, &queries.UpdateWebsitePageSignatureParams{
		ID:     page.ID,
		Sha256: newSha256,
	}); err != nil {
		return nil, slogger.Error(ctx, logger, "failed to update the page signature", err)
	}

	// embed the content
	res, err := emb.Embed(ctx, scrapeResponse.Content)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to embed the content", err)
	}

	encodedHeaders, err := json.Marshal(scrapeResponse.Header)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to encode the page headers", err)
	}

	// lastly upload the vectors to the datastore
	for index, vec := range res.Embeddings {
		// create raw vector object
		vecId, err := dmodel.CreateVector(ctx, &queries.CreateVectorParams{
			Raw:        vec.Raw,
			Embeddings: &vec.Embedding,
			CustomerID: c.ID,
			Metadata:   encodedHeaders,
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to insert the embeddings", err)
		}

		// create a reference to the vector onto the document
		_, err = dmodel.CreateWebsitePageVector(ctx, &queries.CreateWebsitePageVectorParams{
			WebsitePageID: page.ID,
			VectorStoreID: vecId,
			CustomerID:    c.ID,
			Index:         int32(index),
			Metadata:      encodedHeaders,
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to create the vector relationship", err)
		}
	}

	logger.InfoContext(ctx, "Successfully processed page")
	return res.Usage, nil
}
