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
    AND quantity_available > 0;
