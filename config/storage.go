package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Client provides a convenient interface for interacting with R2
type S3Client struct {
	Client *s3.S3
}

// NewS3Client creates a new S3Client instance
func NewS3Client() (*S3Client, error) {
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
