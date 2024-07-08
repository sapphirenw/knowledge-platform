package customer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/datastore"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/webparse"
)

// Wrapper around the `queries.Customer` object that represents the database object
// in order to store some state about the customer when needed
type Customer struct {
	*queries.Customer

	logger *slog.Logger
}

func CreateCustomer(ctx context.Context, logger *slog.Logger, db queries.DBTX, request *createCustomerRequest) (*Customer, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	l := logger.With("request", *request)
	l.InfoContext(ctx, "Creating a new customer ...")

	model := queries.New(db)
	customer, err := model.CreateCustomer(ctx, request.Name)
	if err != nil {
		return nil, fmt.Errorf("error creating the customer")
	}

	l.InfoContext(ctx, "Successfully created new customer", "customer", *customer)

	return &Customer{
		Customer: customer,
		logger:   logger,
	}, nil
}

// grabs a customer from the database using the supplied id
func NewCustomer(ctx context.Context, logger *slog.Logger, id uuid.UUID, db queries.DBTX) (*Customer, error) {
	model := queries.New(db)

	logger.InfoContext(ctx, "Fetching the customer record")

	// get the customer
	c, err := model.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Customer{
		Customer: c,
		logger:   logger.With("customer.ID", c.ID.String(), "customer.Name", c.Name, "customer.Datastore", c.Datastore),
	}, nil
}

// Gets the docstore associated with the customer
func (c *Customer) GetDocstore(ctx context.Context) (docstore.RemoteDocstore, error) {
	switch c.Customer.Datastore {
	case "s3":
		return docstore.NewS3Docstore(ctx, docstore.S3_BUCKET, c.logger)
	default:
		return docstore.NewTODODocstore(c.logger)
	}
}

// Creates a new folder tied to the customer with an optional parent.
func (c *Customer) CreateFolder(ctx context.Context, txn queries.DBTX, args *createFolderRequest) (*queries.Folder, error) {
	if args == nil {
		return nil, fmt.Errorf("the arguments cannot be nil")
	}
	if args.Name == "" {
		return nil, fmt.Errorf("the name cannot be empty")
	}

	logger := c.logger.With("folder.Name", args.Name)

	// parse the owner if applicable
	var parentId pgtype.UUID
	if args.Owner != nil {
		logger = c.logger.With("folder.Owner", args.Owner)
		parentId.Scan(args.Owner)
	}
	logger.InfoContext(ctx, "Creating a new folder ...")

	// create the folder
	model := queries.New(txn)
	folder, err := model.CreateFolder(ctx, &queries.CreateFolderParams{
		ParentID:   parentId,
		CustomerID: c.ID,
		Title:      args.Name,
	})
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			logger.InfoContext(ctx, "The folder already exists, returning to the customer", "id", folder.ID)
			// the folder already exists in this location, return the folder
			return folder, nil
		}
		return nil, fmt.Errorf("there was an issue creating the folder: %v", err)
	}
	logger.InfoContext(ctx, "Successfully created folder", "id", folder.ID)

	return folder, err
}

// Does an 'ls' on a folder
func (c *Customer) ListFolderContents(ctx context.Context, db queries.DBTX, folderId pgtype.UUID) (*listFolderContentsResponse, error) {
	logger := c.logger.With()
	if folderId.Valid {
		logger = logger.With("folderId", folderId.Bytes)
	}
	logger.InfoContext(ctx, "Getting all children of the folder ...")

	model := queries.New(db)

	var err error
	var folder *queries.Folder
	var folders []*queries.Folder
	var documents []*queries.Document

	if !folderId.Valid {
		// get root file information
		folders, err = model.GetRootFoldersByCustomer(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("there was an issue getting the folders: %v", err)
		}
		documents, err = model.GetRootDocumentsByCustomer(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("there was an issue getting the documents: %v", err)
		}
	} else {
		// get self
		uid := utils.PGXUUIDToGoogleUUID(folderId)
		if uid == nil {
			return nil, fmt.Errorf("failed to convert folderId to a google uuid")
		}
		folder, err = model.GetFolder(ctx, *uid)
		if err != nil {
			return nil, fmt.Errorf("this folder does not exist: %s", err)
		}

		// get the information using the folder
		folders, err = model.GetFoldersFromParent(ctx, folderId)
		if err != nil {
			return nil, fmt.Errorf("there was an issue getting the folders: %v", err)
		}
		documents, err = model.GetDocumentsFromParent(ctx, folderId)
		if err != nil {
			return nil, fmt.Errorf("there was an issue getting the documents: %v", err)
		}
	}

	logger.InfoContext(ctx, "Successfully listed folder contents", "folders", len(folders), "documents", len(documents))

	return &listFolderContentsResponse{
		Self:      folder,
		Folders:   folders,
		Documents: documents,
	}, nil
}

/*
Generates pre-signed urls for the user to use to upload to their preferred datastore.
*/
func (c *Customer) GeneratePresignedUrl(ctx context.Context, db queries.DBTX, body *generatePresignedUrlRequest) (*generatePresignedUrlResponse, error) {
	logger := c.logger.With("body", *body)
	logger.InfoContext(ctx, "Generating a presigned url...")

	// get the customer's docstore
	store, err := c.GetDocstore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the document store: %v", err)
	}

	var parentId pgtype.UUID
	if body.ParentId != nil {
		parentId.Scan(*body.ParentId)
	}

	// insert a record into the documents table
	model := queries.New(db)
	d, err := model.CreateDocument(ctx, &queries.CreateDocumentParams{
		ParentID:      parentId,
		CustomerID:    c.ID,
		Filename:      body.Filename,
		Type:          body.Mime,
		SizeBytes:     body.Size,
		Sha256:        body.Signature,
		DatastoreID:   fmt.Sprintf("%s/%s", c.ID.String(), uuid.New().String()),
		DatastoreType: c.Datastore,
	})
	if err != nil {
		return nil, fmt.Errorf("there was an issue creating the document: %v", err)
	}

	// create the document
	doc, err := datastore.NewDocumentFromDocument(ctx, logger, d)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the database document: %s", err)
	}

	// generate the pre-signed url
	url, err := store.GeneratePresignedUrl(ctx, doc.Document, doc.Type, doc.DatastoreID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue generating the presigned url: %v", err)
	}

	logger.InfoContext(ctx, "Successfully generated the pre-signed url")

	return &generatePresignedUrlResponse{
		UploadUrl:  base64.StdEncoding.EncodeToString([]byte(url)),
		Method:     store.GetUploadMethod(),
		DocumentId: d.ID,
	}, nil
}

/*
Function to notify the server that the document upload using the pre-signed url was successful, and the
server can store the record of this object in the datastore.
*/
func (c *Customer) NotifyOfSuccessfulUpload(ctx context.Context, db queries.DBTX, documentId uuid.UUID) error {
	logger := c.logger.With("documentId", documentId)
	logger.InfoContext(ctx, "Marking the document as successfully uploaded")

	// create the database object
	model := queries.New(db)
	_, err := model.MarkDocumentAsUploaded(ctx, documentId)
	if err != nil {
		// TODO -- implement a critical error here that can contain information to be notified by
		return fmt.Errorf("there was an issue inserting the document into the database: %v", err)
	}

	logger.InfoContext(ctx, "Successfully validated document")
	return nil
}

/*
Adds a website for the user, but does not scrape it.
*/
func (c *Customer) HandleWebsite(ctx context.Context, db queries.DBTX, request *handleWebsiteRequest) (*handleWebsiteResponse, error) {
	logger := c.logger.With("domain", request.Domain)
	logger.InfoContext(ctx, "Ingesting the domain...", "whitelist", request.Whitelist, "blacklist", request.Blacklist, "insert", request.Insert)

	// ensure the domain has a scheme
	parsed, err := url.Parse(request.Domain)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to parse the domain", err)
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}

	// compose the new domain
	var parsedDomain string
	if parsed.Host != "" {
		parsedDomain = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	} else if parsed.Path != "" {
		parsedDomain = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Path)
	} else {
		return nil, slogger.Error(ctx, logger, "failed to parse the domain", err, "parsed", parsed)
	}

	logger = logger.With("domain", parsedDomain)
	logger.Info("Cleaned the domain name")

	// parse the domain
	protocol, domain, err := utils.ParseWebsiteInformation(parsedDomain)
	if err != nil {
		return nil, fmt.Errorf("error parsing the website: %v", err)
	}

	// create a site object
	tmpSite := queries.Website{
		CustomerID: c.ID,
		Protocol:   protocol,
		Domain:     domain,
		Blacklist:  request.Blacklist,
		Whitelist:  request.Whitelist,
	}

	// parse the pages from the site
	urls, err := webparse.ParseSitemap(ctx, logger, &tmpSite, 100)
	if err != nil {
		return nil, fmt.Errorf("there was an issue parsing the sitemap: %v", err)
	}

	pages := make([]*queries.WebsitePage, len(urls))

	// send back the parsed data if not an insertion request
	if !request.Insert {
		// create tmp pages
		for i, item := range urls {
			pages[i] = &queries.WebsitePage{
				CustomerID: c.ID,
				Url:        item,
			}
		}

		return &handleWebsiteResponse{
			Site:  &tmpSite,
			Pages: pages,
		}, nil
	}

	// insert the website
	model := queries.New(db)
	site, err := model.CreateWebsite(ctx, &queries.CreateWebsiteParams{
		CustomerID: c.ID,
		Protocol:   protocol,
		Domain:     domain,
		Blacklist:  request.Blacklist,
		Whitelist:  request.Whitelist,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating the website: %v", err)
	}

	// set all existing pages to  not valid, the create call will re-set
	// to valid if the page already exists, and default is TRUE
	if err := model.SetWebsitePagesNotValid(ctx, &queries.SetWebsitePagesNotValidParams{
		CustomerID: c.ID,
		WebsiteID:  site.ID,
	}); err != nil {
		return nil, fmt.Errorf("err setting website pages to not valid: %s", err)
	}

	// insert the pages
	for i, item := range urls {
		page, err := model.CreateWebsitePage(ctx, &queries.CreateWebsitePageParams{
			CustomerID: c.ID,
			WebsiteID:  site.ID,
			Url:        item,
			Sha256:     utils.GenerateFingerprint([]byte(item)), // use a tmp hash until the content is actually ingested
		})
		if err != nil {
			return nil, fmt.Errorf("there was an issue inserting the page: %v", err)
		}
		pages[i] = page
	}

	// delete the records that are not valid, these are the stale records
	if err := model.DeleteWebsitePagesNotValid(ctx, &queries.DeleteWebsitePagesNotValidParams{
		CustomerID: c.ID,
		WebsiteID:  site.ID,
	}); err != nil {
		return nil, fmt.Errorf("error deleting stale records: %s", err)
	}

	return &handleWebsiteResponse{
		Site:  site,
		Pages: pages,
	}, nil
}

func (c *Customer) VectorizeWebsite(ctx context.Context, txn queries.DBTX, site *queries.Website) error {
	logger := c.logger.With("site.ID", site.ID.String(), "site.Domain", site.Domain)
	logger.InfoContext(ctx, "Parsing site ...")

	// get the pages
	model := queries.New(txn)
	pages, err := model.GetWebsitePagesBySite(ctx, site.ID)
	if err != nil {
		return fmt.Errorf("there was an issue fetching the sites: %v", err)
	}

	logger.InfoContext(ctx, "Creating embeddings for each page ...")

	// create an embeddings object
	emb := llm.GetEmbeddings(logger, c.Customer)

	// loop and perform the vectorization
	var wg sync.WaitGroup

	// create a results slice for the data
	results := make([]*vectorizeWebsiteResult, len(pages))
	tokenRecsChannel := make(chan *tokens.UsageRecord, len(pages))
	errs := make(chan error, len(pages))

	for i, item := range pages {
		wg.Add(1)
		go func(index int, page *queries.WebsitePage) {
			defer wg.Done()
			pLogger := logger.With("page", page.Url)

			pLogger.InfoContext(ctx, "Scraping the page ...")

			// scrape the webpage
			response, err := webparse.ScrapeSingle(ctx, pLogger, page)
			if err != nil {
				errs <- fmt.Errorf("error scraping the site: %v", err)
				return
			}

			// create a signature to compare the old vs new
			sig := utils.GenerateFingerprint([]byte(response.Content))
			if page.Sha256 == sig {
				pLogger.InfoContext(ctx, "This website page has not changed")
				return
			} else {
				pLogger.InfoContext(ctx, "The signatures do not match", "oldSHA256", page.Sha256, "newSHA256", sig)
			}

			pLogger.InfoContext(ctx, "Vecorizing the content ...")

			// embed the content
			res, err := emb.Embed(ctx, response.Content)
			if err != nil {
				errs <- fmt.Errorf("error embedding the content: %v", err)
				return
			}

			tokenRecsChannel <- res.Usage

			// write to index in the list
			results[index] = &vectorizeWebsiteResult{
				page:    page,
				headers: response.Header,
				sha256:  sig,
				vectors: res.Embeddings,
			}

			pLogger.InfoContext(ctx, "Successfully processed page")
		}(i, item)
	}

	logger.InfoContext(ctx, "Processed all pages")

	// wait for the routines to finish
	wg.Wait()
	close(errs)
	close(tokenRecsChannel)

	// report vectors
	tokenRecords := make([]*tokens.UsageRecord, 0)
	for item := range tokenRecsChannel {
		tokenRecords = append(tokenRecords, item)
	}
	if err := utils.ReportUsage(ctx, logger, txn, c.ID, tokenRecords, nil); err != nil {
		logger.ErrorContext(ctx, "Failed to log vector usage", "err", err)
	}

	// parse the errors
	var runtimeErr error
	for err := range errs {
		runtimeErr = err
		logger.ErrorContext(ctx, "there was an error vectorizing data", "error", runtimeErr)
	}
	if runtimeErr != nil {
		return fmt.Errorf("there was an issue during vecorization: %v", runtimeErr)
	}

	logger.InfoContext(ctx, "Parsing the result ...")

	// parse the results
	for _, item := range results {
		// ignore pages that did not change
		if item == nil {
			continue
		}

		plogger := logger.With("page", *item.page)

		plogger.InfoContext(ctx, "Updating the web page signature")
		newPage, err := model.UpdateWebsitePageSignature(ctx, &queries.UpdateWebsitePageSignatureParams{
			ID:     item.page.ID,
			Sha256: item.sha256,
		})
		if err != nil {
			return fmt.Errorf("there was an issue updating the page sha256: %v", err)
		}
		item.page = newPage

		// encode the headers from the website page to store with the vectors
		encodedHeaders, err := json.Marshal(item.headers)
		if err != nil {
			return fmt.Errorf("failed to serialize headers: %s", err)
		}

		plogger.InfoContext(ctx, "Uploading page vectors ...")

		// lastly upload the vectors to the datastore
		for index, vec := range item.vectors {
			// create raw vector object
			vecId, err := model.CreateVector(ctx, &queries.CreateVectorParams{
				Raw:        vec.Raw,
				Embeddings: &vec.Embedding,
				CustomerID: c.ID,
				Metadata:   encodedHeaders,
			})
			if err != nil {
				return fmt.Errorf("error inserting the vector object: %s", err)
			}

			// create a reference to the vector onto the document
			_, err = model.CreateWebsitePageVector(ctx, &queries.CreateWebsitePageVectorParams{
				WebsitePageID: item.page.ID,
				VectorStoreID: vecId,
				CustomerID:    c.ID,
				Index:         int32(index),
				Metadata:      encodedHeaders,
			})
			if err != nil {
				return fmt.Errorf("error creating document vector relationship: %s", err)
			}
		}

		plogger.InfoContext(ctx, "Successfully uploaded page")
	}

	logger.InfoContext(ctx, "Successfully vectorized site")

	return nil
}

func (c *Customer) VectorizeAllWebsites(ctx context.Context, txn queries.DBTX) error {
	c.logger.InfoContext(ctx, "Querying all sites ...")
	// get all the sites
	model := queries.New(txn)
	sites, err := model.GetWebsitesByCustomer(ctx, c.ID)
	if err != nil {
		return fmt.Errorf("failed to get websites: %s", err)
	}

	c.logger.InfoContext(ctx, "Processing all sites ...")

	// process all sites
	for _, site := range sites {
		if err := c.VectorizeWebsite(ctx, txn, site); err != nil {
			return fmt.Errorf("error vectorizing site: %s", err)
		}
	}

	c.logger.InfoContext(ctx, "Successfully vectorized all sites")

	return nil
}

func (c *Customer) PurgeDatastore(
	ctx context.Context,
	txn queries.DBTX,
	request *purgeDatastoreRequest,
) error {
	var err error
	logger := c.logger.With()

	// get the datastore
	store, err := c.GetDocstore(ctx)
	if err != nil {
		return fmt.Errorf("error getting the docstore: %s", err)
	}

	// parse the request (default is )
	timestamp := time.Now().UTC().Add(time.Minute * -10)
	if request.Timestamp != nil {
		// parse the time
		time, err := time.Parse("2006-01-02 15:04:05", *request.Timestamp)
		if err != nil {
			return fmt.Errorf("error parsing the time: %s", err)
		}
		timestamp = time
		logger.InfoContext(ctx, "passed timestamp", "timestamp", time)
	}
	// encode into sql type
	var pgtime pgtype.Timestamptz
	err = pgtime.Scan(timestamp)
	if err != nil {
		return fmt.Errorf("error encoding the timestamp into an sql type: %s", err)
	}

	logger = logger.With("timestamp", timestamp)

	model := queries.New(txn)

	// get all documents older than
	docs, err := model.GetDocumentsOlderThan(ctx, &queries.GetDocumentsOlderThanParams{
		CustomerID: c.ID,
		UpdatedAt:  pgtime,
	})
	if err != nil {
		return fmt.Errorf("error getting documents: %s", err)
	}

	logger.InfoContext(ctx, "Attempting to delete all documents from remote datastore", "length", len(docs))

	// delete in go-routines
	var wg sync.WaitGroup
	failedDocIds := make(chan uuid.UUID)
	for _, item := range docs {
		wg.Add(1)
		go func(doc queries.Document) {
			defer wg.Done()
			l := logger.With("doc", doc)
			dsDoc, err := datastore.NewDocumentFromDocument(ctx, l, &doc)
			if err != nil {
				l.ErrorContext(ctx, "failed to create the docstore doc", "error", err)
				failedDocIds <- doc.ID
				return
			}
			l.InfoContext(ctx, "Attempting to delete from datastore")
			if err := store.DeleteFile(ctx, dsDoc.DatastoreID); err != nil {
				l.ErrorContext(ctx, "error deleting from datastore", "error", err)
				failedDocIds <- doc.ID
				return
			}
			l.InfoContext(ctx, "successfully deleted from datastore")
		}(*item)
	}

	wg.Wait()
	close(failedDocIds)

	logger.InfoContext(ctx, "Successfully deleted documents")

	// TODO -- add some better error handling here. For now, ignore failed state in S3
	// for id := range failedDocIds {

	// }

	// get all folders older than
	folders, err := model.GetFoldersOlderThan(ctx, &queries.GetFoldersOlderThanParams{
		CustomerID: c.ID,
		UpdatedAt:  pgtime,
	})
	if err != nil {
		return fmt.Errorf("error getting folders: %s", err)
	}

	logger.InfoContext(ctx, "Attempting to delete all folders from remote datastore", "length", len(folders))

	failedFolderIds := make(chan uuid.UUID)
	for _, item := range folders {
		wg.Add(1)
		go func(folder queries.Folder) {
			defer wg.Done()
			l := logger.With("folder", folder)
			l.InfoContext(ctx, "Attempting to delete from datastore")
			if err := store.DeleteFile(ctx, fmt.Sprintf("%d", folder.ID)); err != nil {
				l.ErrorContext(ctx, "error deleting from datastore", "error", err)
				failedFolderIds <- folder.ID
				return
			}
			l.InfoContext(ctx, "successfully deleted from datastore")
		}(*item)
	}

	wg.Wait()
	close(failedFolderIds)

	logger.InfoContext(ctx, "Successfully deleted folders")
	logger.InfoContext(ctx, "Purging all records from DB before timestamp ...", "timestamp", timestamp)

	// purge all files
	err = model.DeleteDocumentsOlderThan(ctx, &queries.DeleteDocumentsOlderThanParams{
		CustomerID: c.ID,
		UpdatedAt:  pgtime,
	})
	if err != nil {
		return fmt.Errorf("error deleting documents: %s", err)
	}

	// purge all folders
	logger.Info("deleting stale folders ...")
	err = model.DeleteFoldersOlderThan(ctx, &queries.DeleteFoldersOlderThanParams{
		CustomerID: c.ID,
		UpdatedAt:  pgtime,
	})
	if err != nil {
		return fmt.Errorf("error deleting folders: %s", err)
	}

	// purge all website pages
	logger.Info("deleting stale website pages ...")
	if err := model.DeleteWebsitePagesOlderThan(ctx, &queries.DeleteWebsitePagesOlderThanParams{
		CustomerID: c.ID,
		UpdatedAt:  pgtime,
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to delete the website pages", err)
	}

	// delete all websites where page count is 0
	logger.Info("deleting stale websites ...")
	if err := model.DeleteWebsiteEmpty(ctx, c.ID); err != nil {
		return slogger.Error(ctx, logger, "failed to delete the empty websites", err)
	}

	logger.InfoContext(ctx, "Successfully purged datastore")
	return nil
}

func (c *Customer) VectorizeDatastoreOLD(
	ctx context.Context,
	txn queries.DBTX,
) error {
	logger := c.logger.With()
	logger.InfoContext(ctx, "Vectorizing datastore ...")

	// get the embeddings
	emb := llm.GetEmbeddings(logger, c.Customer)

	// create the model object
	model := queries.New(txn)

	logger.InfoContext(ctx, "Getting all documents for the customer")

	// get all documents
	docs, err := model.GetDocumentsByCustomer(ctx, c.ID)
	if err != nil {
		return fmt.Errorf("error getting all the documents: %s", err)
	}

	// create a type to store the embeddings data
	type embeddingResponse struct {
		doc     *datastore.Document
		vectors []*ltypes.EmbeddingsData
	}

	// create an array for the results to be stored in
	vectors := make([]*embeddingResponse, len(docs))

	logger.InfoContext(ctx, "Processing all documents ...", "length", len(docs))

	// process all docs in paralell
	var wg sync.WaitGroup
	tokenRecsChannel := make(chan *tokens.UsageRecord)
	errCh := make(chan error)

	for i, item := range docs {
		wg.Add(1)
		go func(index int, d queries.Document) {
			defer wg.Done()
			l := logger.With("doc", d)

			doc, err := datastore.NewDocumentFromDocument(ctx, l, &d)
			if err != nil {
				l.ErrorContext(ctx, "failed parsing the database doc", "err", err)
				errCh <- err
				return
			}

			// get the cleaned data from the document and a parser
			l.InfoContext(ctx, "Fetching document from datastore ...")
			cleaned, err := doc.GetCleaned(ctx)
			if err != nil {
				l.ErrorContext(ctx, "there was an issue getting the cleaned document", "err", err)
				errCh <- err
				return
			}

			// TODO -- there is potential to get this to work
			// see if the document changed
			// sig := utils.GenerateFingerprint(d.Data)
			// if doc.Sha256 == sig {
			// 	l.InfoContext(ctx, "This document has not changed")
			// 	return
			// } else {
			// 	l.InfoContext(ctx, "The signatures do not match", "old", doc.Sha256, "new", sig)
			// }

			// embed the content
			l.InfoContext(ctx, "Embedding the document ...")
			res, err := emb.Embed(ctx, cleaned.String())
			if err != nil {
				l.ErrorContext(ctx, "error embedding the content", "err", err)
				errCh <- err
				return
			}
			tokenRecsChannel <- res.Usage

			// add to the result
			vectors[index] = &embeddingResponse{
				doc:     doc,
				vectors: res.Embeddings,
			}

			l.InfoContext(ctx, "Successfully processed document")
		}(i, *item)
	}

	// wait for threads
	wg.Wait()
	close(errCh)
	close(tokenRecsChannel)

	// report vectors
	tokenRecords := make([]*tokens.UsageRecord, 0)
	for item := range tokenRecsChannel {
		tokenRecords = append(tokenRecords, item)
	}
	if err := utils.ReportUsage(ctx, logger, txn, c.ID, tokenRecords, nil); err != nil {
		logger.ErrorContext(ctx, "Failed to log vector usage", "err", err)
	}

	logger.InfoContext(ctx, "Successfully processed all documents")

	// check for errors. If one exists in the channel, something went wrong
	for err := range errCh {
		return fmt.Errorf("error vectorizing the data: %s", err)
	}

	// loop over the results and insert the vectors into the database
	logger.InfoContext(ctx, "Inserting all documents into the database")
	for _, item := range vectors {
		// skip errors, though this should not be hit
		if item == nil {
			continue
		}

		l := logger.With("filename", item.doc.Filename, "documentId", item.doc.ID)

		l.InfoContext(ctx, "Inserting document vectors ...")

		// loop over all vector results
		for index, vec := range item.vectors {
			l.DebugContext(ctx, "Processing index", "index", index)
			// create raw vector object
			vecId, err := model.CreateVector(ctx, &queries.CreateVectorParams{
				Raw:        vec.Raw,
				Embeddings: &vec.Embedding,
				CustomerID: c.ID,
			})
			if err != nil {
				return fmt.Errorf("error inserting the vector object: %s", err)
			}

			// create a reference to the vector onto the document
			_, err = model.CreateDocumentVector(ctx, &queries.CreateDocumentVectorParams{
				DocumentID:    item.doc.ID,
				VectorStoreID: vecId,
				CustomerID:    c.ID,
				Index:         int32(index),
			})
			if err != nil {
				return fmt.Errorf("error creating document vector relationship: %s", err)
			}
		}
		l.InfoContext(ctx, "Finished.")
	}

	logger.InfoContext(ctx, "Successfully vectorized the datastore")

	return nil
}

// Deletes ALL objects from the remote datastore. This is more for use in tests
// to quickly delete all objects from the datastore
func (c *Customer) DeleteRemoteDatastore(ctx context.Context) error {
	c.logger.InfoContext(ctx, "Deleting all objects from the remote datastore")

	// get the datastore
	store, err := c.GetDocstore(ctx)
	if err != nil {
		return fmt.Errorf("error getting docstore: %s", err)
	}

	// send request, all customers have a root folder with the id
	if err := store.DeleteRoot(ctx, fmt.Sprintf("%d/", c.Customer.ID)); err != nil {
		return fmt.Errorf("error deleting the root document")
	}

	c.logger.InfoContext(ctx, "Successfully deleted all objects in remote datastore for the customer")

	return nil
}
