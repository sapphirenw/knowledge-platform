package docstore

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type S3Docstore struct {
	bucket string

	logger *slog.Logger
	config *aws.Config
}

func NewS3Docstore(ctx context.Context, bucket string, logger *slog.Logger) (*S3Docstore, error) {
	if logger == nil {
		logger = utils.DefaultLogger()
	}
	lgr := logger.With("s3_bucket", bucket)

	// setup aws config using env
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}

	return &S3Docstore{
		bucket: S3_BUCKET,
		logger: lgr,
		config: &config,
	}, nil
}

func (d *S3Docstore) UploadDocuments(ctx context.Context, customer *queries.Customer, docs []*Doc) []*UploadResponse {
	// create the s3 client
	client := s3.NewFromConfig(*d.config)
	uploader := manager.NewUploader(client)

	// loop and create all the objects
	responses := make([]*UploadResponse, 0)

	for _, doc := range docs {
		// create the response object
		documentId := createUniqueFileId(customer, doc.Filename)
		response := &UploadResponse{FileID: documentId}

		d.logger.InfoContext(ctx, "Uploading file", "filename", doc.Filename, "filetype", doc.Filetype)

		// upload with a docid
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    aws.String(documentId),
			Body:   bytes.NewReader(doc.Data),
		})
		if err != nil {
			response.Error = fmt.Errorf("there was an issue uploading the file: %v", err)
			responses = append(responses, response)
			continue
		}

		d.logger.InfoContext(ctx, "Successfully uploaded file", "filename", doc.Filename, "filetype", doc.Filetype)

		// add response metadata
		response.Document = doc
		response.Url = result.Location
		responses = append(responses, response)
	}

	return responses
}

func (d *S3Docstore) GetDocument(ctx context.Context, customer *queries.Customer, filename string) (*Doc, error) {
	d.logger.InfoContext(ctx, "Donwloading file from s3...", "filename", filename)

	// create the client
	client := s3.NewFromConfig(*d.config)
	downloader := manager.NewDownloader(client)

	buffer := manager.NewWriteAtBuffer([]byte{})

	// upload with a unqiue identifier
	fileId := createUniqueFileId(customer, filename)
	_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return nil, fmt.Errorf("there was an issue downloading the file from s3: %v", err)
	}

	// parse the filetype
	ft, err := parseFileType(filename)
	if err != nil {
		return nil, fmt.Errorf("there was an issue parsing the filetype: %v", err)
	}

	d.logger.InfoContext(ctx, "Successfully downloaded file", "filename", filename)

	return &Doc{
		Filename: parseUniqueFileId(customer, fileId),
		Filetype: ft,
		Data:     buffer.Bytes(),
	}, nil
}

func (d *S3Docstore) DeleteDocument(ctx context.Context, customer *queries.Customer, filename string) error {
	d.logger.InfoContext(ctx, "Deleting file ...", "filename", filename)

	client := s3.NewFromConfig(*d.config)

	fileId := createUniqueFileId(customer, filename)
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return fmt.Errorf("there was an issue deleting the object: %v", err)
	}

	d.logger.InfoContext(ctx, "Successfully deleted file")

	return nil
}
