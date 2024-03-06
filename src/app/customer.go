package app

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/embeddings"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

// Wrapper around the `queries.Customer` object that represents the database object
// in order to store some state about the customer when needed
type Customer struct {
	*queries.Customer

	root   *queries.Folder
	logger *slog.Logger
}

func NewCustomer(ctx context.Context, logger *slog.Logger, id int64, txn pgx.Tx) (*Customer, error) {
	l := logger.With("customerId", id)
	model := queries.New(txn)

	l.InfoContext(ctx, "Fetching the customer record")

	// get the customer
	c, err := model.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	// get the root folder
	f, err := model.GetCustomerRootFolder(ctx, c.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get the root folder: %v", err)
	}

	return &Customer{
		Customer: &c,
		root:     &f,
		logger:   l,
	}, nil
}

type FolderContents struct {
	Self      *queries.Folder
	Folders   []queries.Folder
	Documents []queries.Document
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

func (c *Customer) CreateFolder(ctx context.Context, txn pgx.Tx, args *CreateFolderArgs) (*queries.Folder, error) {
	if args == nil {
		return nil, fmt.Errorf("the arguments cannot be nil")
	}
	if args.Name == "" {
		return nil, fmt.Errorf("the name cannot be empty")
	}
	if args.Owner == nil {
		args.Owner = c.root
	}

	logger := c.logger.With("name", args.Name)
	logger.InfoContext(ctx, "Creating a new folder ...")

	// create the folder
	model := queries.New(txn)
	folder, err := model.CreateFolder(ctx, queries.CreateFolderParams{
		ParentID:   pgtype.Int8{Int64: args.Owner.ID, Valid: true},
		CustomerID: c.ID,
		Title:      args.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("there was an issue creating the folder: %v", err)
	}

	return &folder, err
}

func (c *Customer) GetFolderContents(ctx context.Context, txn pgx.Tx, folder *queries.Folder) (*FolderContents, error) {
	logger := c.logger.With("folder.id", folder.ID, "folder.title", folder.Title)
	logger.InfoContext(ctx, "Listing folder contents")

	model := queries.New(txn)

	// get the folders
	folders, err := model.GetFoldersFromParent(ctx, pgtype.Int8{Int64: folder.ID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the folders: %v", err)
	}

	// get the documents
	documents, err := model.GetDocumentsFromParent(ctx, folder.ID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the documents: %v", err)
	}

	return &FolderContents{
		Self:      folder,
		Folders:   folders,
		Documents: documents,
	}, nil
}

/*
Uploads documents to the customer's datastore, and stores the reference to the object in
the database, but does NOT vectorize the data. The vectorization is done by the function
`ReVectorizeDatastore`.
*/
func (c *Customer) UploadDocuments(ctx context.Context, txn pgx.Tx, folder *queries.Folder, docs []*docstore.Doc) ([]*UploadDocumentsResponse, error) {
	logger := c.logger.With("folder", folder.ID, "numDocuments", len(docs))
	logger.InfoContext(ctx, "Uploading documents ...")

	// create the docstore object based on the customer
	store, err := c.GetDocstore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docstore: %v", err)
	}

	// create the embeddings
	em := embeddings.NewOpenAIEmbeddings(fmt.Sprint(c.ID), &embeddings.OpenAIEmbeddingsOpts{
		Logger: logger,
	})

	// track responses
	response := make([]*UploadDocumentsResponse, len(docs))

	// upload all the documents to postgres
	for idx, item := range docs {
		logger.DebugContext(ctx, "Processing document ...", "filename", item.Filename)
		response[idx] = &UploadDocumentsResponse{}

		// pseudo transaction for the document
		tx, err := txn.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("there was an issue creating the pseudo-transaction: %v", err)
		}

		// create the queries model
		model := queries.New(tx)

		// create an sha256 fingerprint for the document
		hash := sha256.Sum256(item.Data)

		// upload to postgres
		doc, err := model.CreateDocument(ctx, queries.CreateDocumentParams{
			ParentID:   folder.ID,
			CustomerID: c.ID,
			Filename:   item.Filename,
			Type:       string(item.Filetype),
			SizeBytes:  int64(len(item.Data)),
			Sha256:     fmt.Sprintf("%x", hash),
		})
		if err != nil {
			// rollback
			logger.ErrorContext(ctx, "There was an error uploading to postgres", "error", err)
			response[idx].Error = err
			err := tx.Rollback(ctx)
			if err != nil {
				return nil, fmt.Errorf("CRITICAL failed to rollback: %v", err)
			}
			continue
		}
		response[idx].Doc = &doc

		// create embeddings for the document
		vectors, err := em.Embed(ctx, string(item.Data))
		if err != nil {
			// rollback
			logger.ErrorContext(ctx, "There was an error creating the vectors", "error", err)
			response[idx].Error = err
			err := tx.Rollback(ctx)
			if err != nil {
				return nil, fmt.Errorf("CRITICAL failed to rollback: %v", err)
			}
			continue
		}

		// create vector objects for all the vectors retrieved
		for idx, v := range vectors {
			_, err := model.CreateVector(ctx, queries.CreateVectorParams{
				Raw:        v.Raw,
				Embeddings: v.Embedding,
				CustomerID: c.ID,
				DocumentID: doc.ID,
				Index:      int32(idx),
			})
			if err != nil {
				// rollback
				logger.ErrorContext(ctx, "There was an error inserting the vector", "vectorIndex", idx, "error", err)
				response[idx].Error = err
				err := tx.Rollback(ctx)
				if err != nil {
					return nil, fmt.Errorf("CRITICAL failed to rollback: %v", err)
				}
				continue
			}
		}

		// create the model records for the user
		err = em.ReportUsage(ctx, tx, c.Customer)
		if err != nil {
			// not an error worth failing on
			logger.ErrorContext(ctx, "There was an error reporting the usage", "error", err)
		}

		// upload to docstore
		url, err := store.UploadDocument(ctx, c.Customer, item)
		if err != nil {
			// rollback
			logger.ErrorContext(ctx, "There was an error uploading to the document store", "error", err)
			response[idx].Error = err
			err := tx.Rollback(ctx)
			if err != nil {
				return nil, fmt.Errorf("CRITICAL failed to rollback: %v", err)
			}
			continue
		}
		response[idx].Url = url
	}

	return response, nil
}

/*
Deletes a document from the datastore and its vectorization data.
*/
func (c *Customer) DeleteDocument(ctx context.Context, doc *docstore.Doc) {

}

/*
Performs a complete re-vectorization of all the objects that have changed inside the
document store. Compares the objects that are already inside the datastore and their
sha256 values. If they are equal, then nothing is done. If they are not equal, the old
object is deleted and then re-vectorized. If the object in the datastore does not exist
anymore, the vector data is deleted. This operation is quite expensive from compute and
api costs, so the customer should be wary to run this function often
*/
// func (c *Customer) ReVectorizeDatastore(ctx context.Context) {

// }
