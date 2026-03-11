-- Заявки на партнёрство
-- name: ListApplications :many
SELECT * FROM partner_applications;

-- name: FindApplicationByID :one
SELECT * FROM partner_applications WHERE id = $1;

-- name: CreateApplication :one
INSERT INTO partner_applications (
    contact_name, contact_email, contact_phone, business_name, category_code, 
    address, description, status
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

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

-- name: GetPartnerProfile :one
SELECT 
    e.id as employee_id, e.name as employee_name, e.email, e.role,
    p.id as partner_id, p.legal_name, p.brand_name, p.status as partner_status,
    p.commission_rate, p.promo_commission_until,
    l.id as location_id, l.name as location_name, l.address as location_address
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