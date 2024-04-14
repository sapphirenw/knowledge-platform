package customer

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/document"
	"github.com/sapphirenw/ai-content-creation-api/src/embeddings"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/webscrape"
)

// Wrapper around the `queries.Customer` object that represents the database object
// in order to store some state about the customer when needed
type Customer struct {
	*queries.Customer

	// root   *queries.Folder
	logger *slog.Logger
}

func NewCustomer(ctx context.Context, logger *slog.Logger, id int64, db queries.DBTX) (*Customer, error) {
	model := queries.New(db)

	logger.InfoContext(ctx, "Fetching the customer record")

	// get the customer
	c, err := model.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	// get the root folder
	// f, err := model.GetCustomerRootFolder(ctx, c.ID)
	// if err != nil {
	// 	if err.Error() == "no rows in result set" {
	// 		// attempt to create a new root folder
	// 		logger.InfoContext(ctx, "No root folder was found, attempting to create one...")
	// 		f, err = model.CreateFolderRoot(ctx, c.ID)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to create the root folder: %v", err)
	// 		}
	// 		logger.InfoContext(ctx, "Successfully created the root folder")
	// 	} else {
	// 		return nil, fmt.Errorf("could not get the root folder: %v", err)
	// 	}
	// }

	return &Customer{
		Customer: c,
		// root:     f,
		logger: logger.With("customer.ID", c.ID, "customer.Name", c.Name, "customer.Datastore", c.Datastore),
	}, nil
}

// Gets the docstore associated with the customer
func (c *Customer) GetDocstore(ctx context.Context) (docstore.Docstore, error) {
	switch c.Customer.Datastore {
	case "s3":
		return docstore.NewS3Docstore(ctx, docstore.S3_BUCKET, c.logger)
	default:
		return docstore.NewTODODocstore(c.logger)
	}
}

// Creates a new folder tied to the customer with an optional parent.
func (c *Customer) CreateFolder(ctx context.Context, txn pgx.Tx, args *createFolderRequest) (*queries.Folder, error) {
	if args == nil {
		return nil, fmt.Errorf("the arguments cannot be nil")
	}
	if args.Name == "" {
		return nil, fmt.Errorf("the name cannot be empty")
	}

	logger := c.logger.With("folder.Name", args.Name)

	// parse the owner if applicable
	var parentId pgtype.Int8
	if args.Owner != nil {
		logger = c.logger.With("folder.Owner", args.Owner.ID)
		parentId.Scan(args.Owner.ID)
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
		return nil, fmt.Errorf("there was an issue creating the folder: %v", err)
	}
	logger.InfoContext(ctx, "Successfully created folder", "id", folder.ID)

	return folder, err
}

// Does an 'ls' on a folder
func (c *Customer) ListFolderContents(ctx context.Context, db queries.DBTX, folder *queries.Folder) (*listFolderContentsResponse, error) {
	logger := c.logger.With()
	if folder != nil {
		logger = logger.With("folder.ID", folder.ID, "folder.Title", folder.Title)
	}
	logger.InfoContext(ctx, "Getting all children of the folder ...")

	model := queries.New(db)

	var err error
	var folders []*queries.Folder
	var documents []*queries.Document

	if folder == nil {
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
		// get the information using the folder
		folders, err = model.GetFoldersFromParent(ctx, pgtype.Int8{Int64: folder.ID, Valid: true})
		if err != nil {
			return nil, fmt.Errorf("there was an issue getting the folders: %v", err)
		}
		documents, err = model.GetDocumentsFromParent(ctx, pgtype.Int8{Int64: folder.ID, Valid: true})
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
Generates pre-signed urls for the user to use to upload to their preferred datastore. This does not
have any state-chaning effects, as no records are inserted into the database, and no objects
*/
func (c *Customer) GeneratePresignedUrl(ctx context.Context, db queries.DBTX, body *generatePresignedUrlRequest) (*generatePresignedUrlResponse, error) {
	logger := c.logger.With("body", *body)
	logger.InfoContext(ctx, "Generating a presigned url...")

	// get the customer's docstore
	store, err := c.GetDocstore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the document store: %v", err)
	}

	var parentId pgtype.Int8
	if body.ParentId != nil {
		parentId.Scan(body.ParentId)
	}

	// insert a record into the documents table
	model := queries.New(db)
	d, err := model.CreateDocument(ctx, &queries.CreateDocumentParams{
		ParentID:   parentId,
		CustomerID: c.ID,
		Filename:   body.Filename,
		Type:       body.Mime,
		SizeBytes:  body.Size,
		Sha256:     body.Signature,
	})
	if err != nil {
		return nil, fmt.Errorf("there was an issue creating the document: %v", err)
	}

	// generate the pre-signed url
	url, err := store.GeneratePresignedUrl(ctx, c.Customer, &docstore.UploadUrlInput{
		Filename:  body.Filename,
		Mime:      body.Mime,
		Signature: body.Signature,
		Size:      body.Size,
	})
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
func (c *Customer) NotifyOfSuccessfulUpload(ctx context.Context, db queries.DBTX, documentId int64) error {
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
Deletes a document from the datastore and its vectorization data.
*/
func (c *Customer) DeleteDocument(ctx context.Context, doc *document.Doc) {

}

/*
Performs a complete re-vectorization of all the objects that have changed inside the
document store. Compares the objects that are already inside the datastore and their
sha256 values. If they are equal, then nothing is done. If they are not equal, the old
object is deleted and then re-vectorized. If the object in the datastore does not exist
anymore, the vector data is deleted. This operation is quite expensive from compute and
api costs, so the customer should be wary to run this function often
*/
func (c *Customer) ReVectorizeDatastore(ctx context.Context) {
	// fetch all documents from database

	// loop over all documents

	// query for document in s3

	// if exists
	// query vectors. If empty, vectorize. If not,
	// compare fingerprints
	// if different, re-vectorize
	// if same, do nothing

	// if not exists
	// remove document from datastore

	//
	//
	//

	// create the embeddings
	// em := embeddings.NewOpenAIEmbeddings(fmt.Sprint(c.ID), &embeddings.OpenAIEmbeddingsOpts{
	// 	Logger: logger,
	// })

	// create embeddings for the document
	// vectors, err := em.Embed(ctx, string(item.Data))
	// if err != nil {
	// 	// rollback
	// 	logger.ErrorContext(ctx, "There was an error creating the vectors", "error", err)
	// 	response[idx].Error = err
	// 	err := tx.Rollback(ctx)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("CRITICAL failed to rollback: %v", err)
	// 	}
	// 	continue
	// }

	// // create vector objects for all the vectors retrieved
	// for idx, v := range vectors {
	// 	_, err := model.CreateVector(ctx, queries.CreateVectorParams{
	// 		Raw:        v.Raw,
	// 		Embeddings: v.Embedding,
	// 		CustomerID: c.ID,
	// 		DocumentID: doc.ID,
	// 		Index:      int32(idx),
	// 	})
	// 	if err != nil {
	// 		// rollback
	// 		logger.ErrorContext(ctx, "There was an error inserting the vector", "vectorIndex", idx, "error", err)
	// 		response[idx].Error = err
	// 		err := tx.Rollback(ctx)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("CRITICAL failed to rollback: %v", err)
	// 		}
	// 		continue
	// 	}
	// }

	// create the model records for the user
	// err = em.ReportUsage(ctx, tx, c.Customer)
	// if err != nil {
	// 	// not an error worth failing on
	// 	logger.ErrorContext(ctx, "There was an error reporting the usage", "error", err)
	// }
}

/*
Adds a website for the user, but does not scrape it.
*/
func (c *Customer) HandleWebsite(ctx context.Context, db queries.DBTX, request *handleWebsiteRequest) (*handleWebsiteResponse, error) {
	logger := c.logger.With("domain", request.Domain)
	logger.InfoContext(ctx, "Ingesting the domain...")

	// parse the domain
	protocol, domain, err := utils.ParseWebsiteInformation(request.Domain)
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
	urls, err := webscrape.ParseSitemap(ctx, logger, &tmpSite)
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

	return &handleWebsiteResponse{
		Site:  site,
		Pages: pages,
	}, nil
}

func (c *Customer) VectorizeWebsite(ctx context.Context, txn pgx.Tx, site *queries.Website) ([]*vectorizeWebsiteResult, error) {
	logger := c.logger.With("site.ID", site.ID, "site.Domain", site.Domain)
	logger.InfoContext(ctx, "Parsing site")

	// get the pages
	model := queries.New(txn)
	pages, err := model.GetWebsitePagesBySite(ctx, site.ID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue fetching the sites: %v", err)
	}

	// loop and perform the vectorization
	var wg sync.WaitGroup

	// create a results slice for the data
	results := make([]*vectorizeWebsiteResult, len(pages))
	errs := make(chan error, len(pages))

	for i, item := range pages {
		wg.Add(1)
		go func(index int, page *queries.WebsitePage) {
			defer wg.Done()
			pLogger := logger.With("page", page.Url)

			// scrape the webpage
			content, err := webscrape.ScrapeSingle(ctx, pLogger, page)
			if err != nil {
				errs <- fmt.Errorf("error scraping the site: %v", err)
				return
			}

			// create a signature to compare the old vs new
			sig := utils.GenerateFingerprint(content)
			if page.Sha256 == sig {
				pLogger.InfoContext(ctx, "This website page has not changed")
				return
			} else {
				pLogger.InfoContext(ctx, "The signatures do not match", "oldSHA256", page.Sha256, "newSHA256", sig)
			}

			// embed the content
			emb := embeddings.NewOpenAIEmbeddings(fmt.Sprintf("%d", c.ID), &embeddings.OpenAIEmbeddingsOpts{
				Logger: pLogger,
			})
			res, err := emb.Embed(ctx, string(content))
			if err != nil {
				errs <- fmt.Errorf("error embedding the content: %v", err)
				return
			}

			// write to index in the list
			results[index] = &vectorizeWebsiteResult{
				page:    page,
				sha256:  sig,
				vectors: res,
			}
		}(i, item)
	}

	// wait for the routines to finish
	wg.Wait()
	close(errs)

	// parse the errors
	var runtimeErr error
	for err := range errs {
		runtimeErr = err
		logger.ErrorContext(ctx, "there was an error vectorizing data", "error", runtimeErr)
	}
	if runtimeErr != nil {
		return nil, fmt.Errorf("there was an issue during vecorization: %v", runtimeErr)
	}

	// update any website page hashes that changed
	for _, item := range results {
		// ignore pages that did not change
		if item == nil {
			continue
		}

		logger.InfoContext(ctx, "Updating the web page signature", "page", item.page.Url)
		newPage, err := model.UpdateWebsitePageSignature(ctx, &queries.UpdateWebsitePageSignatureParams{
			ID:     item.page.ID,
			Sha256: item.sha256,
		})
		if err != nil {
			return nil, fmt.Errorf("there was an issue updating the page sha256: %v", err)
		}
		item.page = newPage
	}

	// return the results to the consumer
	return results, nil
}
