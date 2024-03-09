package customer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"testing"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
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

	// generate a fingerprint for the file
	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])

	// create the pre-signed url
	uploadInput := &docstore.UploadUrlInput{
		Filename:  doc.Filename,
		Mime:      string(doc.Filetype),
		Signature: hashHex,
		Size:      int64(len(doc.Data)),
	}
	preSignedResp, err := customer.GeneratePresignedUrl(ctx, uploadInput)
	if err != nil {
		t.Error(err)
		return
	}

	// create the upload request
	request, err := http.NewRequest(preSignedResp.Method, preSignedResp.UploadUrl, bytes.NewReader(doc.Data))
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
	err = customer.NotifyOfSuccessfulUpload(ctx, txn, customer.root, uploadInput)
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
