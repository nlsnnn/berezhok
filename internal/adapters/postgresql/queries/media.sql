-- Create a new media file record
-- name: CreateMediaFile :one
INSERT INTO media_files (
    filename, original_filename, storage_key, url, content_type, size_bytes
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- Get media file by ID
-- name: FindMediaFileByID :one
SELECT * FROM media_files WHERE id = $1;

-- Delete media file record
-- name: DeleteMediaFile :exec
DELETE FROM media_files WHERE id = $1;

-- List all media files (paginated)
-- name: ListMediaFiles :many
SELECT * FROM media_files 
ORDER BY uploaded_at DESC
LIMIT $1 OFFSET $2;
