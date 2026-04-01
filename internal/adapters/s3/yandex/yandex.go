package yandex

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	sharedConfig "github.com/nlsnnn/berezhok/internal/shared/config"
)

type Storage struct {
	client        *s3.Client
	bucket        string
	publicBaseURL string
}

func NewStorage(cfg sharedConfig.S3) (*Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = false
	})

	return &Storage{
		client:        client,
		bucket:        cfg.Bucket,
		publicBaseURL: cfg.PublicBaseURL,
	}, nil
}

func (s *Storage) UploadFile(ctx context.Context, file io.Reader, filename, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ext := filepath.Ext(filename)
	storageKey := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	objInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(storageKey),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, objInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return storageKey, nil
}

func (s *Storage) DeleteFile(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func (s *Storage) GetPublicURL(key string) string {
	return fmt.Sprintf("%s/%s", s.publicBaseURL, key)
}
