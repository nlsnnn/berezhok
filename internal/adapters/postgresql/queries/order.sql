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

-- name: GetOrderDetailsByID :one
SELECT
    o.id,
    o.user_id,
    o.status,
    o.pickup_code,
    COALESCE(o.qr_code_url, '') AS qr_code_url,
    o.amount,
    o.pickup_time_start,
    o.pickup_time_end,
    o.created_at,
    o.partner_confirmed_at,
    sb.name AS box_name,
    COALESCE(sb.image_url, '') AS box_image_url,
    l.name AS location_name,
    l.address AS location_address,
    COALESCE(l.phone, '') AS location_phone,
    ST_Y(l.location::geometry) AS location_lat,
    ST_X(l.location::geometry) AS location_lng
FROM orders o
JOIN surprise_boxes sb ON sb.id = o.box_id
JOIN locations l ON l.id = o.location_id
WHERE o.id = $1;

-- name: GetPartnerOrderByPickupCode :one
SELECT
    o.id,
    o.pickup_code,
    o.status,
    sb.name AS box_name,
    COALESCE(sb.image_url, '') AS box_image_url,
    u.phone AS customer_phone,
    COALESCE(u.name, '') AS customer_name,
    o.pickup_time_start,
    o.pickup_time_end,
    o.created_at
FROM orders o
JOIN surprise_boxes sb ON sb.id = o.box_id
JOIN users u ON u.id = o.user_id
JOIN locations l ON l.id = o.location_id
WHERE o.pickup_code = $1
  AND l.partner_id = $2
ORDER BY o.created_at DESC
LIMIT 1;

-- name: GetPartnerOrderByID :one
SELECT
    o.id,
    o.status
FROM orders o
JOIN locations l ON l.id = o.location_id
WHERE o.id = $1
  AND l.partner_id = $2;

-- name: MarkOrderPickedUp :execrows
UPDATE orders
SET status = 'completed',
    picked_up_at = NOW(),
    picked_up_confirmed_by = $2,
    updated_at = NOW()
WHERE id = $1
  AND status = 'confirmed';

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

-- name: ListOrdersByCustomerIDFiltered :many
SELECT
    o.id, o.status, o.pickup_code, o.amount,
    o.pickup_time_start, o.created_at,
    sb.name AS box_name,
    l.name AS location_name,
    false AS has_review
FROM orders o
JOIN surprise_boxes sb ON o.box_id = sb.id
JOIN locations l ON o.location_id = l.id
WHERE o.user_id = $1
  AND ($2 = '' OR o.status::text = $2)
ORDER BY o.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountOrdersByCustomerID :one
SELECT COUNT(*)
FROM orders o
WHERE o.user_id = $1
  AND ($2 = '' OR o.status::text = $2);
