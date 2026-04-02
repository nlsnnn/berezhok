CREATE EXTENSION IF NOT EXISTS "postgis";

-- Категории заведений (справочник)
CREATE TABLE location_categories (
    code VARCHAR(50) PRIMARY KEY,
    name_ru VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    icon_url TEXT,
    color VARCHAR(7),
    sort_order INT DEFAULT 0
);

-- Предзаполнение категорий
INSERT INTO location_categories (code, name_ru, color, sort_order) VALUES
    ('bakery', 'Пекарня', '#FF6B6B', 1),
    ('cafe', 'Кафе', '#4ECDC4', 2),
    ('restaurant', 'Ресторан', '#45B7D1', 3),
    ('grocery', 'Продуктовый магазин', '#FFA07A', 4),
    ('hotel', 'Отель', '#98D8C8', 5);

-- Заведения (точки партнёров)
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    category_code VARCHAR(50) NOT NULL REFERENCES location_categories(code),

    name VARCHAR(200) NOT NULL,
    address TEXT NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    phone VARCHAR(15),

    logo_url TEXT,
    cover_image_url TEXT,
    gallery_urls TEXT[] DEFAULT '{}',

    working_hours JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive' CHECK (status IN ('active', 'inactive', 'closed')),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_locations_partner ON locations(partner_id);
CREATE INDEX idx_locations_category ON locations(category_code);
CREATE INDEX idx_locations_geography ON locations USING GIST(location);
CREATE INDEX idx_locations_status ON locations(status) WHERE status = 'active';

CREATE TRIGGER set_locations_updated_at BEFORE UPDATE ON locations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE locations IS 'Заведения (точки партнёров)';

-- Сотрудники партнёров
CREATE TABLE partner_employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'employee' CHECK (role IN ('owner', 'manager', 'employee')),
    name VARCHAR(100),

    must_change_password BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_employees_partner ON partner_employees(partner_id);
CREATE INDEX idx_employees_location ON partner_employees(location_id);
CREATE INDEX idx_employees_email ON partner_employees(email);

CREATE TRIGGER set_employees_updated_at BEFORE UPDATE ON partner_employees
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE partner_employees IS 'Сотрудники партнёров';
