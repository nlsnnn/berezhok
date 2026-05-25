CREATE TABLE location_pins (
    code       VARCHAR(50) PRIMARY KEY,
    name_ru    VARCHAR(100) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE location_selected_pins (
    location_id UUID        NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    pin_code    VARCHAR(50) NOT NULL REFERENCES location_pins(code) ON DELETE CASCADE,
    PRIMARY KEY (location_id, pin_code)
);

CREATE INDEX idx_location_selected_pins_location ON location_selected_pins(location_id);

INSERT INTO location_pins (code, name_ru, sort_order) VALUES
    ('baked_goods', 'Выпечка',     1),
    ('desserts',    'Десерты',     2),
    ('sweets',      'Сладости',    3),
    ('pastries',    'Пирожные',    4),
    ('cakes',       'Торты',       5),
    ('sandwiches',  'Сэндвичи',    6),
    ('ready_food',  'Готовая еда', 7),
    ('beverages',   'Напитки',     8),
    ('salads',      'Салаты',      9),
    ('sushi',       'Суши',       10);
