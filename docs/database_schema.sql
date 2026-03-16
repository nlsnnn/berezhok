-- ============================================
-- Бережок — Полная схема базы данных
-- PostgreSQL 15+ с PostGIS
-- ============================================

-- Включаем расширения
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- ============================================
-- Функция для автообновления updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ============================================
-- МОДУЛЬ: CUSTOMERS (Клиенты)
-- ============================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(15) UNIQUE NOT NULL,
    name VARCHAR(100) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_phone ON users(phone);

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Push токены для уведомлений
CREATE TABLE user_push_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    platform VARCHAR(10) NOT NULL CHECK (platform IN ('ios', 'android')),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_push_tokens_user_id ON user_push_tokens(user_id);
CREATE INDEX idx_push_tokens_active ON user_push_tokens(user_id, active) WHERE active = TRUE;

-- ============================================
-- МОДУЛЬ: ADMINS (Администраторы)
-- ============================================

CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'admin' CHECK (role IN ('super_admin', 'admin', 'support')),
    permissions TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_users_email ON admin_users(email);
CREATE INDEX idx_admin_users_active ON admin_users(is_active);

-- Аудит действий администраторов
CREATE TABLE admin_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_user_id UUID NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    details JSONB,
    ip_address INET,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_admin ON admin_audit_log(admin_user_id);
CREATE INDEX idx_audit_log_created ON admin_audit_log(created_at DESC);
CREATE INDEX idx_audit_log_entity ON admin_audit_log(entity_type, entity_id);

-- ============================================
-- МОДУЛЬ: PARTNERS (Партнёры)
-- ============================================

-- Заявки на партнёрство
CREATE TABLE partner_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contact_name VARCHAR(100) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(15) NOT NULL,
    business_name VARCHAR(200) NOT NULL,
    category_code VARCHAR(50),
    address TEXT,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by UUID REFERENCES admin_users(id),
    reviewed_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_applications_status ON partner_applications(status);
CREATE INDEX idx_applications_created ON partner_applications(created_at DESC);

-- Партнёры (юридические лица)
CREATE TABLE partners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_partners_status ON partners(status);
CREATE INDEX idx_partners_parent ON partners(parent_partner_id);

CREATE TRIGGER update_partners_updated_at BEFORE UPDATE ON partners
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

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

-- Счета для выплат (может быть несколько у партнёра)
CREATE TABLE partner_payout_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    account_name VARCHAR(100),
    bank_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    bik VARCHAR(9) NOT NULL,
    correspondent_account VARCHAR(20),
    is_primary BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(partner_id, account_number)
);

CREATE INDEX idx_payout_accounts_partner ON partner_payout_accounts(partner_id);
CREATE INDEX idx_payout_accounts_primary ON partner_payout_accounts(partner_id, is_primary) WHERE is_primary = TRUE;

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
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

CREATE TRIGGER update_locations_updated_at BEFORE UPDATE ON locations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Сотрудники партнёров
CREATE TABLE partner_employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'employee' CHECK (role IN ('owner', 'manager', 'employee')),
    name VARCHAR(100),
    
    must_change_password BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_employees_partner ON partner_employees(partner_id);
CREATE INDEX idx_employees_location ON partner_employees(location_id);
CREATE INDEX idx_employees_email ON partner_employees(email);

-- ============================================
-- МОДУЛЬ: CATALOG (Каталог сюрприз-боксов)
-- ============================================

CREATE TABLE surprise_boxes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    original_price NUMERIC(10, 2),
    discount_price NUMERIC(10, 2) NOT NULL,
    quantity_available INT NOT NULL DEFAULT 0,
    
    pickup_time_start TIME NOT NULL,
    pickup_time_end TIME NOT NULL,
    
    image_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive' CHECK (status IN ('active', 'inactive', 'sold_out')),
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT positive_quantity CHECK (quantity_available >= 0),
    CONSTRAINT valid_price CHECK (discount_price > 0 AND (original_price IS NULL OR discount_price < original_price))
);

CREATE INDEX idx_boxes_location_id ON surprise_boxes(location_id);
CREATE INDEX idx_boxes_status ON surprise_boxes(status) WHERE status = 'active';

CREATE TRIGGER update_boxes_updated_at BEFORE UPDATE ON surprise_boxes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- МОДУЛЬ: ORDERS (Заказы)
-- ============================================

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    box_id UUID NOT NULL REFERENCES surprise_boxes(id) ON DELETE RESTRICT,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    
    -- Код для получения
    pickup_code VARCHAR(8) UNIQUE NOT NULL,
    qr_code_url TEXT,
    
    -- Финансы
    amount NUMERIC(10, 2) NOT NULL,
    
    -- Время получения
    pickup_time_start TIMESTAMP NOT NULL,
    pickup_time_end TIMESTAMP NOT NULL,
    
    -- Статусы
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'confirmed', 'picked_up', 'completed', 'cancelled', 'refunded', 'disputed')),
    
    -- Подтверждение партнёром
    partner_confirmation_deadline TIMESTAMP,
    partner_confirmed_at TIMESTAMP,
    partner_confirmed_by UUID REFERENCES partner_employees(id),
    
    -- Отмена
    cancellation_reason TEXT,
    cancelled_at TIMESTAMP,
    
    -- Выдача
    employee_confirmed_at TIMESTAMP,
    employee_confirmed_by UUID REFERENCES partner_employees(id),
    
    -- Получение клиентом
    user_confirmed_at TIMESTAMP,
    auto_completed_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_box_id ON orders(box_id);
CREATE INDEX idx_orders_location_id ON orders(location_id);
CREATE INDEX idx_orders_pickup_code ON orders(pickup_code);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_pickup_time ON orders(pickup_time_start, pickup_time_end);
CREATE INDEX idx_orders_confirmation_deadline ON orders(partner_confirmation_deadline) WHERE status = 'paid';

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- История статусов заказа (для аудита)
CREATE TABLE order_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_by_user_id UUID REFERENCES users(id),
    changed_by_employee_id UUID REFERENCES partner_employees(id),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_order_status_history_created ON order_status_history(created_at DESC);

-- ============================================
-- МОДУЛЬ: PAYMENTS (Платежи)
-- ============================================

-- Платежи (интеграция с ЮKassa)
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    
    -- Интеграция с ЮKassa
    external_payment_id VARCHAR(255) UNIQUE,
    payment_url TEXT,
    
    amount NUMERIC(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'succeeded', 'cancelled', 'failed')),
    
    payment_method JSONB,
    
    paid_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_external_id ON payments(external_payment_id);
CREATE INDEX idx_payments_status ON payments(status);

-- Выплаты партнёрам
CREATE TABLE payouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE RESTRICT,
    payout_account_id UUID NOT NULL REFERENCES partner_payout_accounts(id) ON DELETE RESTRICT,
    
    -- Финансы
    gross_amount NUMERIC(10, 2) NOT NULL,
    commission_rate NUMERIC(5, 4) NOT NULL,
    commission_amount NUMERIC(10, 2) NOT NULL,
    net_amount NUMERIC(10, 2) NOT NULL,
    
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    orders_count INT NOT NULL,
    
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    
    payout_details JSONB,
    
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payouts_partner_id ON payouts(partner_id);
CREATE INDEX idx_payouts_period ON payouts(period_start, period_end);
CREATE INDEX idx_payouts_status ON payouts(status);

-- Связь выплаты с заказами
CREATE TABLE payout_orders (
    payout_id UUID NOT NULL REFERENCES payouts(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    amount NUMERIC(10, 2) NOT NULL,
    PRIMARY KEY (payout_id, order_id)
);

CREATE INDEX idx_payout_orders_order_id ON payout_orders(order_id);

-- ============================================
-- МОДУЛЬ: REVIEWS (Отзывы)
-- ============================================

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT one_review_per_order UNIQUE (order_id)
);

CREATE INDEX idx_reviews_location_id ON reviews(location_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_created ON reviews(created_at DESC);

-- ============================================
-- МОДУЛЬ: DISPUTES (Споры)
-- ============================================

CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    
    initiated_by VARCHAR(10) NOT NULL CHECK (initiated_by IN ('user', 'partner')),
    
    reason TEXT NOT NULL,
    evidence_urls TEXT[] DEFAULT '{}',
    
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'under_review', 'resolved_for_user', 'resolved_for_partner')),
    
    resolution TEXT,
    resolved_by UUID REFERENCES admin_users(id),
    resolved_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disputes_order_id ON disputes(order_id);
CREATE INDEX idx_disputes_status ON disputes(status);

-- Сообщения в споре (если нужен чат)
CREATE TABLE dispute_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dispute_id UUID NOT NULL REFERENCES disputes(id) ON DELETE CASCADE,
    sender_type VARCHAR(10) NOT NULL CHECK (sender_type IN ('user', 'admin', 'partner')),
    sender_id UUID NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dispute_messages_dispute_id ON dispute_messages(dispute_id);
CREATE INDEX idx_dispute_messages_created ON dispute_messages(created_at ASC);

-- ============================================
-- МОДУЛЬ: NOTIFICATIONS (Уведомления)
-- ============================================

-- Лог отправленных уведомлений
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    type VARCHAR(50) NOT NULL CHECK (type IN ('sms', 'push', 'email')),
    channel VARCHAR(50) NOT NULL,
    
    recipient VARCHAR(255) NOT NULL,
    content JSONB NOT NULL,
    
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed')),
    error_message TEXT,
    
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_type_channel ON notifications(type, channel);
CREATE INDEX idx_notifications_created ON notifications(created_at DESC);

-- ============================================
-- ПРЕДСТАВЛЕНИЯ (Views) для удобства
-- ============================================

-- Активные заказы на сегодня по партнёру
CREATE VIEW partner_today_orders AS
SELECT 
    o.id,
    o.pickup_code,
    o.status,
    o.amount,
    o.pickup_time_start,
    o.pickup_time_end,
    u.phone as customer_phone,
    u.name as customer_name,
    l.id as location_id,
    l.name as location_name,
    sb.name as box_name,
    p.id as partner_id
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN locations l ON o.location_id = l.id
JOIN partners p ON l.partner_id = p.id
JOIN surprise_boxes sb ON o.box_id = sb.id
WHERE DATE(o.pickup_time_start) = CURRENT_DATE
  AND o.status IN ('paid', 'confirmed', 'picked_up');

-- Статистика заведений
CREATE VIEW location_stats AS
SELECT 
    l.id as location_id,
    l.name,
    COUNT(DISTINCT o.id) as total_orders,
    COALESCE(AVG(r.rating), 0) as avg_rating,
    COUNT(DISTINCT r.id) as total_reviews,
    SUM(CASE WHEN o.status = 'completed' THEN o.amount ELSE 0 END) as total_revenue
FROM locations l
LEFT JOIN orders o ON l.id = o.location_id
LEFT JOIN reviews r ON l.id = r.location_id
GROUP BY l.id, l.name;

-- ============================================
-- ФУНКЦИИ для бизнес-логики
-- ============================================

-- Функция генерации уникального кода заказа
CREATE OR REPLACE FUNCTION generate_pickup_code()
RETURNS VARCHAR(8) AS $$
DECLARE
    chars TEXT := 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';
    result VARCHAR(8) := '';
    i INT;
BEGIN
    FOR i IN 1..8 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Функция автоотмены неподтверждённых заказов
CREATE OR REPLACE FUNCTION cancel_expired_confirmations()
RETURNS TABLE(cancelled_order_id UUID) AS $$
BEGIN
    RETURN QUERY
    UPDATE orders
    SET status = 'cancelled',
        cancellation_reason = 'partner_confirmation_timeout',
        cancelled_at = NOW()
    WHERE status = 'paid'
      AND partner_confirmation_deadline < NOW()
    RETURNING id;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- КОММЕНТАРИИ к таблицам
-- ============================================

COMMENT ON TABLE users IS 'Клиенты мобильного приложения';
COMMENT ON TABLE admin_users IS 'Администраторы платформы';
COMMENT ON TABLE partners IS 'Партнёры (юридические лица)';
COMMENT ON TABLE locations IS 'Заведения партнёров (физические точки)';
COMMENT ON TABLE surprise_boxes IS 'Сюрприз-боксы для продажи';
COMMENT ON TABLE orders IS 'Заказы клиентов';
COMMENT ON TABLE payments IS 'Платежи через ЮKassa';
COMMENT ON TABLE payouts IS 'Выплаты партнёрам';
COMMENT ON TABLE reviews IS 'Отзывы клиентов о заведениях';
COMMENT ON TABLE disputes IS 'Споры между клиентами и партнёрами';

-- ============================================
-- Конец схемы
-- ============================================