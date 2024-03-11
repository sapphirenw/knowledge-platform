package customer

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/document"
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
	customer, err := NewCustomer(ctx, logger, testingutils.TEST_CUSTOMER_ID, txn)
	if err != nil {
		t.Error(err)
		return
	}

	// create 2 folders in the root
	f1, err := customer.CreateFolder(ctx, txn, &createFolderRequest{
		Owner: nil,
		Name:  "folder1",
	})
	if err != nil {
		t.Error(err)
		return
	}
	f2, err := customer.CreateFolder(ctx, txn, &createFolderRequest{
		Owner: nil,
		Name:  "folder2",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// create a folder in folder2
	_, err = customer.CreateFolder(ctx, txn, &createFolderRequest{
		Owner: f2,
		Name:  "folder3",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// describe the folders
	resp1, err := customer.ListFolderContents(ctx, txn, nil)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 2, len(resp1.Folders))
	assert.Empty(t, resp1.Documents)

	resp2, err := customer.ListFolderContents(ctx, txn, f1)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Empty(t, resp2.Folders)
	assert.Empty(t, resp2.Documents)

	resp3, err := customer.ListFolderContents(ctx, txn, f2)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 1, len(resp3.Folders))
	assert.Empty(t, resp3.Documents)

}

func TestCustomerUploadDocument(t *testing.T) {
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

	// create the wrapper customer object
	customer, err := NewCustomer(ctx, logger, testingutils.TEST_CUSTOMER_ID, txn)
	if err != nil {
		t.Error(err)
		return
	}

	// create a doc
	filename := "file1.txt"
	data, err := os.ReadFile(fmt.Sprintf("../../resources/%s", filename))
	if err != nil {
		t.Error(err)
		return
	}

	// create a new document type
	doc, err := document.NewDoc(filename, data)
	if err != nil {
		t.Error(err)
		return
	}

	// create the pre-signed url
	uploadInput := generatePresignedUrlRequest{
		Filename:  doc.Filename,
		Mime:      string(doc.Filetype),
		Signature: utils.GenerateFingerprint(doc.Data),
		Size:      int64(len(doc.Data)),
	}
	preSignedResp, err := customer.GeneratePresignedUrl(ctx, txn, &uploadInput)
	if err != nil {
		t.Error(err)
		return
	}

	// decode the request url
	url, err := base64.StdEncoding.DecodeString(preSignedResp.UploadUrl)
	if err != nil {
		t.Error(err)
	}

	// create the upload request
	request, err := http.NewRequest(preSignedResp.Method, string(url), bytes.NewReader(doc.Data))
	if err != nil {
		t.Error(err)
		return
	}

	// set the headers
	request.Header.Set("Content-Type", string(doc.Filetype))
	client := &http.Client{}

	// send the request
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Error("the status code was not 400")
		fmt.Println(response.StatusCode)
		fmt.Println(response)
	}

	// notify the server of the success
	err = customer.NotifyOfSuccessfulUpload(ctx, txn, preSignedResp.DocumentId)
	if err != nil {
		t.Error(err)
		return
	}

	// delete the document from the docstore
	store, err := customer.GetDocstore(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	err = store.DeleteDocument(ctx, customer.Customer, doc.Filename)
	if err != nil {
		t.Error(err)
		return
	}
}
