-- name: CreateReview :one
INSERT INTO reviews (
        order_id,
        user_id,
        location_id,
        rating,
        comment
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListLocationReviews :many
SELECT r.id,
    r.rating,
    COALESCE(r.comment, '') AS comment,
    COALESCE(u.name, '') AS user_name,
    r.created_at
FROM reviews r
    JOIN users u ON u.id = r.user_id
WHERE r.location_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;
-- name: CountLocationReviews :one
SELECT COUNT(*)
FROM reviews
WHERE location_id = $1;
