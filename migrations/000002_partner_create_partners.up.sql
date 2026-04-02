
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Партнёры (юридические лица)
CREATE TABLE partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    legal_name VARCHAR(200) NOT NULL,
    brand_name VARCHAR(200),
    logo_url TEXT,
    parent_partner_id UUID REFERENCES partners(id) ON DELETE SET NULL,
    account_type VARCHAR(20) DEFAULT 'independent' CHECK (account_type IN ('independent', 'network_head', 'franchise')),

    -- Комиссии
    commission_rate NUMERIC(5, 4) NOT NULL DEFAULT 0.20 CHECK (commission_rate >= 0 AND commission_rate <= 1),
    promo_commission_rate NUMERIC(5, 4) CHECK (promo_commission_rate >= 0 AND promo_commission_rate <= 1),
    promo_commission_until DATE,

    status VARCHAR(20) NOT NULL DEFAULT 'pending_documents' CHECK (status IN ('pending_documents', 'active', 'suspended', 'blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_partners_status ON partners(status);
CREATE INDEX idx_partners_parent ON partners(parent_partner_id);

CREATE TRIGGER set_partners_updated_at BEFORE UPDATE ON partners
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Юридическая информация партнёров
CREATE TABLE partner_legal_info (
    partner_id UUID PRIMARY KEY REFERENCES partners(id) ON DELETE CASCADE,
    inn VARCHAR(12) UNIQUE NOT NULL,
    ogrn VARCHAR(15),
    kpp VARCHAR(9),
    legal_address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_legal_info_inn ON partner_legal_info(inn);

CREATE TRIGGER set_legal_info_updated_at BEFORE UPDATE ON partner_legal_info
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE partners IS 'Партнёры (юридические лица)';
COMMENT ON TABLE partner_legal_info IS 'Юридическая информация партнёров';
