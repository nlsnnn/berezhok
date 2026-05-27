-- name: GetEcoAggregateByUserID :many
-- Aggregates completed orders for a customer grouped by location category.
-- "Completed" = any order that reached terminal status 'completed', regardless
-- of how it got there (partner scan, auto-complete, admin override, etc.).
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
  AND o.status = 'completed'
GROUP BY l.category_code;
