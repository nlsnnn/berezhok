-- Реквизиты получателя выплат

-- name: GetPayoutDestination :one
SELECT * FROM partner_payout_destinations WHERE partner_id = $1;

-- name: UpsertPayoutDestination :one
INSERT INTO partner_payout_destinations (
    partner_id, type, sbp_phone, sbp_bank_id, recipient_name
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (partner_id) DO UPDATE SET
    type            = EXCLUDED.type,
    sbp_phone       = EXCLUDED.sbp_phone,
    sbp_bank_id     = EXCLUDED.sbp_bank_id,
    recipient_name  = EXCLUDED.recipient_name,
    updated_at      = NOW()
RETURNING *;


-- Активные партнёры с тарифной информацией для воркера расчёта

-- name: ListActivePartnersForPayout :many
SELECT id, commission_rate, promo_commission_rate, promo_commission_until
FROM partners
WHERE status = 'active'
ORDER BY id;


-- Завершённые заказы партнёра за период, ещё не включённые ни в одну выплату

-- name: ListUnsettledCompletedOrders :many
SELECT o.id, o.amount, o.updated_at
FROM orders o
JOIN locations l ON l.id = o.location_id
LEFT JOIN partner_payout_orders po ON po.order_id = o.id
WHERE l.partner_id = $1
  AND o.status = 'completed'
  AND o.updated_at >= $2
  AND o.updated_at <  $3
  AND po.order_id IS NULL
ORDER BY o.updated_at ASC;


-- Создание выплаты и связки с заказами

-- name: CreatePayout :one
INSERT INTO partner_payouts (
    id, partner_id, period_start, period_end,
    gross_amount, commission_amount, commission_rate_applied, net_amount,
    status, provider, idempotency_key
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: AddPayoutOrder :exec
INSERT INTO partner_payout_orders (payout_id, order_id, order_amount, commission_part)
VALUES ($1, $2, $3, $4);


-- Диспатч: захват pending-выплат с блокировкой

-- name: LockPendingPayoutsForDispatch :many
SELECT id
FROM partner_payouts
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: GetPayoutByID :one
SELECT * FROM partner_payouts WHERE id = $1;

-- name: MarkPayoutProcessing :exec
UPDATE partner_payouts
SET status = 'processing'
WHERE id = $1 AND status = 'pending';

-- name: MarkPayoutCompleted :exec
UPDATE partner_payouts
SET status              = 'completed',
    provider_payout_id  = $2,
    processed_at        = NOW()
WHERE id = $1;

-- name: MarkPayoutFailed :exec
UPDATE partner_payouts
SET status         = 'failed',
    error_message  = $2,
    processed_at   = NOW()
WHERE id = $1;

-- name: ResetPayoutToPending :exec
UPDATE partner_payouts
SET status        = 'pending',
    error_message = NULL,
    processed_at  = NULL
WHERE id = $1 AND status = 'failed';

-- name: SetProviderPayoutID :exec
UPDATE partner_payouts
SET provider_payout_id = $2
WHERE id = $1;

-- name: ListProcessingPayouts :many
SELECT id, provider_payout_id
FROM partner_payouts
WHERE status = 'processing'
  AND provider_payout_id IS NOT NULL
ORDER BY created_at ASC
LIMIT $1;


-- История выплат партнёра

-- name: ListPayoutsByPartner :many
SELECT * FROM partner_payouts
WHERE partner_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPayoutsByPartner :one
SELECT COUNT(*)::bigint FROM partner_payouts WHERE partner_id = $1;

-- name: ListPayoutOrders :many
SELECT * FROM partner_payout_orders WHERE payout_id = $1;
