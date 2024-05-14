package storage

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// S3Client provides a convenient interface for interacting with R2
type S3Client struct {
	Client *s3.S3
}

//go:generate mockgen -source=storage.go -destination=mock/storage.go -package=mock
type S3Interface interface {
	PutObject(format string, imageData *bytes.Buffer) (string, error)
}

// NewS3Client creates a new S3Client instance
func NewS3Client() (S3Interface, error) {
	// Retrieve AWS credentials (use environment variables or a secure configuration)
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("auto"), // R2 uses "auto" for the region
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String("https://" + accountID + ".r2.cloudflarestorage.com"),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %w", err)
	}

	// Create S3 client
	s3Client := s3.New(sess)
	return &S3Client{Client: s3Client}, nil
}

func (c *S3Client) PutObject(format string, imageData *bytes.Buffer) (string, error) {
	bucketName := os.Getenv("R2_BUCKET_NAME") // Get bucket name from env vars

	objectName := fmt.Sprintf("%s.%s", uuid.NewString(), format)

	if _, err := c.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   bytes.NewReader(imageData.Bytes()),
	}); err != nil {
		return "", err
	}
	objectURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", os.Getenv("CLOUDFLARE_ACCOUNT_ID"), bucketName, objectName)

	return objectURL, nil
}
