-- name: CreatePayment :one
WITH inserted AS (
INSERT INTO payments (
    order_id, provider_payment_id, payment_url, method, provider, amount, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (order_id) DO NOTHING
RETURNING *
)
SELECT *
FROM inserted
UNION ALL
SELECT *
FROM payments
WHERE order_id = $1
  AND NOT EXISTS (SELECT 1 FROM inserted)
LIMIT 1;

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: UpdatePaymentStatus :one
UPDATE payments SET
    status = $2,
    paid_at = COALESCE($3, paid_at)
WHERE id = $1 RETURNING *;

-- name: CreateEvent :one
INSERT INTO payment_events (
    payment_id, event_type, payload
) VALUES (
    $1, $2, $3
) RETURNING *;
