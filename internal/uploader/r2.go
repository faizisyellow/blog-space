package uploader

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
	S3Client *s3.Client
	Bucket   string
}

func NewR2Client(bucketName, accountId, accessKeyId, accessKeySecret string) *R2Client {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return &R2Client{
		S3Client: client,
		Bucket:   bucketName,
	}
}

func (r2c *R2Client) UploadFile(ctx context.Context, key string, file io.ReadCloser, contentType string, errUpload chan<- error) {

	defer close(errUpload)

	select {
	case <-ctx.Done():
		errUpload <- ctx.Err()
		return
	default:
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(r2c.Bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	_, err := r2c.S3Client.PutObject(ctx, input)
	errUpload <- err
}
