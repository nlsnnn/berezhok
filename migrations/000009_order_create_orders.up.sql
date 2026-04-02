-- 24.03.2026 20:33:00

-- Enum для статуса заказа
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'confirmed', 'completed', 'picked_up', 'cancelled', 'refunded', 'disputed');

-- Таблица заказов
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    box_id UUID NOT NULL REFERENCES surprise_boxes(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,

    -- Код для получения
    pickup_code VARCHAR(8) UNIQUE NOT NULL,
    qr_code_url TEXT,

    -- Финансы
    amount NUMERIC(10, 2) NOT NULL,

    -- Время получения
    pickup_time_start TIMESTAMPTZ NOT NULL,
    pickup_time_end TIMESTAMPTZ NOT NULL,

    status order_status NOT NULL DEFAULT 'pending',

    -- Подтверждение партнёром
    partner_confirmation_deadline TIMESTAMPTZ NOT NULL,
    partner_confirmed_at TIMESTAMPTZ,
    partner_confirmed_by UUID REFERENCES partner_employees(id) ON DELETE SET NULL,

    -- Отмена
    cancellation_reason TEXT,
    cancelled_at TIMESTAMPTZ,

    -- Выдача
    picked_up_at TIMESTAMPTZ,
    picked_up_confirmed_by UUID REFERENCES partner_employees(id) ON DELETE SET NULL,

    -- Получение клиентом
    user_confirmed_at TIMESTAMPTZ,
    auto_completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Индексы для оптимизации запросов
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_box_id ON orders(box_id);
CREATE INDEX idx_orders_location_id ON orders(location_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_pickup_time ON orders(pickup_time_start, pickup_time_end);
CREATE INDEX idx_orders_confirmation_deadline ON orders(partner_confirmation_deadline) WHERE status = 'paid';
CREATE INDEX idx_orders_customer_status ON orders(user_id, status);

-- Триггер для обновления updated_at
CREATE TRIGGER set_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE orders IS 'Заказы';
