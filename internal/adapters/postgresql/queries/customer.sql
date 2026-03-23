-- Create new customer
-- name: CreateCustomer :one
INSERT INTO users (phone) 
VALUES ($1)
RETURNING id;

-- Get customer by ID
-- name: FindCustomerByID :one
SELECT * FROM users WHERE id = $1;

-- Get customer by phone
-- name: FindCustomerByPhone :one
SELECT * FROM users WHERE phone = $1;
