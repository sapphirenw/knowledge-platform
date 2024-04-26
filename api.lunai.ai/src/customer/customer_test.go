package customer

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"
)

// import (
// 	"bytes"
// 	"context"
// 	"encoding/base64"
// 	"fmt"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/jackc/pgx/v5"
// 	db "github.com/sapphirenw/ai-content-creation-api/src/database"
// 	"github.com/sapphirenw/ai-content-creation-api/src/queries"
// 	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
// 	"github.com/sapphirenw/ai-content-creation-api/src/utils"
// 	"github.com/stretchr/testify/assert"
// )

func TestCustomerCreate(t *testing.T) {
	_, _, _, c := testInit(t)
	fmt.Printf("Created customerId: %d\n", c.ID)
}

func TestCustomerDocumentStore(t *testing.T) {
	ctx, _, pool, c := testInit(t)
	var err error

	// create the docstore
	err = uploadToDocstore(ctx, c, nil, "../../resources", pool)
	require.NoError(t, err)

	// vectorize the docstore
	err = c.VectorizeDatastore(ctx, pool)
	require.NoError(t, err)

	// TODO -- query docstore

	// remove the remote docstore
	err = c.DeleteRemoteDatastore(ctx)
	require.NoError(t, err)
}

func TestCustomerWebsites(t *testing.T) {
	ctx, _, pool, c := testInit(t)

	// parse some rules
	noRules, err := c.HandleWebsite(ctx, pool, &handleWebsiteRequest{
		Domain:    "https://crosschecksports.com",
		Blacklist: []string{},
		Whitelist: []string{},
		Insert:    false,
	})
	require.NoError(t, err)
	whitelist, err := c.HandleWebsite(ctx, pool, &handleWebsiteRequest{
		Domain:    "https://crosschecksports.com",
		Blacklist: []string{},
		Whitelist: []string{"docs"},
		Insert:    false,
	})
	require.NoError(t, err)
	blacklist, err := c.HandleWebsite(ctx, pool, &handleWebsiteRequest{
		Domain:    "https://crosschecksports.com",
		Blacklist: []string{"docs"},
		Whitelist: []string{},
		Insert:    false,
	})
	require.NoError(t, err)

	// assertions
	require.Less(t, len(whitelist.Pages), len(noRules.Pages))
	require.Less(t, len(blacklist.Pages), len(noRules.Pages))
	require.Less(t, len(whitelist.Pages), len(blacklist.Pages))

	// insert a website
	site, err := c.HandleWebsite(ctx, pool, &handleWebsiteRequest{
		Domain:    "https://crosschecksports.com",
		Blacklist: []string{},
		Whitelist: []string{"docs"},
		Insert:    true,
	})
	require.NoError(t, err)

	// vectorize the website
	err = c.VectorizeAllWebsites(ctx, pool)
	require.NoError(t, err)

	// query the vectors
	model := queries.New(pool)
	vectors, err := model.ListWebsitePageVectors(ctx, c.ID)
	require.NoError(t, err)

	// ensure the correct number of vectors was inserted
	rootVectors := 0
	for _, item := range vectors {
		if item.Index == 0 {
			rootVectors += 1
		}
	}
	require.Equal(t, len(site.Pages), rootVectors)

	// TODO -- run a query against the vector store
	vecQueryResponse, err := c.QueryVectorStore(ctx, pool, &queryVectorStoreRequest{
		Query: "How to create a team",
		K:     3,
	})
	require.NoError(t, err)

	fmt.Println("\n\n++++ DOCS:")
	for _, item := range vecQueryResponse.Documents {
		fmt.Println("- " + item.Filename)
	}

	fmt.Println("\n\n++++ PAGES:")
	for _, item := range vecQueryResponse.WebsitePages {
		fmt.Println("- " + item.Url)
	}

}

func uploadToDocstore(ctx context.Context, c *Customer, parentId *int64, directory string, db *pgxpool.Pool) error {
	// get all files in dir
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	// loop over all files
	for _, file := range files {
		if file.IsDir() {
			// create the folder
			var owner int64 = 0
			if parentId != nil {
				owner = *parentId
			}
			folder, err := c.CreateFolder(ctx, db, &createFolderRequest{
				Owner: owner,
				Name:  file.Name(),
			})
			if err != nil {
				return err
			}

			// parse the children
			err = uploadToDocstore(ctx, c, &folder.ID, fmt.Sprintf("%s/%s", directory, file.Name()), db)
			if err != nil {
				return err
			}
		} else {
			// read the file
			data, err := os.ReadFile(directory + "/" + file.Name())
			if err != nil {
				return err
			}

			filetype, err := docstore.ParseFileType(file.Name())
			if err != nil {
				return fmt.Errorf("failed to parse filetype: %s", err)
			}

			// create the pre-signed url
			preSignedResp, err := c.GeneratePresignedUrl(ctx, db, &generatePresignedUrlRequest{
				ParentId:  parentId,
				Filename:  file.Name(),
				Mime:      string(filetype),
				Signature: utils.GenerateFingerprint(data),
				Size:      int64(len(data)),
			})
			if err != nil {
				return err
			}

			// use the upload url to upload the doc
			// decode the request url
			url, err := base64.StdEncoding.DecodeString(preSignedResp.UploadUrl)
			if err != nil {
				return fmt.Errorf("failed to decode url: %s", err)
			}

			// create the upload request
			request, err := http.NewRequest(preSignedResp.Method, string(url), bytes.NewReader(data))
			if err != nil {
				return fmt.Errorf("failed to create the request: %s", err)
			}

			// set the headers
			request.Header.Set("Content-Type", string(filetype))
			client := &http.Client{}

			// send the request
			response, err := client.Do(request)
			if err != nil {
				return fmt.Errorf("failed to send the request: %s", err)
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				fmt.Println("the status code was not 400")
				fmt.Println(response.StatusCode)
				fmt.Println(response)
				return fmt.Errorf("the request failed")
			}

			// notify the server of the success
			if err := c.NotifyOfSuccessfulUpload(ctx, db, preSignedResp.DocumentId); err != nil {
				return fmt.Errorf("failed to notify of successful upload: %s", err)
			}
		}
	}

	return nil
}

func remotels(ctx context.Context, c *Customer, parentId *int64, indent int, db *pgxpool.Pool) error {
	response, err := c.ListFolderContents(ctx, db, parentId)
	if err != nil {
		return err
	}

	// print parent
	spacing := strings.Repeat(" ", indent)
	if response.Self == nil {
		fmt.Printf("%s - root/\n", spacing)
	} else {
		fmt.Printf("%s - %s/\n", spacing, response.Self.Title)
	}

	// print docs
	for _, item := range response.Documents {
		fmt.Printf("%s  - %s\n", spacing, item.Filename)
	}

	// print folders
	for _, item := range response.Folders {
		if err := remotels(ctx, c, &item.ID, indent+2, db); err != nil {
			return err
		}
	}

	return nil
}

// sets up the required structure for a test to run properly
func testInit(t *testing.T) (context.Context, *slog.Logger, *pgxpool.Pool, *Customer) {
	// base vars
	ctx := context.Background()
	logger := utils.DefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)

	// create the customer
	c, err := CreateCustomer(ctx, logger, pool, &createCustomerRequest{
		Name: "test-customer",
	})
	require.NoError(t, err)

	// create the customer object wrapper
	customer, err := NewCustomer(ctx, logger, c.ID, pool)
	require.NoError(t, err)

	return ctx, logger, pool, customer
}
