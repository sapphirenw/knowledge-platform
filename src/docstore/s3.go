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

	// internal state to avoid making multiple copies
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
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

	// create an s3 client
	client := s3.NewFromConfig(config)

	return &S3Docstore{
		bucket: S3_BUCKET,
		logger: lgr,
		config: &config,
		client: client,
	}, nil
}

func (d *S3Docstore) UploadDocument(ctx context.Context, customer *queries.Customer, doc *Doc) (string, error) {
	// create the s3 client
	if d.uploader == nil {
		d.uploader = manager.NewUploader(d.client)
	}

	// create a unique id
	documentId := createUniqueFileId(customer, doc.Filename)
	d.logger.InfoContext(ctx, "Uploading file", "filename", doc.Filename, "filetype", doc.Filetype)

	// upload with a docid
	result, err := d.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(documentId),
		Body:   bytes.NewReader(doc.Data),
	})
	if err != nil {
		return "", fmt.Errorf("there was an issue uploading the file: %v", err)
	}

	d.logger.InfoContext(ctx, "Successfully uploaded file", "filename", doc.Filename, "filetype", doc.Filetype)

	return result.Location, nil
}

func (d *S3Docstore) GetDocument(ctx context.Context, customer *queries.Customer, filename string) (*Doc, error) {
	d.logger.InfoContext(ctx, "Donwloading file from s3...", "filename", filename)

	if d.downloader == nil {
		d.downloader = manager.NewDownloader(d.client)
	}

	buffer := manager.NewWriteAtBuffer([]byte{})

	// upload with a unqiue identifier
	fileId := createUniqueFileId(customer, filename)
	_, err := d.downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
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

	fileId := createUniqueFileId(customer, filename)
	_, err := d.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return fmt.Errorf("there was an issue deleting the object: %v", err)
	}

	d.logger.InfoContext(ctx, "Successfully deleted file")

	return nil
}
