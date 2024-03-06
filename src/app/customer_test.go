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
