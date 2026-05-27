-- Create new customer
-- name: CreateCustomer :one
INSERT INTO users (phone)
VALUES ($1)
RETURNING id;

-- Get customer by ID
-- name: FindCustomerByID :one
SELECT * FROM users WHERE id = $1;

-- Get customer by phone
-- name: FindCustomerByPhone :one
SELECT * FROM users WHERE phone = $1;

-- Get customer profile with stats
-- name: GetCustomerProfile :one
SELECT
    u.id,
    u.phone,
    u.name,
    u.created_at,
    u.updated_at,
    COALESCE(order_stats.orders_count, 0) AS orders_count,
    COALESCE(review_stats.reviews_count, 0) AS reviews_count,
    COALESCE(order_stats.saved_amount, 0)::numeric AS saved_amount
FROM users u
LEFT JOIN (
    SELECT
        o.user_id,
        COUNT(*) FILTER (WHERE o.status NOT IN ('cancelled', 'refunded', 'disputed')) AS orders_count,
        COALESCE(SUM(
            CASE
                WHEN o.status IN ('paid', 'confirmed', 'completed', 'picked_up')
                    THEN COALESCE(sb.original_price, o.amount) - o.amount
                ELSE 0
            END
        ), 0) AS saved_amount
    FROM orders o
    LEFT JOIN surprise_boxes sb ON sb.id = o.box_id
    GROUP BY o.user_id
) AS order_stats ON order_stats.user_id = u.id
LEFT JOIN (
    SELECT
        r.user_id,
        COUNT(*) AS reviews_count
    FROM reviews r
    GROUP BY r.user_id
) AS review_stats ON review_stats.user_id = u.id
WHERE u.id = $1;

-- Update customer profile
-- name: UpdateCustomerProfile :one
UPDATE users
SET name = $2
WHERE id = $1
RETURNING *;
