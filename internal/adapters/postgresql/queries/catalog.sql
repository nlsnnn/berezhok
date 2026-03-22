-- Create a new box
-- name: CreateBox :one
INSERT INTO surprise_boxes (
    name, description, original_price, discount_price, quantity_available,
    pickup_time_start, pickup_time_end, image_url, status, location_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- Update an existing box
-- name: UpdateBox :one
UPDATE surprise_boxes SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    original_price = COALESCE($4, original_price),
    discount_price = COALESCE($5, discount_price),
    quantity_available = COALESCE($6, quantity_available),
    pickup_time_start = COALESCE($7, pickup_time_start),
    pickup_time_end = COALESCE($8, pickup_time_end),
    image_url = COALESCE($9, image_url),
    status = COALESCE($10, status)
WHERE id = $1 RETURNING *;

-- Get a box by ID
-- name: FindBoxByID :one
SELECT * FROM surprise_boxes WHERE id = $1;

-- List boxes by location ID
-- name: ListBoxesByLocationID :many
SELECT * FROM surprise_boxes WHERE location_id = $1;

-- List active boxes by location ID
-- name: ListActiveBoxesByLocationID :many
SELECT * FROM surprise_boxes 
WHERE location_id = $1 AND status = 'active' AND quantity_available > 0;

-- List boxes by partner ID
-- name: ListBoxesByPartnerID :many
SELECT sb.* FROM surprise_boxes sb
JOIN locations l ON sb.location_id = l.id
WHERE l.partner_id = $1;

-- Delete a box
-- name: DeleteBox :exec
DELETE FROM surprise_boxes WHERE id = $1;

