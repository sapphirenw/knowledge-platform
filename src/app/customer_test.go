package app

import (
	"context"
	"fmt"
	"os"
	"testing"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestCustomerFolderStructure(t *testing.T) {
	ctx := context.TODO()
	logger := utils.DefaultLogger()

	// get the db pool
	pool, err := db.GetPool()
	if err != nil {
		t.Error(err)
		return
	}

	// start a txn
	txn, err := pool.Begin(ctx)
	if err != nil {
		t.Error(t)
	}
	defer txn.Rollback(ctx)

	// get the customer
	c, err := testingutils.CreateTestCustomer(txn)
	if err != nil {
		t.Error(err)
		return
	}

	// create the wrapper customer object
	customer, err := NewCustomer(ctx, logger, c.ID, txn)
	if err != nil {
		t.Error(err)
		return
	}

	// create 2 folders in the root
	f1, err := customer.CreateFolder(ctx, txn, &CreateFolderArgs{
		Owner: customer.root,
		Name:  "folder1",
	})
	if err != nil {
		t.Error(err)
		return
	}
	f2, err := customer.CreateFolder(ctx, txn, &CreateFolderArgs{
		Owner: customer.root,
		Name:  "folder2",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// create a folder in folder2
	_, err = customer.CreateFolder(ctx, txn, &CreateFolderArgs{
		Owner: f2,
		Name:  "folder3",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// describe the folders
	resp1, err := customer.GetFolderContents(ctx, txn, customer.root)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 2, len(resp1.Folders))
	assert.Empty(t, resp1.Documents)

	resp2, err := customer.GetFolderContents(ctx, txn, f1)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Empty(t, resp2.Folders)
	assert.Empty(t, resp2.Documents)

	resp3, err := customer.GetFolderContents(ctx, txn, f2)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 1, len(resp3.Folders))
	assert.Empty(t, resp3.Documents)

}

func TestCustomerInsertDocuments(t *testing.T) {
	ctx := context.TODO()
	logger := utils.DefaultLogger()

	// get the db pool
	pool, err := db.GetPool()
	if err != nil {
		t.Error(err)
		return
	}

	// start a txn
	txn, err := pool.Begin(ctx)
	if err != nil {
		t.Error(t)
	}
	defer txn.Rollback(ctx)

	// get the customer
	c, err := testingutils.CreateTestCustomer(txn)
	if err != nil {
		t.Error(err)
		return
	}

	// create the wrapper customer object
	customer, err := NewCustomer(ctx, logger, c.ID, txn)
	if err != nil {
		t.Error(err)
		return
	}

	// create a doc
	data, err := os.ReadFile("/Users/jakelanders/code/ai-content-creation/api/resources/oregon.txt")
	if err != nil {
		t.Error(err)
		return
	}
	doc := &docstore.Doc{
		Filename: "test.txt",
		Filetype: docstore.FT_txt,
		Data:     data,
	}

	// upload the documents
	resp, err := customer.UploadDocuments(ctx, txn, customer.root, []*docstore.Doc{doc})
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range resp {
		if item.Error != nil {
			t.Error(item.Error)
		}
		fmt.Println("URL:", item.Url)
		fmt.Println("ID:", item.Doc.ID)
	}
	assert.NotEmpty(t, resp)
}
