-- name: ListApplications :many
SELECT * FROM partner_applications;

-- name: FindApplicationByID :one
SELECT * FROM partner_applications WHERE id = $1;

-- name: CreateApplication :one
INSERT INTO partner_applications (
    contact_name, contact_email, contact_phone, business_name, category_code, 
    address, description, status
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateApplication :exec
UPDATE partner_applications
SET status = $1, reviewed_at = $2, rejection_reason = $3
WHERE id = $4;

-- name: DeleteApplication :exec
DELETE FROM partner_applications WHERE id = $1;