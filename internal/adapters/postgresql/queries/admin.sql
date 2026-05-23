-- name: FindAdminByEmail :one
SELECT *
FROM admin_users
WHERE email = $1;

-- name: FindAdminByID :one
SELECT *
FROM admin_users
WHERE id = $1;

-- name: CreateAdminUser :one
INSERT INTO admin_users (
    email,
    password_hash,
    name,
    role,
    is_active
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateAdminLastLogin :exec
UPDATE admin_users
SET last_login_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: ListAdminUsers :many
SELECT *
FROM admin_users
WHERE (
    sqlc.arg(search)::text = ''
    OR email ILIKE '%' || sqlc.arg(search)::text || '%'
    OR name ILIKE '%' || sqlc.arg(search)::text || '%'
)
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminUsers :one
SELECT COUNT(*)
FROM admin_users
WHERE (
    sqlc.arg(search)::text = ''
    OR email ILIKE '%' || sqlc.arg(search)::text || '%'
    OR name ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: UpdateAdminUser :one
UPDATE admin_users
SET name = COALESCE(sqlc.narg(name), name),
    role = COALESCE(sqlc.narg(role), role),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeactivateAdminUser :execrows
UPDATE admin_users
SET is_active = FALSE,
    updated_at = NOW()
WHERE id = $1
  AND is_active = TRUE;

-- name: CreateAdminAuditLog :one
INSERT INTO admin_audit_log (
    admin_user_id,
    action,
    entity_type,
    entity_id,
    details,
    ip_address
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListAdminAuditLog :many
SELECT al.*,
       au.email AS admin_email,
       au.name AS admin_name
FROM admin_audit_log al
JOIN admin_users au ON au.id = al.admin_user_id
WHERE (
    sqlc.arg(action_filter)::text = ''
    OR al.action = sqlc.arg(action_filter)::text
)
AND (
    sqlc.arg(entity_type_filter)::text = ''
    OR al.entity_type = sqlc.arg(entity_type_filter)::text
)
ORDER BY al.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminAuditLog :one
SELECT COUNT(*)
FROM admin_audit_log al
WHERE (
    sqlc.arg(action_filter)::text = ''
    OR al.action = sqlc.arg(action_filter)::text
)
AND (
    sqlc.arg(entity_type_filter)::text = ''
    OR al.entity_type = sqlc.arg(entity_type_filter)::text
);

-- name: ListAdminApplications :many
SELECT *
FROM partner_applications
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR status = sqlc.arg(status_filter)::text
)
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminApplications :one
SELECT COUNT(*)
FROM partner_applications
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR status = sqlc.arg(status_filter)::text
);

-- name: MarkApplicationReviewed :execrows
UPDATE partner_applications
SET reviewed_by = $2,
    reviewed_at = NOW()
WHERE id = $1;

-- name: ListAdminPartners :many
SELECT p.id,
       p.brand_name,
       COALESCE(pli.legal_name, '') AS legal_name,
       p.logo_url,
       p.account_type,
       p.commission_rate,
       p.promo_commission_rate,
       p.promo_commission_until,
       p.status,
       p.created_at,
       p.updated_at,
       COUNT(DISTINCT l.id)::int AS locations_count,
       COUNT(DISTINCT o.id)::int AS total_orders,
       COALESCE(SUM(o.amount) FILTER (WHERE o.status IN ('completed', 'picked_up')), 0)::numeric AS total_revenue
FROM partners p
LEFT JOIN partner_legal_info pli ON pli.partner_id = p.id
LEFT JOIN locations l ON l.partner_id = p.id
LEFT JOIN orders o ON o.location_id = l.id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR p.status = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(search)::text = ''
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR pli.legal_name ILIKE '%' || sqlc.arg(search)::text || '%'
)
GROUP BY p.id, pli.legal_name
ORDER BY p.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminPartners :one
SELECT COUNT(DISTINCT p.id)
FROM partners p
LEFT JOIN partner_legal_info pli ON pli.partner_id = p.id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR p.status = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(search)::text = ''
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR pli.legal_name ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: GetAdminPartnerByID :one
SELECT p.id,
       p.brand_name,
       COALESCE(pli.legal_name, '') AS legal_name,
       p.logo_url,
       p.parent_partner_id,
       p.account_type,
       p.commission_rate,
       p.promo_commission_rate,
       p.promo_commission_until,
       p.status,
       p.created_at,
       p.updated_at,
       COALESCE(pli.inn, '') AS inn,
       COALESCE(pli.ogrn, '') AS ogrn,
       COALESCE(pli.kpp, '') AS kpp,
       COALESCE(pli.legal_address, '') AS legal_address,
       COUNT(DISTINCT l.id)::int AS locations_count,
       COUNT(DISTINCT o.id)::int AS total_orders,
       COALESCE(SUM(o.amount) FILTER (WHERE o.status IN ('completed', 'picked_up')), 0)::numeric AS total_revenue
FROM partners p
LEFT JOIN partner_legal_info pli ON pli.partner_id = p.id
LEFT JOIN locations l ON l.partner_id = p.id
LEFT JOIN orders o ON o.location_id = l.id
WHERE p.id = $1
GROUP BY p.id, pli.partner_id;

-- name: UpdateAdminPartner :one
UPDATE partners
SET status = COALESCE(sqlc.narg(status), status),
    commission_rate = COALESCE(sqlc.narg(commission_rate), commission_rate),
    promo_commission_rate = sqlc.narg(promo_commission_rate),
    promo_commission_until = sqlc.narg(promo_commission_until),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListAdminLocations :many
SELECT l.id,
       l.partner_id,
       p.brand_name AS partner_name,
       l.category_code,
       l.name,
       l.address,
       COALESCE(l.phone, '') AS phone,
       COALESCE(l.logo_url, '') AS logo_url,
       COALESCE(l.cover_image_url, '') AS cover_image_url,
       l.status,
       l.created_at,
       l.updated_at,
       COUNT(DISTINCT sb.id)::int AS boxes_count,
       COUNT(DISTINCT o.id)::int AS total_orders
FROM locations l
JOIN partners p ON p.id = l.partner_id
LEFT JOIN surprise_boxes sb ON sb.location_id = l.id
LEFT JOIN orders o ON o.location_id = l.id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR l.status = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(partner_id_filter)::text = ''
    OR l.partner_id = sqlc.arg(partner_id_filter)::uuid
)
AND (
    sqlc.arg(search)::text = ''
    OR l.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR l.address ILIKE '%' || sqlc.arg(search)::text || '%'
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
)
GROUP BY l.id, p.brand_name
ORDER BY l.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminLocations :one
SELECT COUNT(*)
FROM locations l
JOIN partners p ON p.id = l.partner_id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR l.status = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(partner_id_filter)::text = ''
    OR l.partner_id = sqlc.arg(partner_id_filter)::uuid
)
AND (
    sqlc.arg(search)::text = ''
    OR l.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR l.address ILIKE '%' || sqlc.arg(search)::text || '%'
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: GetAdminLocationByID :one
SELECT l.id,
       l.partner_id,
       p.brand_name AS partner_name,
       l.category_code,
       l.name,
       l.address,
       COALESCE(l.phone, '') AS phone,
       COALESCE(l.logo_url, '') AS logo_url,
       COALESCE(l.cover_image_url, '') AS cover_image_url,
       l.gallery_urls,
       l.working_hours,
       l.status,
       l.created_at,
       l.updated_at,
       COUNT(DISTINCT sb.id)::int AS boxes_count,
       COUNT(DISTINCT o.id)::int AS total_orders
FROM locations l
JOIN partners p ON p.id = l.partner_id
LEFT JOIN surprise_boxes sb ON sb.location_id = l.id
LEFT JOIN orders o ON o.location_id = l.id
WHERE l.id = $1
GROUP BY l.id, p.brand_name;

-- name: UpdateAdminLocationStatus :one
UPDATE locations
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListAdminBoxes :many
SELECT sb.id,
       sb.location_id,
       l.name AS location_name,
       l.partner_id,
       p.brand_name AS partner_name,
       sb.name,
       COALESCE(sb.description, '') AS description,
       sb.original_price,
       sb.discount_price,
       sb.quantity_available,
       sb.pickup_time_start,
       sb.pickup_time_end,
       COALESCE(sb.image_url, '') AS image_url,
       sb.status,
       sb.created_at,
       sb.updated_at,
       COUNT(o.id)::int AS total_orders
FROM surprise_boxes sb
JOIN locations l ON l.id = sb.location_id
JOIN partners p ON p.id = l.partner_id
LEFT JOIN orders o ON o.box_id = sb.id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR sb.status = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(location_id_filter)::text = ''
    OR sb.location_id = sqlc.arg(location_id_filter)::uuid
)
AND (
    sqlc.arg(search)::text = ''
    OR sb.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR l.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
)
GROUP BY sb.id, l.id, p.id
ORDER BY sb.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminBoxes :one
SELECT COUNT(*)
FROM surprise_boxes sb
JOIN locations l ON l.id = sb.location_id
JOIN partners p ON p.id = l.partner_id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR sb.status = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(location_id_filter)::text = ''
    OR sb.location_id = sqlc.arg(location_id_filter)::uuid
)
AND (
    sqlc.arg(search)::text = ''
    OR sb.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR l.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: GetAdminBoxByID :one
SELECT sb.id,
       sb.location_id,
       l.name AS location_name,
       l.partner_id,
       p.brand_name AS partner_name,
       sb.name,
       COALESCE(sb.description, '') AS description,
       sb.original_price,
       sb.discount_price,
       sb.quantity_available,
       sb.pickup_time_start,
       sb.pickup_time_end,
       COALESCE(sb.image_url, '') AS image_url,
       sb.status,
       sb.created_at,
       sb.updated_at,
       COUNT(o.id)::int AS total_orders
FROM surprise_boxes sb
JOIN locations l ON l.id = sb.location_id
JOIN partners p ON p.id = l.partner_id
LEFT JOIN orders o ON o.box_id = sb.id
WHERE sb.id = $1
GROUP BY sb.id, l.id, p.id;

-- name: UpdateAdminBoxStatus :one
UPDATE surprise_boxes
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListAdminCustomers :many
SELECT u.id,
       u.phone,
       u.name,
       u.created_at,
       u.updated_at,
       COUNT(o.id)::int AS total_orders,
       COALESCE(SUM(o.amount) FILTER (WHERE o.status IN ('completed', 'picked_up')), 0)::numeric AS total_spent
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE (
    sqlc.arg(search)::text = ''
    OR u.phone ILIKE '%' || sqlc.arg(search)::text || '%'
    OR u.name ILIKE '%' || sqlc.arg(search)::text || '%'
)
GROUP BY u.id
ORDER BY u.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminCustomers :one
SELECT COUNT(*)
FROM users u
WHERE (
    sqlc.arg(search)::text = ''
    OR u.phone ILIKE '%' || sqlc.arg(search)::text || '%'
    OR u.name ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: GetAdminCustomerByID :one
SELECT u.id,
       u.phone,
       u.name,
       u.created_at,
       u.updated_at,
       COUNT(o.id)::int AS total_orders,
       COALESCE(SUM(o.amount) FILTER (WHERE o.status IN ('completed', 'picked_up')), 0)::numeric AS total_spent
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE u.id = $1
GROUP BY u.id;

-- name: ListAdminOrders :many
SELECT o.id,
       o.user_id,
       u.phone AS customer_phone,
       COALESCE(u.name, '') AS customer_name,
       o.box_id,
       sb.name AS box_name,
       o.location_id,
       l.name AS location_name,
       l.partner_id,
       p.brand_name AS partner_name,
       o.pickup_code,
       o.amount,
       o.pickup_time_start,
       o.pickup_time_end,
       o.status,
       o.created_at,
       o.updated_at
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN surprise_boxes sb ON sb.id = o.box_id
JOIN locations l ON l.id = o.location_id
JOIN partners p ON p.id = l.partner_id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR o.status::text = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(search)::text = ''
    OR o.pickup_code ILIKE '%' || sqlc.arg(search)::text || '%'
    OR u.phone ILIKE '%' || sqlc.arg(search)::text || '%'
    OR sb.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR l.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
)
ORDER BY o.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminOrders :one
SELECT COUNT(*)
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN surprise_boxes sb ON sb.id = o.box_id
JOIN locations l ON l.id = o.location_id
JOIN partners p ON p.id = l.partner_id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR o.status::text = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(search)::text = ''
    OR o.pickup_code ILIKE '%' || sqlc.arg(search)::text || '%'
    OR u.phone ILIKE '%' || sqlc.arg(search)::text || '%'
    OR sb.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR l.name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR p.brand_name ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: GetAdminOrderByID :one
SELECT o.*,
       u.phone AS customer_phone,
       COALESCE(u.name, '') AS customer_name,
       sb.name AS box_name,
       l.name AS location_name,
       l.address AS location_address,
       l.partner_id,
       p.brand_name AS partner_name
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN surprise_boxes sb ON sb.id = o.box_id
JOIN locations l ON l.id = o.location_id
JOIN partners p ON p.id = l.partner_id
WHERE o.id = $1;

-- name: ListAdminPayments :many
SELECT pay.id,
       pay.order_id,
       pay.provider_payment_id,
       pay.payment_url,
       pay.method,
       pay.provider,
       pay.amount,
       pay.status,
       pay.paid_at,
       pay.created_at,
       pay.updated_at,
       o.pickup_code,
       u.phone AS customer_phone,
       l.name AS location_name,
       p.brand_name AS partner_name
FROM payments pay
JOIN orders o ON o.id = pay.order_id
JOIN users u ON u.id = o.user_id
JOIN locations l ON l.id = o.location_id
JOIN partners p ON p.id = l.partner_id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR pay.status::text = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(search)::text = ''
    OR pay.provider_payment_id ILIKE '%' || sqlc.arg(search)::text || '%'
    OR o.pickup_code ILIKE '%' || sqlc.arg(search)::text || '%'
    OR u.phone ILIKE '%' || sqlc.arg(search)::text || '%'
)
ORDER BY pay.created_at DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAdminPayments :one
SELECT COUNT(*)
FROM payments pay
JOIN orders o ON o.id = pay.order_id
JOIN users u ON u.id = o.user_id
WHERE (
    sqlc.arg(status_filter)::text = ''
    OR pay.status::text = sqlc.arg(status_filter)::text
)
AND (
    sqlc.arg(search)::text = ''
    OR pay.provider_payment_id ILIKE '%' || sqlc.arg(search)::text || '%'
    OR o.pickup_code ILIKE '%' || sqlc.arg(search)::text || '%'
    OR u.phone ILIKE '%' || sqlc.arg(search)::text || '%'
);

-- name: GetAdminPaymentByID :one
SELECT pay.*,
       o.pickup_code,
       u.phone AS customer_phone,
       l.name AS location_name,
       p.brand_name AS partner_name
FROM payments pay
JOIN orders o ON o.id = pay.order_id
JOIN users u ON u.id = o.user_id
JOIN locations l ON l.id = o.location_id
JOIN partners p ON p.id = l.partner_id
WHERE pay.id = $1;

-- name: ListAdminPaymentEvents :many
SELECT *
FROM payment_events
WHERE payment_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountAdminPaymentEvents :one
SELECT COUNT(*)
FROM payment_events
WHERE payment_id = $1;

-- name: GetAdminStats :one
SELECT (SELECT COUNT(*) FROM users)::int AS customers_total,
       (SELECT COUNT(*) FROM partners)::int AS partners_total,
       (SELECT COUNT(*) FROM partners WHERE status = 'active')::int AS partners_active,
       (SELECT COUNT(*) FROM locations)::int AS locations_total,
       (SELECT COUNT(*) FROM surprise_boxes)::int AS boxes_total,
       (SELECT COUNT(*) FROM orders)::int AS orders_total,
       (SELECT COUNT(*) FROM orders WHERE status = 'completed')::int AS orders_completed,
       (SELECT COUNT(*) FROM orders WHERE status = 'cancelled')::int AS orders_cancelled,
       (SELECT COUNT(*) FROM orders WHERE status = 'disputed')::int AS orders_disputed,
       (SELECT COALESCE(SUM(amount), 0)::numeric FROM orders WHERE status IN ('completed', 'picked_up')) AS gross_revenue,
       (SELECT COUNT(*) FROM payments)::int AS payments_total,
       (SELECT COUNT(*) FROM payments WHERE status = 'succeeded')::int AS payments_succeeded;
