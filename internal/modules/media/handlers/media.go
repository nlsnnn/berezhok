package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/modules/media/domain"
	mediaErrors "github.com/nlsnnn/berezhok/internal/modules/media/errors"
	"github.com/nlsnnn/berezhok/internal/modules/media/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/modules/media/service"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

const maxUploadSize = 10 * 1024 * 1024 // 10MB

type mediaHandler struct {
	mediaService MediaService
	log          *slog.Logger
}

type MediaService interface {
	UploadFile(ctx context.Context, input service.UploadFileInput) (*domain.MediaFile, error)
}

func NewMediaHandler(mediaService MediaService, log *slog.Logger) *mediaHandler {
	return &mediaHandler{
		mediaService: mediaService,
		log:          log,
	}
}

func (h *mediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	const op = "media.handler.Upload"
	log := h.log.With(slog.String("op", op))

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Warn("failed to parse multipart form", sl.Err(err))
		response.BadRequest(w, "file size exceeds 10MB limit")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Warn("no file provided", sl.Err(err))
		response.BadRequest(w, "no file provided")
		return
	}
	defer func() { _ = file.Close() }()

	mediaFile, err := h.mediaService.UploadFile(r.Context(), service.UploadFileInput{
		File:        file,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
	})
	if err != nil {
		switch {
		case errors.Is(err, mediaErrors.ErrNoFileProvided):
			response.BadRequest(w, "no file provided")
		case errors.Is(err, mediaErrors.ErrInvalidFileType):
			response.BadRequest(w, "invalid file type, allowed: jpg, png, webp, pdf")
		case errors.Is(err, mediaErrors.ErrFileTooLarge):
			response.BadRequest(w, "file size exceeds 10MB limit")
		default:
			log.Error("failed to upload file", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Created(w, dto.UploadResponse{
		URL:      mediaFile.URL,
		Filename: mediaFile.OriginalFilename,
		Size:     mediaFile.SizeBytes,
	})
}
