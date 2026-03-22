package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/media/domain"
	mediaErrors "github.com/nlsnnn/berezhok/internal/modules/media/errors"
)

type MediaRepo struct {
	q *sqlc.Queries
}

func NewMediaRepo(q *sqlc.Queries) *MediaRepo {
	return &MediaRepo{q: q}
}

func (r *MediaRepo) CreateMediaFile(ctx context.Context, media *domain.MediaFile) error {
	m, err := r.q.CreateMediaFile(ctx, sqlc.CreateMediaFileParams{
		Filename:         media.Filename,
		OriginalFilename: media.OriginalFilename,
		StorageKey:       media.StorageKey,
		Url:              media.URL,
		ContentType:      media.ContentType,
		SizeBytes:        media.SizeBytes,
	})
	if err != nil {
		return err
	}

	media.ID = m.ID
	media.UploadedAt = m.UploadedAt.Time
	return nil
}

func (r *MediaRepo) GetMediaFileByID(ctx context.Context, id string) (*domain.MediaFile, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, mediaErrors.ErrInvalidMediaID
	}

	m, err := r.q.FindMediaFileByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, mediaErrors.ErrMediaNotFound
		}
		return nil, err
	}

	return &domain.MediaFile{
		ID:               m.ID,
		Filename:         m.Filename,
		OriginalFilename: m.OriginalFilename,
		StorageKey:       m.StorageKey,
		URL:              m.Url,
		ContentType:      m.ContentType,
		SizeBytes:        m.SizeBytes,
		UploadedAt:       m.UploadedAt.Time,
	}, nil
}

func (r *MediaRepo) DeleteMediaFile(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return mediaErrors.ErrInvalidMediaID
	}

	_, err = r.q.FindMediaFileByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return mediaErrors.ErrMediaNotFound
		}
		return err
	}

	return r.q.DeleteMediaFile(ctx, uid)
}
