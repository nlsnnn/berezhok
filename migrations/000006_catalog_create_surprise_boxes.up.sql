-- 17.03.2026 10:20:00

CREATE TABLE IF NOT EXISTS surprise_boxes (
    id SERIAL PRIMARY KEY,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    description TEXT,
    original_price NUMERIC(10, 2),
    discount_price NUMERIC(10, 2) NOT NULL,
    quantity_available INT NOT NULL DEFAULT 0,

    pickup_time_start TIME NOT NULL,
    pickup_time_end TIME NOT NULL,

    image_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive' CHECK (status IN ('draft', 'active', 'inactive', 'sold_out')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    CONSTRAINT positive_quantity CHECK (quantity_available >= 0),
    CONSTRAINT valid_price CHECK (discount_price > 0 AND (original_price IS NULL OR discount_price < original_price))
);

CREATE INDEX idx_boxes_location_id ON surprise_boxes(location_id);
CREATE INDEX idx_boxes_status ON surprise_boxes(status) WHERE status = 'active';

CREATE TRIGGER set_surprise_boxes_updated_at BEFORE UPDATE ON surprise_boxes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();


COMMENT ON TABLE surprise_boxes IS 'Сюрприз-боксы для продажи';
