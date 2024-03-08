package docstore

import (
	"context"
	"fmt"
	"testing"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/document"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/stretchr/testify/assert"
)

func TestS3Docstore(t *testing.T) {
	ctx := context.TODO()

	// create the database conn
	pool, err := db.GetPool()
	if err != nil {
		t.Error(err)
	}
	txn, err := pool.Begin(ctx)
	if err != nil {
		t.Error(err)
	}

	// create the test customer
	customer, err := testingutils.CreateTestCustomer(txn)
	if err != nil {
		t.Error(err)
	}

	// create the docstore object
	ds, err := NewS3Docstore(ctx, S3_BUCKET, nil)
	if err != nil {
		t.Error(err)
	}

	// dummy doc
	doc, err := document.NewDoc("helloworld.txt", []byte("This is some text from the document"))
	if err != nil {
		t.Error(err)
	}

	// insert the document
	url, err := ds.UploadDocument(ctx, customer, doc)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("URL:", url)

	// get the document
	retrievedDoc, err := ds.GetDocument(ctx, customer, doc.Filename)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, doc.Filename, retrievedDoc.Filename)
	assert.Equal(t, doc.Filetype, retrievedDoc.Filetype)
	assert.Equal(t, doc.Data, retrievedDoc.Data)

	// delete the document
	err = ds.DeleteDocument(ctx, customer, doc.Filename)
	assert.Nil(t, err)
}
