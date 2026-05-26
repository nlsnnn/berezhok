-- Location queries for customer app

-- Search locations
-- name: SearchLocations :many
SELECT
    l.id,
    l.partner_id,
    l.name,
    l.address,
    l.phone,
    l.logo_url,
    l.cover_image_url,
    l.gallery_urls,
    l.working_hours,
    l.status,
    l.category_code,
    lc.name_ru as category_name,
    lc.icon_url as category_icon_url,
    lc.color as category_color,
    ST_X(l.location::geometry) as longitude,
    ST_Y(l.location::geometry) as latitude,
    COALESCE((SELECT COUNT(*) FROM surprise_boxes sb WHERE sb.location_id = l.id AND sb.status = 'active' AND sb.quantity_available > 0 AND sb.pickup_time_end > CURRENT_TIME), 0)::int as active_boxes_count,
    COALESCE(
        (SELECT json_agg(json_build_object('code', lp.code, 'name_ru', lp.name_ru) ORDER BY lp.sort_order)
         FROM location_selected_pins lsp
         JOIN location_pins lp ON lp.code = lsp.pin_code
         WHERE lsp.location_id = l.id),
        '[]'::json
    ) AS pins,
    l.created_at,
    l.updated_at
FROM locations l
JOIN location_categories lc ON l.category_code = lc.code
WHERE l.status = 'active'
    AND (sqlc.narg('category_code')::varchar IS NULL OR l.category_code = sqlc.narg('category_code')::varchar)
ORDER BY l.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- Count active locations for pagination
-- name: CountActiveLocations :one
SELECT COUNT(*)
FROM locations
WHERE status = 'active'
    AND (sqlc.narg('category_code')::varchar IS NULL OR category_code = sqlc.narg('category_code')::varchar);

-- Get location details by ID with category info
-- name: GetLocationDetailsByID :one
SELECT
    l.id,
    l.partner_id,
    l.name,
    l.address,
    l.phone,
    l.logo_url,
    l.cover_image_url,
    l.gallery_urls,
    l.working_hours,
    l.status,
    l.category_code,
    lc.name_ru as category_name,
    lc.icon_url as category_icon_url,
    lc.color as category_color,
    ST_X(l.location::geometry) as longitude,
    ST_Y(l.location::geometry) as latitude,
    COALESCE(
        (SELECT json_agg(json_build_object('code', lp.code, 'name_ru', lp.name_ru) ORDER BY lp.sort_order)
         FROM location_selected_pins lsp
         JOIN location_pins lp ON lp.code = lsp.pin_code
         WHERE lsp.location_id = l.id),
        '[]'::json
    ) AS pins,
    l.created_at,
    l.updated_at
FROM locations l
JOIN location_categories lc ON l.category_code = lc.code
WHERE l.id = $1;

-- Count active boxes by location ID
-- name: CountActiveBoxesByLocationID :one
SELECT COUNT(*)
FROM surprise_boxes
WHERE location_id = $1
    AND status = 'active'
    AND quantity_available > 0
    AND pickup_time_end > CURRENT_TIME;
