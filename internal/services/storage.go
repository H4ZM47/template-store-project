package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appconfig "template-store/internal/config"
)

// StorageService defines the interface for file storage operations.
type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

// S3StorageService provides file storage services using AWS S3.
type S3StorageService struct {
	client     *s3.Client
	bucketName string
	region     string
}

// NewStorageService creates a new S3StorageService.
func NewStorageService(cfg *appconfig.Config) (*S3StorageService, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.AWS.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3StorageService{
		client:     s3.NewFromConfig(awsCfg),
		bucketName: cfg.AWS.S3Bucket,
		region:     cfg.AWS.Region,
	}, nil
}

// UploadFile uploads a file to S3 and returns its URL.
func (s *S3StorageService) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(file.Filename),
		Body:   src,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct the URL of the uploaded object
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, file.Filename)

	return url, nil
}
