-- 25.03.2026 3:36:00

CREATE TYPE payment_status AS ENUM ('pending', 'succeeded', 'cancelled', 'failed', 'refunded');

CREATE TYPE payment_method AS ENUM ('bank_card', 'sbp', 'cash', 'wallet', 'other');
CREATE TYPE payment_provider AS ENUM ('yookassa', 'tinkoff', 'stripe');

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,

    -- Платёжная система может возвращать свой ID платежа для последующей сверки
    provider_payment_id VARCHAR(255) UNIQUE,
    -- URL для оплаты (redirect)
    payment_url TEXT, 
    -- Метод оплаты, выбранный пользователем (может быть NULL, если ещё не выбран)
    method payment_method,
    provider payment_provider DEFAULT 'yookassa',

    -- Сумма для оплаты (может отличаться от суммы заказа из-за скидок, налогов и т.д.)
    amount NUMERIC(10, 2) NOT NULL,

    -- Статус платежа
    status payment_status NOT NULL DEFAULT 'pending',

    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_paid_at ON payments(paid_at);

CREATE TRIGGER set_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE payments IS 'Платежи за заказы';

CREATE TABLE IF NOT EXISTS payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_events_payment_id ON payment_events(payment_id);

COMMENT ON TABLE payment_events IS 'События, связанные с платежами (для аудита и отладки)';