package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

// Wrapper around the `queries.Customer` object that represents the database object
// in order to store some state about the customer when needed
type Customer struct {
	*queries.Customer

	db     *queries.DBTX
	logger *slog.Logger
}

type FolderContents struct {
	Self      *queries.Folder
	Folders   []queries.Folder
	Documents []queries.Document
}

func (c *Customer) GetFolderContents(ctx context.Context, folder *queries.Folder) (*FolderContents, error) {
	logger := c.logger.With("folder.ID", folder.ID, "folder.Title", folder.Title)
	logger.InfoContext(ctx, "Listing folder contents")

	model := queries.New(*c.db)

	// get the folders
	folders, err := model.GetFoldersFromParent(ctx, folder.ID)
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
func (c *Customer) UploadDocuments(ctx context.Context, folder *queries.Folder, docs []*docstore.Doc) {

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
func (c *Customer) ReVectorizeDatastore(ctx context.Context) {

}
