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

-- name: CheckEmailExists :one
SELECT EXISTS (
    SELECT 1 FROM partner_employees WHERE email = $1
    UNION
    SELECT 1 FROM partner_applications WHERE contact_email = $1 AND status = 'pending'
) AS email_exists;


-- name: GetPartnerProfile :one
SELECT 
    e.id as employee_id, e.name as employee_name, e.email, e.role, e.created_at as employee_created_at, e.must_change_password,
    p.id as partner_id, p.legal_name, p.brand_name, p.status as partner_status,
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

-- name: CreatePartner :one
INSERT INTO partners (
    legal_name, brand_name, logo_url, parent_partner_id, account_type, 
    commission_rate, promo_commission_rate, promo_commission_until, status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdatePartner :exec
UPDATE partners
SET legal_name = $1, brand_name = $2, logo_url = $3, parent_partner_id = $4, 
    account_type = $5, commission_rate = $6, promo_commission_rate = $7, 
    promo_commission_until = $8, status = $9
WHERE id = $10;


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