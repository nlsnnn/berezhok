-- Реквизиты получателя для выплат (1:1 с партнёром).
-- Отдельная таблица, не в partner_legal_info, т.к. жизненный цикл
-- и верификация другие (платёжный реквизит, а не юридический документ).
CREATE TABLE partner_payout_destinations (
    partner_id      UUID PRIMARY KEY REFERENCES partners(id) ON DELETE CASCADE,
    type            VARCHAR(20) NOT NULL CHECK (type IN ('sbp')),
    sbp_phone       VARCHAR(20),
    sbp_bank_id     VARCHAR(20),
    recipient_name  VARCHAR(200),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_payout_destinations_updated_at BEFORE UPDATE
    ON partner_payout_destinations FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE partner_payout_destinations IS 'Реквизиты партнёра для приёма выплат (SBP).';

-- Выплата = финансовый документ за период.
CREATE TABLE partner_payouts (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id               UUID NOT NULL REFERENCES partners(id) ON DELETE RESTRICT,
    period_start             TIMESTAMPTZ NOT NULL,
    period_end               TIMESTAMPTZ NOT NULL,
    gross_amount             NUMERIC(12,2) NOT NULL,
    commission_amount        NUMERIC(12,2) NOT NULL,
    commission_rate_applied  NUMERIC(5,4)  NOT NULL,
    net_amount               NUMERIC(12,2) NOT NULL,
    status                   VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending','processing','completed','failed')),
    provider                 VARCHAR(20) NOT NULL DEFAULT 'yookassa',
    provider_payout_id       VARCHAR(100),
    idempotency_key          VARCHAR(100) UNIQUE,
    error_message            TEXT,
    processed_at             TIMESTAMPTZ,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_partner_payouts_partner_status ON partner_payouts(partner_id, status);
CREATE INDEX idx_partner_payouts_status         ON partner_payouts(status);
CREATE INDEX idx_partner_payouts_partner_created ON partner_payouts(partner_id, created_at DESC);

CREATE TRIGGER set_partner_payouts_updated_at BEFORE UPDATE
    ON partner_payouts FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE partner_payouts IS 'Выплаты партнёрам (финансовые документы).';

-- Связка заказ ↔ выплата (PK на order_id гарантирует, что один заказ
-- может попасть только в одну выплату — safe retry калькуляции).
CREATE TABLE partner_payout_orders (
    payout_id        UUID NOT NULL REFERENCES partner_payouts(id) ON DELETE CASCADE,
    order_id         UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    order_amount     NUMERIC(12,2) NOT NULL,
    commission_part  NUMERIC(12,2) NOT NULL,
    PRIMARY KEY (order_id)
);

CREATE INDEX idx_payout_orders_payout ON partner_payout_orders(payout_id);

COMMENT ON TABLE partner_payout_orders IS 'Связка заказов с выплатой (snapshot сумм на момент расчёта).';
