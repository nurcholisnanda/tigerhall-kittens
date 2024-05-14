// pkg/imagehandler/image.go
package imagehandler

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg" // Import the necessary image formats
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nfnt/resize"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/storage"
)

// ImageHandler defines an interface for image-related operations
//
//go:generate mockgen -source=image.go -destination=mock/image.go -package=mock
type ImageHandler interface {
	ResizeAndUpload(ctx context.Context, inputImage *graphql.Upload) (string, error)
}

// imageService implements the ImageHandler interface
type imageService struct {
	s3Client storage.S3Interface
}

// NewImageService creates a new ImageService instance
func NewImageService(s3Client storage.S3Interface) ImageHandler {
	return &imageService{s3Client: s3Client}
}

// ResizeAndUpload resizes the image and uploads it to S3
func (s *imageService) ResizeAndUpload(ctx context.Context, inputImage *graphql.Upload) (string, error) {
	imageData, readErr := io.ReadAll(inputImage.File)
	if readErr != nil {
		logger.Logger(ctx).Error(ctx, "readErr", readErr)

		fmt.Printf("error from file %v", readErr)
	}

	logger := logger.Logger(ctx) // Get contextual logger

	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.Error(ctx, "Error decoding image", "error", err)
		return "", fmt.Errorf("error decoding image: %v", err)
	}
	logger.Info(ctx, "Image format:", format)

	resizedImage := resize.Resize(250, 200, img, resize.Lanczos3)

	// Encode resized image to JPEG format
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		logger.Error(ctx, "Error encoding image", "error", err)
		return "", fmt.Errorf("error encoding image: %v", err)
	}

	// Upload to R2
	objectURL, err := s.s3Client.PutObject(format, buf)
	if err != nil {
		logger.Error(ctx, "Error uploading image to R2", "error", err)
		return "", fmt.Errorf("error uploading image to R2: %w", err)
	}
	logger.Info(ctx, "Image uploaded to R2", "objectURL:", objectURL)

	return objectURL, nil
}
