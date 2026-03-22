package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"path/filepath"
	"strings"

	"github.com/nlsnnn/berezhok/internal/modules/media/domain"
	mediaErrors "github.com/nlsnnn/berezhok/internal/modules/media/errors"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
)

var allowedContentTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

type UploadFileInput struct {
	File        io.Reader
	Filename    string
	ContentType string
	Size        int64
}

type mediaService struct {
	storage  Storage
	mediaRepo MediaRepository
	log      *slog.Logger
}

type Storage interface {
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	DeleteFile(ctx context.Context, key string) error
	GetPublicURL(key string) string
}

type MediaRepository interface {
	CreateMediaFile(ctx context.Context, media *domain.MediaFile) error
	GetMediaFileByID(ctx context.Context, id string) (*domain.MediaFile, error)
	DeleteMediaFile(ctx context.Context, id string) error
}

func NewMediaService(storage Storage, mediaRepo MediaRepository, log *slog.Logger) *mediaService {
	return &mediaService{
		storage:  storage,
		mediaRepo: mediaRepo,
		log:      log,
	}
}

func (s *mediaService) UploadFile(ctx context.Context, input UploadFileInput) (*domain.MediaFile, error) {
	if input.File == nil {
		return nil, mediaErrors.ErrNoFileProvided
	}

	if input.Size > maxFileSize {
		return nil, mediaErrors.ErrFileTooLarge
	}

	contentType := s.detectContentType(input.ContentType, input.Filename)
	if !allowedContentTypes[contentType] {
		return nil, mediaErrors.ErrInvalidFileType
	}

	storageKey, err := s.storage.UploadFile(ctx, input.File, input.Filename, contentType)
	if err != nil {
		s.log.Error("failed to upload file to storage", "error", err)
		return nil, fmt.Errorf("%w: %v", mediaErrors.ErrUploadFailed, err)
	}

	publicURL := s.storage.GetPublicURL(storageKey)

	mediaFile := &domain.MediaFile{
		Filename:         filepath.Base(storageKey),
		OriginalFilename: input.Filename,
		StorageKey:       storageKey,
		URL:              publicURL,
		ContentType:      contentType,
		SizeBytes:        input.Size,
	}

	if err := s.mediaRepo.CreateMediaFile(ctx, mediaFile); err != nil {
		s.log.Error("failed to save media file to database", "error", err)
		if delErr := s.storage.DeleteFile(ctx, storageKey); delErr != nil {
			s.log.Error("failed to cleanup uploaded file after db error", "error", delErr)
		}
		return nil, fmt.Errorf("failed to save media file: %w", err)
	}

	return mediaFile, nil
}

func (s *mediaService) detectContentType(providedType, filename string) string {
	if providedType != "" && providedType != "application/octet-stream" {
		return providedType
	}

	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return strings.Split(mimeType, ";")[0]
	}

	return "application/octet-stream"
}
