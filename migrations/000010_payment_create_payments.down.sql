DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS payments;

DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS payment_provider;

DROP TRIGGER IF EXISTS set_payments_updated_at ON payments;
DROP FUNCTION IF EXISTS set_payments_updated_at();

DROP INDEX IF EXISTS idx_payments_order_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_provider_payment_id;
DROP INDEX IF EXISTS idx_payments_paid_at;
DROP INDEX IF EXISTS idx_payment_events_payment_id;