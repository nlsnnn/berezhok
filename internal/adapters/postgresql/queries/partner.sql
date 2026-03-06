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
