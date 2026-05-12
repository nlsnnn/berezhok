-- Заявки на партнёрство
-- name: ListApplications :many
SELECT * FROM partner_applications;

-- name: FindApplicationByID :one
SELECT * FROM partner_applications WHERE id = $1;

-- name: CreateApplication :one
INSERT INTO partner_applications (
    contact_name, contact_email, contact_phone, business_name, category_code,
    address, description, status, latitude, longitude
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: UpdateApplication :exec
UPDATE partner_applications
SET status = $1, reviewed_at = $2, rejection_reason = $3
WHERE id = $4;

-- name: DeleteApplication :exec
DELETE FROM partner_applications WHERE id = $1;


-- Партнёры (юридические лица)
-- name: ListPartners :many
SELECT * FROM partners;

-- name: FindPartnerByID :one
SELECT * FROM partners WHERE id = $1;

-- name: UpdatePartnerStatusByID :exec
UPDATE partners
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpsertPartnerLegalInfo :one
INSERT INTO partner_legal_info (
    partner_id, inn, ogrn, kpp, legal_address, legal_name
)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (partner_id)
DO UPDATE SET
    inn = EXCLUDED.inn,
    ogrn = EXCLUDED.ogrn,
    kpp = EXCLUDED.kpp,
    legal_address = EXCLUDED.legal_address,
    legal_name = EXCLUDED.legal_name,
    updated_at = NOW()
RETURNING *;

-- name: CheckEmailExists :one
SELECT EXISTS (
    SELECT 1 FROM partner_employees WHERE email = $1
    UNION
    SELECT 1 FROM partner_applications WHERE contact_email = $1 AND status = 'pending'
) AS email_exists;


-- name: GetPartnerProfile :one
SELECT
    e.id as employee_id, e.name as employee_name, e.email, e.role, e.created_at as employee_created_at, e.must_change_password,
    p.id as partner_id, p.brand_name, p.status as partner_status,
    CASE
        WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
        ELSE p.commission_rate
    END AS commission_rate,
    p.promo_commission_until,
    p.created_at as partner_created_at,
    l.id as location_id, l.name as location_name, l.address as location_address, l.created_at as location_created_at
FROM partner_employees e
JOIN partners p ON e.partner_id = p.id
LEFT JOIN locations l ON e.location_id = l.id
WHERE e.id = $1;

-- name: GetPartnerDashboardTodayStats :one
SELECT
    COUNT(*) FILTER (WHERE o.status = 'paid')::bigint AS pending_confirmation,
    COUNT(*) FILTER (WHERE o.status = 'confirmed')::bigint AS confirmed,
    COUNT(*) FILTER (WHERE o.status = 'picked_up')::bigint AS picked_up,
    COUNT(*) FILTER (WHERE o.status = 'completed')::bigint AS completed
FROM orders o
JOIN locations l ON l.id = o.location_id
WHERE l.partner_id = $1
  AND o.created_at >= date_trunc('day', NOW())
  AND o.created_at < date_trunc('day', NOW()) + interval '1 day';

-- name: GetPartnerDashboardWeekStats :one
WITH commission AS (
    SELECT
        CASE
            WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
            ELSE p.commission_rate
        END AS rate
    FROM partners p
    WHERE p.id = $1
),
order_stats AS (
    SELECT
        COUNT(*) FILTER (WHERE o.status = 'completed')::bigint AS orders_completed,
        COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0)::bigint AS gross_revenue,
        COALESCE(SUM(o.amount * (1 - c.rate)) FILTER (WHERE o.status = 'completed'), 0)::bigint AS net_revenue
    FROM orders o
    JOIN locations l ON l.id = o.location_id
    CROSS JOIN commission c
    WHERE l.partner_id = $1
      AND o.created_at >= NOW() - interval '7 days'
),
review_stats AS (
    SELECT COALESCE(AVG(r.rating), 0)::double precision AS avg_rating
    FROM reviews r
    JOIN locations l ON l.id = r.location_id
    WHERE l.partner_id = $1
      AND r.created_at >= NOW() - interval '7 days'
)
SELECT
    order_stats.orders_completed,
    order_stats.gross_revenue,
    order_stats.net_revenue,
    review_stats.avg_rating
FROM order_stats, review_stats;

-- name: GetPartnerDashboardFinance :one
WITH commission AS (
    SELECT
        CASE
            WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
            ELSE p.commission_rate
        END AS rate
    FROM partners p
    WHERE p.id = $1
)
SELECT
    COALESCE(SUM(o.amount * (1 - c.rate)) FILTER (WHERE o.status = 'completed'), 0)::bigint AS balance_pending,
    (date_trunc('week', NOW()) + interval '1 week')::timestamptz AS next_payout_date
FROM orders o
JOIN locations l ON l.id = o.location_id
CROSS JOIN commission c
WHERE l.partner_id = $1;

-- name: GetPartnerStatsSummary :one
WITH commission AS (
    SELECT
        CASE
            WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
            ELSE p.commission_rate
        END AS rate
    FROM partners p
    WHERE p.id = sqlc.arg(partner_id)
),
order_stats AS (
    SELECT
        COUNT(*)::bigint AS orders_total,
        COUNT(*) FILTER (WHERE o.status = 'completed')::bigint AS orders_completed,
        COUNT(*) FILTER (WHERE o.status = 'cancelled')::bigint AS orders_cancelled,
        COUNT(*) FILTER (WHERE o.status = 'paid')::bigint AS orders_pending_confirmation,
        COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0)::bigint AS gross_revenue,
        COALESCE(SUM(o.amount * (1 - c.rate)) FILTER (WHERE o.status = 'completed'), 0)::bigint AS net_revenue,
        COALESCE(AVG(o.amount) FILTER (WHERE o.status = 'completed'), 0)::double precision AS avg_order_value
    FROM orders o
    JOIN locations l ON l.id = o.location_id
    CROSS JOIN commission c
    WHERE l.partner_id = sqlc.arg(partner_id)
      AND o.created_at >= sqlc.arg(start_date)
      AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
      AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
      AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
),
review_stats AS (
    SELECT
        COALESCE(AVG(r.rating), 0)::double precision AS avg_rating,
        COUNT(*)::bigint AS reviews_count
    FROM reviews r
    JOIN locations l ON l.id = r.location_id
    WHERE l.partner_id = sqlc.arg(partner_id)
      AND r.created_at >= sqlc.arg(start_date)
      AND r.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
      AND (sqlc.narg(location_id)::uuid IS NULL OR r.location_id = sqlc.narg(location_id)::uuid)
)
SELECT
    order_stats.orders_total,
    order_stats.orders_completed,
    order_stats.orders_cancelled,
    order_stats.orders_pending_confirmation,
    order_stats.gross_revenue,
    order_stats.net_revenue,
    order_stats.avg_order_value,
    review_stats.avg_rating,
    review_stats.reviews_count
FROM order_stats, review_stats;

-- name: GetPartnerStatsTimeline :many
WITH commission AS (
    SELECT
        CASE
            WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
            ELSE p.commission_rate
        END AS rate
    FROM partners p
    WHERE p.id = sqlc.arg(partner_id)
),
series AS (
    SELECT generate_series(sqlc.arg(start_date)::timestamptz, sqlc.arg(end_date)::timestamptz, interval '1 day') AS day
),
aggregated AS (
    SELECT
        date_trunc('day', o.created_at) AS day,
        COUNT(*)::bigint AS orders_total,
        COUNT(*) FILTER (WHERE o.status = 'completed')::bigint AS orders_completed,
        COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0)::bigint AS gross_revenue,
        COALESCE(SUM(o.amount * (1 - c.rate)) FILTER (WHERE o.status = 'completed'), 0)::bigint AS net_revenue
    FROM orders o
    JOIN locations l ON l.id = o.location_id
    CROSS JOIN commission c
    WHERE l.partner_id = sqlc.arg(partner_id)
      AND o.created_at >= sqlc.arg(start_date)
      AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
      AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
      AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
    GROUP BY 1
)
SELECT
    series.day::timestamptz AS date,
    COALESCE(aggregated.orders_total, 0)::bigint AS orders_total,
    COALESCE(aggregated.orders_completed, 0)::bigint AS orders_completed,
    COALESCE(aggregated.gross_revenue, 0)::bigint AS gross_revenue,
    COALESCE(aggregated.net_revenue, 0)::bigint AS net_revenue
FROM series
LEFT JOIN aggregated ON aggregated.day = series.day
ORDER BY series.day ASC;

-- name: GetPartnerStatsStatusBreakdown :many
WITH totals AS (
    SELECT COUNT(*)::double precision AS total_count
    FROM orders o
    JOIN locations l ON l.id = o.location_id
    WHERE l.partner_id = sqlc.arg(partner_id)
      AND o.created_at >= sqlc.arg(start_date)
      AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
      AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
      AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
)
SELECT
    o.status::text AS status,
    COUNT(*)::bigint AS count,
    CASE
        WHEN totals.total_count = 0 THEN 0::double precision
        ELSE COUNT(*)::double precision / totals.total_count
    END::double precision AS share
FROM orders o
JOIN locations l ON l.id = o.location_id
CROSS JOIN totals
WHERE l.partner_id = sqlc.arg(partner_id)
  AND o.created_at >= sqlc.arg(start_date)
  AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
  AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
  AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
GROUP BY o.status, totals.total_count
ORDER BY count DESC, o.status::text ASC;

-- name: GetPartnerStatsTopLocations :many
WITH commission AS (
    SELECT
        CASE
            WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
            ELSE p.commission_rate
        END AS rate
    FROM partners p
    WHERE p.id = sqlc.arg(partner_id)
),
review_stats AS (
    SELECT
        r.location_id,
        COALESCE(AVG(r.rating), 0)::double precision AS avg_rating
    FROM reviews r
    JOIN locations l ON l.id = r.location_id
    WHERE l.partner_id = sqlc.arg(partner_id)
      AND r.created_at >= sqlc.arg(start_date)
      AND r.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
      AND (sqlc.narg(location_id)::uuid IS NULL OR r.location_id = sqlc.narg(location_id)::uuid)
    GROUP BY r.location_id
)
SELECT
    l.id AS location_id,
    l.name,
    l.address,
    COUNT(o.id)::bigint AS orders_total,
    COUNT(o.id) FILTER (WHERE o.status = 'completed')::bigint AS orders_completed,
    COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0)::bigint AS gross_revenue,
    COALESCE(SUM(o.amount * (1 - c.rate)) FILTER (WHERE o.status = 'completed'), 0)::bigint AS net_revenue,
    COALESCE(review_stats.avg_rating, 0)::double precision AS avg_rating
FROM locations l
LEFT JOIN orders o ON o.location_id = l.id
    AND o.created_at >= sqlc.arg(start_date)
    AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
    AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
CROSS JOIN commission c
LEFT JOIN review_stats ON review_stats.location_id = l.id
WHERE l.partner_id = sqlc.arg(partner_id)
  AND (sqlc.narg(location_id)::uuid IS NULL OR l.id = sqlc.narg(location_id)::uuid)
GROUP BY l.id, l.name, l.address, review_stats.avg_rating
ORDER BY
    CASE WHEN sqlc.arg(sort)::text = 'revenue_desc' THEN COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0) END DESC,
    CASE WHEN sqlc.arg(sort)::text = 'orders_desc' THEN COUNT(o.id) END DESC,
    CASE WHEN sqlc.arg(sort)::text = 'rating_desc' THEN COALESCE(review_stats.avg_rating, 0) END DESC,
    CASE WHEN sqlc.arg(sort)::text = 'name_asc' THEN l.name END ASC,
    l.name ASC;

-- name: GetPartnerStatsTopBoxes :many
WITH commission AS (
    SELECT
        CASE
            WHEN p.promo_commission_until >= NOW() THEN COALESCE(p.promo_commission_rate, p.commission_rate)
            ELSE p.commission_rate
        END AS rate
    FROM partners p
    WHERE p.id = sqlc.arg(partner_id)
)
SELECT
    sb.id AS box_id,
    sb.name,
    COALESCE(sb.image_url, '') AS image_url,
    l.name AS location_name,
    COUNT(o.id)::bigint AS orders_total,
    COUNT(o.id) FILTER (WHERE o.status = 'completed')::bigint AS orders_completed,
    COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0)::bigint AS gross_revenue,
    COALESCE(SUM(o.amount * (1 - c.rate)) FILTER (WHERE o.status = 'completed'), 0)::bigint AS net_revenue
FROM surprise_boxes sb
JOIN locations l ON l.id = sb.location_id
LEFT JOIN orders o ON o.box_id = sb.id
    AND o.created_at >= sqlc.arg(start_date)
    AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
    AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
    AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
CROSS JOIN commission c
WHERE l.partner_id = sqlc.arg(partner_id)
  AND (sqlc.narg(location_id)::uuid IS NULL OR l.id = sqlc.narg(location_id)::uuid)
GROUP BY sb.id, sb.name, sb.image_url, l.name
ORDER BY
    CASE WHEN sqlc.arg(sort)::text = 'revenue_desc' THEN COALESCE(SUM(o.amount) FILTER (WHERE o.status = 'completed'), 0) END DESC,
    CASE WHEN sqlc.arg(sort)::text = 'orders_desc' THEN COUNT(o.id) END DESC,
    CASE WHEN sqlc.arg(sort)::text = 'name_asc' THEN sb.name END ASC,
    sb.name ASC;

-- name: ListPartnerStatsOrders :many
SELECT
  o.id,
  o.status,
  o.pickup_code,
  o.amount,
  o.created_at,
  sb.name AS box_name,
  COALESCE(sb.image_url, '') AS box_image_url,
  COALESCE(u.phone, '') AS customer_phone,
  COALESCE(u.name, '') AS customer_name,
  l.id AS location_id,
  l.name AS location_name,
  l.address AS location_address,
  o.pickup_time_start,
  o.pickup_time_end
FROM orders o
JOIN surprise_boxes sb ON o.box_id = sb.id
JOIN users u ON o.user_id = u.id
JOIN locations l ON o.location_id = l.id
WHERE l.partner_id = sqlc.arg(partner_id)
  AND o.created_at >= sqlc.arg(start_date)
  AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
  AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
  AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text)
ORDER BY
  CASE WHEN sqlc.arg(sort)::text = 'created_at_desc' THEN o.created_at END DESC,
  CASE WHEN sqlc.arg(sort)::text = 'created_at_asc' THEN o.created_at END ASC,
  CASE WHEN sqlc.arg(sort)::text = 'amount_desc' THEN o.amount END DESC,
  CASE WHEN sqlc.arg(sort)::text = 'amount_asc' THEN o.amount END ASC,
  CASE WHEN sqlc.arg(sort)::text = 'pickup_time_desc' THEN o.pickup_time_start END DESC,
  CASE WHEN sqlc.arg(sort)::text = 'pickup_time_asc' THEN o.pickup_time_start END ASC,
  o.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountPartnerStatsOrders :one
SELECT COUNT(*)
FROM orders o
JOIN locations l ON o.location_id = l.id
WHERE l.partner_id = sqlc.arg(partner_id)
  AND o.created_at >= sqlc.arg(start_date)
  AND o.created_at < sqlc.arg(end_date)::timestamptz + interval '1 day'
  AND (sqlc.narg(location_id)::uuid IS NULL OR o.location_id = sqlc.narg(location_id)::uuid)
  AND (sqlc.arg(status)::text = '' OR o.status::text = sqlc.arg(status)::text);

-- name: CreatePartner :one
INSERT INTO partners (
    brand_name, logo_url, parent_partner_id, account_type,
    commission_rate, promo_commission_rate, promo_commission_until, status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdatePartner :exec
UPDATE partners
SET brand_name = $1, logo_url = $2, parent_partner_id = $3,
    account_type = $4, commission_rate = $5, promo_commission_rate = $6,
    promo_commission_until = $7, status = $8
WHERE id = $9;


-- Сотрудники партнёров
-- name: ListPartnerEmployees :many
SELECT * FROM partner_employees;

-- name: ListEmployeesByPartnerID :many
SELECT * FROM partner_employees WHERE partner_id = $1;

-- name: FindPartnerEmployeeByID :one
SELECT * FROM partner_employees WHERE id = $1;

-- name: FindPartnerEmployeeByEmail :one
SELECT * FROM partner_employees WHERE email = $1;

-- name: CreatePartnerEmployee :one
INSERT INTO partner_employees (
    partner_id, location_id, email, password_hash, role, name
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdatePartnerEmployee :exec
UPDATE partner_employees
SET partner_id = $1, location_id = $2, email = $3, password_hash = $4,
    role = $5, name = $6
WHERE id = $7;

-- name: UpdatePartnerEmployeePassword :exec
UPDATE partner_employees
SET password_hash = $1, must_change_password = $2
WHERE id = $3;

-- name: DeletePartnerEmployee :exec
DELETE FROM partner_employees WHERE id = $1;


-- Локации партнёров
-- name: ListLocations :many
SELECT * FROM locations;

-- name: FindLocationByID :one
SELECT * FROM locations WHERE id = $1;

-- name: FindLocationsByPartnerID :many
SELECT * FROM locations WHERE partner_id = $1;

-- name: CreateLocation :one
INSERT INTO locations (name, address, partner_id, category_code, status, location)
VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326))
RETURNING *;

-- name: UpdateLocation :one
UPDATE locations
SET
    name = COALESCE($2, name),
    address = COALESCE($3, address),
    category_code = COALESCE($4, category_code),
    logo_url = COALESCE($5, logo_url),
    cover_image_url = COALESCE($6, cover_image_url),
    working_hours = COALESCE($7, working_hours),
    gallery_urls = COALESCE($8, gallery_urls),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ActivateLocation :exec
UPDATE locations SET status = 'active', updated_at = NOW() WHERE id = $1;

-- name: DeactivateLocation :exec
UPDATE locations SET status = 'inactive', updated_at = NOW() WHERE id = $1;

-- name: CloseLocation :exec
UPDATE locations SET status = 'closed', updated_at = NOW() WHERE id = $1;

-- name: UpdateLocationStatus :exec
UPDATE locations SET status = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateLocationWorkingHours :exec
UPDATE locations SET working_hours = $2, updated_at = NOW() WHERE id = $1
RETURNING *;

-- name: DeleteLocation :exec
DELETE FROM locations WHERE id = $1;

-- name: FindCategoryByCode :one
SELECT * FROM location_categories WHERE code = $1;
