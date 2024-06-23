package docstore

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	downloader *manager.Downloader
}

func NewS3Docstore(ctx context.Context, bucket string, logger *slog.Logger) (*S3Docstore, error) {
	if logger == nil {
		logger = utils.DefaultLogger()
	}
	l := logger.With("docstore", "s3", "s3_bucket", bucket)

	// TODO -- make client configurations global instead of creaing new everytime

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

func (d *S3Docstore) GeneratePresignedUrl(
	ctx context.Context,
	doc *queries.Document,
	contentType string,
	remoteId string,
) (string, error) {
	l := d.logger.With("doc", doc.Filename)
	l.InfoContext(ctx, "Generating a presigned url ...")
	// Set the desired parameters for the pre-signed URL
	presignClient := s3.NewPresignClient(d.client)
	params := &s3.PutObjectInput{
		Bucket:         aws.String(d.bucket),
		Key:            aws.String(remoteId),
		ContentType:    aws.String(contentType),
		ChecksumSHA256: aws.String(doc.Sha256),
	}

	resp, err := presignClient.PresignPutObject(ctx, params, func(o *s3.PresignOptions) {
		o.Expires = time.Minute * 10
	})
	if err != nil {
		return "", fmt.Errorf("there was an issue generating the pre-signed url: %v", err)
	}

	l.InfoContext(ctx, "Successfully generated pre-signed url")

	// Return the pre-signed URL
	return resp.URL, nil
}

func (d *S3Docstore) DownloadFile(ctx context.Context, uniqueId string) ([]byte, error) {
	l := d.logger.With("uniqueId", uniqueId)
	l.InfoContext(ctx, "Downloading the file from the remote datastore ...")

	if d.downloader == nil {
		d.downloader = manager.NewDownloader(d.client)
	}

	buffer := manager.NewWriteAtBuffer([]byte{})

	// upload with a unqiue identifier
	_, err := d.downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(uniqueId),
	})
	if err != nil {
		return nil, fmt.Errorf("there was an issue downloading the file from s3: %v | uniqueId=%s", err, uniqueId)
	}

	l.InfoContext(ctx, "Successfully downloaded file")
	return buffer.Bytes(), nil
}

func (d *S3Docstore) DeleteFile(ctx context.Context, uniqueId string) error {
	l := d.logger.With("uniqueId", uniqueId)
	l.InfoContext(ctx, "Deleting the file from the remote docstore ...")

	_, err := d.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(uniqueId),
	})
	if err != nil {
		return fmt.Errorf("there was an issue deleting the object: %v", err)
	}

	l.InfoContext(ctx, "Successfully deleted file")

	return nil
}

func (d *S3Docstore) DeleteRoot(ctx context.Context, prefix string) error {
	l := d.logger.With("prefix", prefix)
	l.InfoContext(ctx, "Deleting all keys under this prefix ...")

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(d.bucket),
		Prefix: aws.String(prefix),
	}

	// Iterate through the list of objects
	objects, err := d.client.ListObjectsV2(context.TODO(), listInput)
	if err != nil {
		return fmt.Errorf("error listing all objects: %s", err)
	}

	l.InfoContext(ctx, "Successfully found all objects", "length", len(objects.Contents))

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

	l.InfoContext(ctx, "Successfully deleted root folder")

	return nil
}

func (d *S3Docstore) GetUploadMethod() string {
	return "PUT"
}
