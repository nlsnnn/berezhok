-- name: FindActiveOrdersByLocationId :many
SELECT o.* FROM orders o
JOIN surprise_boxes sb ON o.box_id = sb.id
WHERE o.location_id = $1 AND o.status IN ('pending', 'paid', 'confirmed') AND sb.quantity_available > 0;

-- name: CreateOrder :one
INSERT INTO orders (
    user_id, box_id, location_id, pickup_code, qr_code_url, amount,
    pickup_time_start, pickup_time_end, status, partner_confirmation_deadline
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrdersByCustomerID :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: ReserveBox :execrows
UPDATE surprise_boxes 
SET quantity_available = quantity_available - 1
WHERE id = $1 AND quantity_available > 0 AND status = 'active';
