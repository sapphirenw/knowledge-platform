package docstore

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sapphirenw/ai-content-creation-api/src/document"
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
	l := logger.With("docstore", "s3", "s3_bucket", bucket)

	// setup aws config using env
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}

	// create an s3 client
	client := s3.NewFromConfig(config)

	return &S3Docstore{
		bucket: S3_BUCKET,
		logger: l,
		config: &config,
		client: client,
	}, nil
}

func (d *S3Docstore) UploadDocument(ctx context.Context, customer *queries.Customer, doc *document.Doc) (string, error) {
	// create the s3 client
	if d.uploader == nil {
		d.uploader = manager.NewUploader(d.client)
	}

	// create a unique id
	documentId := createUniqueFileId(customer, doc.Filename, doc.ParentId)
	d.logger.InfoContext(ctx, "Uploading file", "filename", doc.Filename, "filetype", doc.Filetype)

	// upload with a docid
	result, err := d.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(documentId),
		Body:   bytes.NewBuffer(doc.Data),
	})
	if err != nil {
		return "", fmt.Errorf("there was an issue uploading the file: %v", err)
	}

	d.logger.InfoContext(ctx, "Successfully uploaded file", "filename", doc.Filename, "filetype", doc.Filetype)

	return result.Location, nil
}

func (d *S3Docstore) GetDocument(ctx context.Context, customer *queries.Customer, parentId *int64, filename string) (*document.Doc, error) {
	d.logger.InfoContext(ctx, "Donwloading file from s3...", "parentId", parentId, "filename", filename)

	if d.downloader == nil {
		d.downloader = manager.NewDownloader(d.client)
	}

	buffer := manager.NewWriteAtBuffer([]byte{})

	// upload with a unqiue identifier
	fileId := createUniqueFileId(customer, filename, parentId)
	_, err := d.downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return nil, fmt.Errorf("there was an issue downloading the file from s3: %v | fileId=%s", err, fileId)
	}

	doc, err := document.NewDoc(parentId, parseUniqueFileId(fileId), buffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error decoding document: %s", err)
	}

	d.logger.InfoContext(ctx, "Successfully downloaded file", "parentId", parentId, "filename", filename)
	return doc, nil
}

func (d *S3Docstore) DeleteDocument(ctx context.Context, customer *queries.Customer, parentId *int64, filename string) error {
	d.logger.InfoContext(ctx, "Deleting file ...", "filename", filename)

	fileId := createUniqueFileId(customer, filename, parentId)
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

func (d *S3Docstore) DeleteRoot(ctx context.Context, customer *queries.Customer) error {
	d.logger.InfoContext(ctx, "Deleting all documents for the customer", "customer.ID", customer.ID)

	prefix := fmt.Sprintf("%d/", customer.ID)
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(d.bucket),
		Prefix: aws.String(prefix),
	}

	// Iterate through the list of objects
	objects, err := d.client.ListObjectsV2(context.TODO(), listInput)
	if err != nil {
		return fmt.Errorf("error listing all objects: %s", err)
	}

	d.logger.InfoContext(ctx, "Successfully found all objects", "length", len(objects.Contents))

	for _, object := range objects.Contents {
		// Delete each object
		delInput := &s3.DeleteObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    object.Key,
		}
		_, err := d.client.DeleteObject(context.TODO(), delInput)
		if err != nil {
			return fmt.Errorf("failed to delete object: %s", err)
		}
	}

	d.logger.InfoContext(ctx, "Successfully deleted root folder")

	return nil
}

func (d *S3Docstore) GeneratePresignedUrl(
	ctx context.Context,
	customer *queries.Customer,
	input *UploadUrlInput,
) (string, error) {
	documentId := createUniqueFileId(customer, input.Filename, input.ParentId)
	// Set the desired parameters for the pre-signed URL
	presignClient := s3.NewPresignClient(d.client)
	params := &s3.PutObjectInput{
		Bucket:         aws.String(d.bucket),
		Key:            &documentId,
		ContentType:    aws.String(input.Mime),
		ChecksumSHA256: aws.String(input.Signature),
	}

	resp, err := presignClient.PresignPutObject(ctx, params, func(o *s3.PresignOptions) {
		o.Expires = time.Minute * 10
	})
	if err != nil {
		return "", fmt.Errorf("there was an issue generating the pre-signed url: %v", err)
	}

	// Return the pre-signed URL
	return resp.URL, nil
}

func (d *S3Docstore) GetUploadMethod() string {
	return "PUT"
}
