-- name: CreatePayment :one
INSERT INTO payments (
    order_id, provider_payment_id, payment_url, method, provider, amount, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments WHERE order_id = $1;

-- name: UpdatePaymentStatus :one
UPDATE payments SET
    status = $2,
    paid_at = COALESCE($3, paid_at)
WHERE id = $1 RETURNING *;
