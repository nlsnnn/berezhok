package domain

import (
	"time"

	"github.com/google/uuid"
)

type MediaFile struct {
	ID               uuid.UUID
	Filename         string
	OriginalFilename string
	StorageKey       string
	URL              string
	ContentType      string
	SizeBytes        int64
	UploadedAt       time.Time
}
