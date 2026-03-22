package errors

import "errors"

var (
	ErrNoFileProvided    = errors.New("no file provided")
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrFileTooLarge      = errors.New("file size exceeds limit")
	ErrUploadFailed      = errors.New("failed to upload file")
	ErrInvalidMediaID    = errors.New("invalid media id")
	ErrMediaNotFound     = errors.New("media file not found")
	ErrDeleteFailed      = errors.New("failed to delete file")
)
