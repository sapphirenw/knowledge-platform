package docstore

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"
)

func TestS3Docstore(t *testing.T) {
	ctx := context.Background()
	logger := utils.DefaultLogger()
	store, err := NewS3Docstore(ctx, S3_BUCKET, logger)
	require.NoError(t, err)

	cid, err := uuid.NewV7()
	require.NoError(t, err)
	customer := &queries.Customer{
		ID: cid,
	}

	// read a document
	filename := "s3.txt"
	data, err := os.ReadFile("../../resources/" + filename)
	require.NoError(t, err)

	filetype, err := ParseFileType(filename)
	require.NoError(t, err)

	// create a doc from this data
	doc, err := NewDocumentFromRaw(customer, filename, data)
	require.NoError(t, err)

	// create the pre-signed url
	url, err := store.GeneratePresignedUrl(ctx, doc)
	require.NoError(t, err)

	// create the upload request
	request, err := http.NewRequest(store.GetUploadMethod(), string(url), bytes.NewReader(data))
	require.NoError(t, err)

	// set the headers
	request.Header.Set("Content-Type", string(filetype))
	client := &http.Client{}

	// send the request
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("the status code was not 400")
		fmt.Println(response.StatusCode)
		fmt.Println(response)
		t.FailNow()
	}

	// delete the doc
	err = store.DeleteFile(ctx, doc.UniqueID)
	require.NoError(t, err)
}
