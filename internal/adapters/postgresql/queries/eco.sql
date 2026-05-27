-- name: GetEcoAggregateByUserID :many
-- Aggregates picked-up orders for a customer grouped by location category.
-- "Picked up" = orders that physically reached the customer (picked_up_at IS NOT NULL).
-- Returned per-category counts let the application layer multiply by the
-- average-weight coefficient that lives in Go code (see internal/modules/eco/coefficients.go).
SELECT
    l.category_code                                              AS category_code,
    COUNT(*)                                                     AS picked_count,
    COALESCE(SUM(sb.original_price - sb.discount_price), 0)::numeric AS savings
FROM orders o
    JOIN surprise_boxes sb ON sb.id = o.box_id
    JOIN locations l ON l.id = o.location_id
WHERE o.user_id = $1
  AND o.picked_up_at IS NOT NULL
GROUP BY l.category_code;
