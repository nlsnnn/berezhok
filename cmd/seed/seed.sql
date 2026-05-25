TRUNCATE TABLE
    admin_audit_log,
    payment_events,
    payments,
    reviews,
    partner_payout_orders,
    partner_payouts,
    partner_payout_destinations,
    orders,
    media_files,
    partner_employees,
    surprise_boxes,
    locations,
    partner_legal_info,
    users,
    partners,
    partner_applications,
    admin_users
RESTART IDENTITY CASCADE;

-- ─── ADMIN USERS ──────────────────────────────────────────────────────────────

INSERT INTO admin_users (id, email, password_hash, name, role, is_active, created_at, updated_at) VALUES
    ('01000000-0000-0000-0000-000000000001', 'admin@berezhok.local',     '__ADMIN_PASSWORD_HASH__', 'Главный администратор', 'super_admin', TRUE, NOW() - INTERVAL '90 days', NOW() - INTERVAL '1 day'),
    ('01000000-0000-0000-0000-000000000002', 'moderator@berezhok.local', '__ADMIN_PASSWORD_HASH__', 'Модератор площадки',    'admin',       TRUE, NOW() - INTERVAL '60 days', NOW() - INTERVAL '3 days'),
    ('01000000-0000-0000-0000-000000000003', 'support@berezhok.local',   '__ADMIN_PASSWORD_HASH__', 'Поддержка клиентов',    'support',     TRUE, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day');

-- ─── PARTNER APPLICATIONS ─────────────────────────────────────────────────────

INSERT INTO partner_applications (
    id, contact_name, contact_email, contact_phone, business_name, category_code,
    address, description, status, latitude, longitude, reviewed_at, rejection_reason, created_at
) VALUES
    ('10000000-0000-0000-0000-000000000001', 'Ирина Соколова',  'pending.partner@berezhok.local',   '+79991110001', 'Булочка у дома',     'bakery',     'Москва, ул. Покровка, 19',       'Небольшая ремесленная пекарня рядом с метро.',    'pending',  55.7570, 37.6476, NULL,                          NULL,                                   NOW() - INTERVAL '2 days'),
    ('10000000-0000-0000-0000-000000000002', 'Олег Миронов',    'rejected.partner@berezhok.local',  '+79991110002', 'Сытный угол',        'restaurant', 'Москва, ул. Маросейка, 8',       'Семейный ресторан с авторской кухней.',           'rejected', 55.7576, 37.6358, NOW() - INTERVAL '5 days',     'Не приложены юридические документы.',          NOW() - INTERVAL '7 days'),
    ('10000000-0000-0000-0000-000000000003', 'Алина Тарасова',  'sushi.new@berezhok.local',         '+79991110003', 'Роллы Фреш',         'restaurant', 'Москва, ул. Новослободская, 3',  'Сеть суши-ресторанов с доставкой и самовывозом.', 'pending',  55.7760, 37.5998, NULL,                          NULL,                                   NOW() - INTERVAL '1 day'),
    ('10000000-0000-0000-0000-000000000004', 'Роман Ефимов',    'pizza.new@berezhok.local',         '+79991110004', 'Пицца Экспресс',     'restaurant', 'Москва, Дмитровское шоссе, 9',  'Пиццерия с открытой кухней.',                     'pending',  55.8150, 37.5700, NULL,                          NULL,                                   NOW() - INTERVAL '3 days'),
    ('10000000-0000-0000-0000-000000000005', 'Виктор Орлов',    'rejected2.partner@berezhok.local', '+79991110005', 'Фастфуд Стрит',      'restaurant', 'Москва, Варшавское шоссе, 100', 'Сеть фастфуда.',                                  'rejected', 55.6810, 37.6210, NOW() - INTERVAL '10 days',    'Деятельность не соответствует условиям платформы.', NOW() - INTERVAL '12 days'),
    ('10000000-0000-0000-0000-000000000006', 'Светлана Власова', 'flower.new@berezhok.local',       '+79991110006', 'Цветочный рынок',    'grocery',    'Москва, пр-т Мира, 55',         'Цветы и растения с коротким сроком реализации.',  'approved', 55.7930, 37.6360, NOW() - INTERVAL '15 days',    NULL,                                   NOW() - INTERVAL '20 days');

-- ─── PARTNERS ─────────────────────────────────────────────────────────────────

INSERT INTO partners (id, brand_name, logo_url, parent_partner_id, account_type, commission_rate, promo_commission_rate, promo_commission_until, status, created_at, updated_at) VALUES
    ('20000000-0000-0000-0000-000000000001', 'Хлеб и Кофе',      'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=800',  NULL,                                   'independent',  0.1800, 0.1200, CURRENT_DATE + INTERVAL '30 days', 'active',    NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 day'),
    ('20000000-0000-0000-0000-000000000002', 'Вечерняя Трапеза', 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=800',  NULL,                                   'independent',  0.2000, NULL,   NULL,                              'active',    NOW() - INTERVAL '45 days', NOW() - INTERVAL '2 days'),
    ('20000000-0000-0000-0000-000000000003', 'Городской Вкус',   'https://images.unsplash.com/photo-1520607162513-77705c0f0d4a?w=800',  NULL,                                   'network_head', 0.1600, 0.1000, CURRENT_DATE + INTERVAL '14 days', 'active',    NOW() - INTERVAL '90 days', NOW() - INTERVAL '3 days'),
    ('20000000-0000-0000-0000-000000000004', 'Роллы и Суши',     'https://images.unsplash.com/photo-1553621042-f6e147245754?w=800',    NULL,                                   'independent',  0.1800, NULL,   NULL,                              'active',    NOW() - INTERVAL '20 days', NOW() - INTERVAL '1 day'),
    ('20000000-0000-0000-0000-000000000005', 'Пицца Плюс',       'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=800',  NULL,                                   'independent',  0.2000, NULL,   NULL,                              'suspended', NOW() - INTERVAL '10 days', NOW() - INTERVAL '1 day');

-- ─── PARTNER LEGAL INFO ───────────────────────────────────────────────────────

INSERT INTO partner_legal_info (
    partner_id, inn, ogrn, kpp, legal_address, legal_name,
    verification_status, verification_comment, verified_by, verified_at, created_at, updated_at
) VALUES
    ('20000000-0000-0000-0000-000000000001', '7701000001', '1027701000001', '770101001', 'Москва, ул. Земляной Вал, 10',       'ООО Хлеб и Кофе',       'verified', NULL,                          '01000000-0000-0000-0000-000000000001', NOW() - INTERVAL '55 days', NOW() - INTERVAL '60 days', NOW() - INTERVAL '55 days'),
    ('20000000-0000-0000-0000-000000000002', '7701000002', '1027701000002', '770102001', 'Москва, ул. Тверская, 24',           'ООО Вечерняя Трапеза',  'verified', NULL,                          '01000000-0000-0000-0000-000000000001', NOW() - INTERVAL '40 days', NOW() - INTERVAL '45 days', NOW() - INTERVAL '40 days'),
    ('20000000-0000-0000-0000-000000000003', '7701000003', '1027701000003', '770103001', 'Москва, Ленинградский пр-т, 32',     'ООО Сеть Городской Вкус','pending',  NULL,                          NULL,                                  NULL,                       NOW() - INTERVAL '90 days', NOW() - INTERVAL '12 days'),
    ('20000000-0000-0000-0000-000000000004', '7701000004', '1027701000004', '770104001', 'Москва, ул. Новослободская, 5',      'ООО Роллы и Суши',      'pending',  NULL,                          NULL,                                  NULL,                       NOW() - INTERVAL '20 days', NOW() - INTERVAL '5 days'),
    ('20000000-0000-0000-0000-000000000005', '7701000005', '1027701000005', '770105001', 'Москва, Варшавское шоссе, 102',      'ООО Пицца Плюс',        'failed',   'ИНН не прошёл проверку ФНС.', '01000000-0000-0000-0000-000000000002', NOW() - INTERVAL '8 days',  NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days');

-- ─── LOCATIONS ────────────────────────────────────────────────────────────────

INSERT INTO locations (
    id, partner_id, category_code, name, address, location, phone,
    logo_url, cover_image_url, gallery_urls, working_hours, status, created_at, updated_at
) VALUES
    (
        '30000000-0000-0000-0000-000000000001',
        '20000000-0000-0000-0000-000000000001', 'bakery',
        'Хлеб и Кофе на Курской', 'Москва, ул. Земляной Вал, 21',
        ST_SetSRID(ST_MakePoint(37.6598, 55.7578), 4326), '+74950000001',
        'https://images.unsplash.com/photo-1483695028939-5bb13f8648b0?w=800',
        'https://images.unsplash.com/photo-1517686469429-8bdb88b9f907?w=1200',
        ARRAY['https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200','https://images.unsplash.com/photo-1517433670267-08bbd4be890f?w=1200'],
        '{"mon":{"open":"08:00","close":"21:00"},"tue":{"open":"08:00","close":"21:00"},"wed":{"open":"08:00","close":"21:00"},"thu":{"open":"08:00","close":"21:00"},"fri":{"open":"08:00","close":"22:00"},"sat":{"open":"09:00","close":"22:00"},"sun":{"open":"09:00","close":"20:00"}}'::jsonb,
        'active', NOW() - INTERVAL '58 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000002',
        '20000000-0000-0000-0000-000000000001', 'cafe',
        'Хлеб и Кофе на Бауманской', 'Москва, ул. Бауманская, 33/2',
        ST_SetSRID(ST_MakePoint(37.6790, 55.7732), 4326), '+74950000002',
        'https://images.unsplash.com/photo-1445116572660-236099ec97a0?w=800',
        'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=1200',
        ARRAY['https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=1200'],
        '{"mon":{"open":"08:30","close":"21:30"},"tue":{"open":"08:30","close":"21:30"},"wed":{"open":"08:30","close":"21:30"},"thu":{"open":"08:30","close":"21:30"},"fri":{"open":"08:30","close":"22:00"},"sat":{"open":"09:00","close":"22:00"},"sun":{"open":"09:00","close":"21:00"}}'::jsonb,
        'active', NOW() - INTERVAL '52 days', NOW() - INTERVAL '2 days'
    ),
    (
        '30000000-0000-0000-0000-000000000003',
        '20000000-0000-0000-0000-000000000002', 'restaurant',
        'Вечерняя Трапеза на Тверской', 'Москва, ул. Тверская, 18к1',
        ST_SetSRID(ST_MakePoint(37.6065, 55.7662), 4326), '+74950000003',
        'https://images.unsplash.com/photo-1552566626-52f8b828add9?w=800',
        'https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=1200',
        ARRAY['https://images.unsplash.com/photo-1424847651672-bf20a4b0982b?w=1200'],
        '{"mon":{"open":"11:00","close":"23:00"},"tue":{"open":"11:00","close":"23:00"},"wed":{"open":"11:00","close":"23:00"},"thu":{"open":"11:00","close":"23:00"},"fri":{"open":"11:00","close":"00:00"},"sat":{"open":"12:00","close":"00:00"},"sun":{"open":"12:00","close":"22:00"}}'::jsonb,
        'active', NOW() - INTERVAL '43 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000004',
        '20000000-0000-0000-0000-000000000003', 'grocery',
        'Городской Вкус Экспресс', 'Москва, ул. Сущёвский Вал, 5с1',
        ST_SetSRID(ST_MakePoint(37.6038, 55.7929), 4326), '+74950000004',
        'https://images.unsplash.com/photo-1542838132-92c53300491e?w=800',
        'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200',
        ARRAY['https://images.unsplash.com/photo-1516594798947-e65505dbb29d?w=1200'],
        '{"mon":{"open":"07:30","close":"23:00"},"tue":{"open":"07:30","close":"23:00"},"wed":{"open":"07:30","close":"23:00"},"thu":{"open":"07:30","close":"23:00"},"fri":{"open":"07:30","close":"23:30"},"sat":{"open":"08:00","close":"23:30"},"sun":{"open":"08:00","close":"23:00"}}'::jsonb,
        'active', NOW() - INTERVAL '86 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000005',
        '20000000-0000-0000-0000-000000000003', 'cafe',
        'Городской Вкус Кофейня', 'Москва, Ленинградский пр-т, 36',
        ST_SetSRID(ST_MakePoint(37.5458, 55.7903), 4326), '+74950000005',
        'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800',
        'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=1200',
        ARRAY['https://images.unsplash.com/photo-1445116572660-236099ec97a0?w=1200'],
        '{"mon":{"open":"08:00","close":"21:00"},"tue":{"open":"08:00","close":"21:00"},"wed":{"open":"08:00","close":"21:00"},"thu":{"open":"08:00","close":"21:00"},"fri":{"open":"08:00","close":"22:00"},"sat":{"open":"09:00","close":"22:00"},"sun":{"open":"09:00","close":"20:00"}}'::jsonb,
        'draft', NOW() - INTERVAL '80 days', NOW() - INTERVAL '5 days'
    ),
    (
        '30000000-0000-0000-0000-000000000006',
        '20000000-0000-0000-0000-000000000004', 'restaurant',
        'Роллы и Суши на Новослободской', 'Москва, ул. Новослободская, 11',
        ST_SetSRID(ST_MakePoint(37.5990, 55.7754), 4326), '+74950000006',
        'https://images.unsplash.com/photo-1553621042-f6e147245754?w=800',
        'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200',
        ARRAY['https://images.unsplash.com/photo-1559410545-0bdcd187e0a6?w=1200'],
        '{"mon":{"open":"11:00","close":"23:00"},"tue":{"open":"11:00","close":"23:00"},"wed":{"open":"11:00","close":"23:00"},"thu":{"open":"11:00","close":"23:00"},"fri":{"open":"11:00","close":"00:00"},"sat":{"open":"12:00","close":"00:00"},"sun":{"open":"12:00","close":"22:00"}}'::jsonb,
        'active', NOW() - INTERVAL '18 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000007',
        '20000000-0000-0000-0000-000000000005', 'restaurant',
        'Пицца Плюс на Варшавке', 'Москва, Варшавское шоссе, 98',
        ST_SetSRID(ST_MakePoint(37.6200, 55.6820), 4326), '+74950000007',
        'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=800',
        'https://images.unsplash.com/photo-1574071318508-1cdbab80d002?w=1200',
        ARRAY['https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=1200'],
        '{"mon":{"open":"10:00","close":"22:00"},"tue":{"open":"10:00","close":"22:00"},"wed":{"open":"10:00","close":"22:00"},"thu":{"open":"10:00","close":"22:00"},"fri":{"open":"10:00","close":"23:00"},"sat":{"open":"11:00","close":"23:00"},"sun":{"open":"11:00","close":"22:00"}}'::jsonb,
        'inactive', NOW() - INTERVAL '8 days', NOW() - INTERVAL '1 day'
    );

-- ─── PARTNER EMPLOYEES ────────────────────────────────────────────────────────

INSERT INTO partner_employees (
    id, partner_id, location_id, email, password_hash, role, name,
    must_change_password, last_login_at, created_at, updated_at
) VALUES
    ('40000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'owner@berezhok.local',         '__PARTNER_PASSWORD_HASH__', 'owner',    'Анна Орлова',      FALSE, NOW() - INTERVAL '2 hours',  NOW() - INTERVAL '58 days', NOW() - INTERVAL '2 hours'),
    ('40000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000002', 'employee@berezhok.local',       '__PARTNER_PASSWORD_HASH__', 'employee', 'Илья Кузнецов',    FALSE, NOW() - INTERVAL '6 hours',  NOW() - INTERVAL '50 days', NOW() - INTERVAL '6 hours'),
    ('40000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000003', 'dinner.owner@berezhok.local',   '__PARTNER_PASSWORD_HASH__', 'owner',    'Мария Белова',     FALSE, NOW() - INTERVAL '1 day',    NOW() - INTERVAL '42 days', NOW() - INTERVAL '1 day'),
    ('40000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000004', 'coffee.owner@berezhok.local',   '__PARTNER_PASSWORD_HASH__', 'owner',    'Денис Смирнов',    FALSE, NOW() - INTERVAL '3 hours',  NOW() - INTERVAL '85 days', NOW() - INTERVAL '3 hours'),
    ('40000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000005', 'coffee.manager@berezhok.local', '__PARTNER_PASSWORD_HASH__', 'manager',  'Екатерина Иванова',FALSE, NOW() - INTERVAL '9 hours',  NOW() - INTERVAL '70 days', NOW() - INTERVAL '9 hours'),
    ('40000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'bakery.manager@berezhok.local', '__PARTNER_PASSWORD_HASH__', 'manager',  'Сергей Попов',     TRUE,  NULL,                        NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
    ('40000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000003', 'dinner.staff@berezhok.local',   '__PARTNER_PASSWORD_HASH__', 'employee', 'Алексей Фёдоров',  TRUE,  NULL,                        NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'),
    ('40000000-0000-0000-0000-000000000008', '20000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000006', 'sushi.owner@berezhok.local',    '__PARTNER_PASSWORD_HASH__', 'owner',    'Алина Тарасова',   FALSE, NOW() - INTERVAL '5 hours',  NOW() - INTERVAL '18 days', NOW() - INTERVAL '5 hours'),
    ('40000000-0000-0000-0000-000000000009', '20000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000007', 'pizza.owner@berezhok.local',    '__PARTNER_PASSWORD_HASH__', 'owner',    'Роман Ефимов',     FALSE, NOW() - INTERVAL '2 days',  NOW() - INTERVAL '8 days',  NOW() - INTERVAL '2 days');

-- ─── MEDIA FILES ──────────────────────────────────────────────────────────────

INSERT INTO media_files (id, filename, original_filename, storage_key, url, content_type, size_bytes, uploaded_at) VALUES
    ('50000000-0000-0000-0000-000000000001', 'croissant-box.jpg',  'croissant-box.jpg',  'seed/croissant-box.jpg',  'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200', 'image/jpeg', 245120, NOW() - INTERVAL '20 days'),
    ('50000000-0000-0000-0000-000000000002', 'dinner-box.jpg',     'dinner-box.jpg',     'seed/dinner-box.jpg',     'https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=1200', 'image/jpeg', 318600, NOW() - INTERVAL '18 days'),
    ('50000000-0000-0000-0000-000000000003', 'sushi-box.jpg',      'sushi-box.jpg',      'seed/sushi-box.jpg',      'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200', 'image/jpeg', 276000, NOW() - INTERVAL '10 days'),
    ('50000000-0000-0000-0000-000000000004', 'grocery-box.jpg',    'grocery-box.jpg',    'seed/grocery-box.jpg',    'https://images.unsplash.com/photo-1516594798947-e65505dbb29d?w=1200', 'image/jpeg', 198400, NOW() - INTERVAL '15 days'),
    ('50000000-0000-0000-0000-000000000005', 'pizza-box.jpg',      'pizza-box.jpg',      'seed/pizza-box.jpg',      'https://images.unsplash.com/photo-1574071318508-1cdbab80d002?w=1200', 'image/jpeg', 302200, NOW() - INTERVAL '5 days');

-- ─── SURPRISE BOXES ───────────────────────────────────────────────────────────

INSERT INTO surprise_boxes (
    id, location_id, name, description, original_price, discount_price,
    quantity_available, pickup_time_start, pickup_time_end, image_url, status, created_at, updated_at
) VALUES
    ('60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'Утренний хлебный бокс',      'Свежая выпечка, хлеб дня и небольшой сладкий десерт.',          890.00,  390.00, 8, '20:00', '21:30', 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200', 'active',   NOW() - INTERVAL '25 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 'Сладкий вечер',              'Набор десертов и слойка из витрины конца дня.',                  760.00,  320.00, 4, '19:30', '21:00', 'https://images.unsplash.com/photo-1517433670267-08bbd4be890f?w=1200', 'active',   NOW() - INTERVAL '16 days', NOW() - INTERVAL '5 hours'),
    ('60000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', 'Кофе и выпечка',             'Сэндвич, булочка и напиток на выбор.',                           950.00,  430.00, 6, '19:00', '21:30', 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=1200', 'active',   NOW() - INTERVAL '20 days', NOW() - INTERVAL '3 hours'),
    ('60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'Ресторанный ужин',           'Горячее блюдо, гарнир и десерт-сюрприз.',                       1800.00, 790.00, 5, '21:00', '22:30', 'https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=1200', 'active',   NOW() - INTERVAL '22 days', NOW() - INTERVAL '2 hours'),
    ('60000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000003', 'Ланч-бокс шефа',             'Салат, суп и главное блюдо по выбору кухни.',                   1450.00, 690.00, 0, '15:00', '16:30', 'https://images.unsplash.com/photo-1544025162-d76694265947?w=1200', 'sold_out', NOW() - INTERVAL '14 days', NOW() - INTERVAL '1 hour'),
    ('60000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000004', 'Фреш-маркет бокс',           'Овощи, салаты, молочные продукты и снеки с коротким сроком.',   1200.00, 520.00, 7, '20:30', '22:00', 'https://images.unsplash.com/photo-1516594798947-e65505dbb29d?w=1200', 'active',   NOW() - INTERVAL '18 days', NOW() - INTERVAL '4 hours'),
    ('60000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000004', 'Семейный продуктовый набор', 'Собранный набор из кулинарии и свежих продуктов.',               2100.00, 990.00, 3, '21:00', '22:45', 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200', 'active',   NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 hours'),
    ('60000000-0000-0000-0000-000000000008', '30000000-0000-0000-0000-000000000005', 'Кофейня позднего вечера',    'Кофе дня и десерты. Локация выключена для теста неактивной.',    990.00,  390.00, 5, '19:00', '20:30', 'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=1200', 'inactive', NOW() - INTERVAL '9 days',  NOW() - INTERVAL '2 days'),
    ('60000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000002', 'Черновик сезонного бокса',   'Черновик для проверки отображения неактивных карточек.',        1100.00, 450.00, 2, '18:00', '19:30', 'https://images.unsplash.com/photo-1445116572660-236099ec97a0?w=1200', 'draft',    NOW() - INTERVAL '6 days',  NOW() - INTERVAL '6 days'),
    ('60000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000006', 'Суши-сет вечерний',          'Роллы и нигири из остатков сервисного ланча.',                  1600.00, 680.00, 6, '21:00', '22:30', 'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200', 'active',   NOW() - INTERVAL '15 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000011', '30000000-0000-0000-0000-000000000006', 'Ланч-бокс суши',             'Лёгкий набор на обед: роллы и мисо-суп.',                       980.00,  420.00, 4, '14:30', '16:00', 'https://images.unsplash.com/photo-1559410545-0bdcd187e0a6?w=1200', 'active',   NOW() - INTERVAL '12 days', NOW() - INTERVAL '3 hours'),
    ('60000000-0000-0000-0000-000000000012', '30000000-0000-0000-0000-000000000007', 'Пицца дня',                  'Остатки дня: целые пиццы и слайсы.',                            1300.00, 560.00, 0, '21:30', '22:45', 'https://images.unsplash.com/photo-1574071318508-1cdbab80d002?w=1200', 'inactive', NOW() - INTERVAL '6 days',  NOW() - INTERVAL '1 day');

-- ─── CUSTOMERS (USERS) ────────────────────────────────────────────────────────

INSERT INTO users (id, phone, name, created_at, updated_at) VALUES
    ('70000000-0000-0000-0000-000000000001', '+79990000001', 'Иван Петров',       NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day'),
    ('70000000-0000-0000-0000-000000000002', '+79990000002', 'Ольга Романова',    NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 days'),
    ('70000000-0000-0000-0000-000000000003', '+79990000003', 'Павел Лебедев',     NOW() - INTERVAL '21 days', NOW() - INTERVAL '5 days'),
    ('70000000-0000-0000-0000-000000000004', '+79990000004', 'Марина Крылова',    NOW() - INTERVAL '14 days', NOW() - INTERVAL '4 days'),
    ('70000000-0000-0000-0000-000000000005', '+79990000005', 'Тимур Сафин',       NOW() - INTERVAL '7 days',  NOW() - INTERVAL '1 day'),
    ('70000000-0000-0000-0000-000000000006', '+79990000006', 'Алексей Новиков',   NOW() - INTERVAL '55 days', NOW() - INTERVAL '1 day'),
    ('70000000-0000-0000-0000-000000000007', '+79990000007', 'Светлана Попова',   NOW() - INTERVAL '50 days', NOW() - INTERVAL '3 days'),
    ('70000000-0000-0000-0000-000000000008', '+79990000008', 'Дмитрий Козлов',    NOW() - INTERVAL '48 days', NOW() - INTERVAL '2 days'),
    ('70000000-0000-0000-0000-000000000009', '+79990000009', 'Наталья Морозова',  NOW() - INTERVAL '44 days', NOW() - INTERVAL '6 days'),
    ('70000000-0000-0000-0000-000000000010', '+79990000010', 'Артём Волков',      NOW() - INTERVAL '40 days', NOW() - INTERVAL '2 days'),
    ('70000000-0000-0000-0000-000000000011', '+79990000011', 'Юлия Соколова',     NOW() - INTERVAL '35 days', NOW() - INTERVAL '1 day'),
    ('70000000-0000-0000-0000-000000000012', '+79990000012', 'Максим Жуков',      NOW() - INTERVAL '28 days', NOW() - INTERVAL '5 days');

-- ─── ORDERS ───────────────────────────────────────────────────────────────────
-- Поля: id, user_id, box_id, location_id, pickup_code, qr_code_url, amount,
--       pickup_time_start, pickup_time_end, status,
--       partner_confirmation_deadline, partner_confirmed_at, partner_confirmed_by,
--       cancellation_reason, cancelled_at,
--       picked_up_at, picked_up_confirmed_by, user_confirmed_at, auto_completed_at,
--       created_at, updated_at

INSERT INTO orders (
    id, user_id, box_id, location_id, pickup_code, qr_code_url, amount,
    pickup_time_start, pickup_time_end, status,
    partner_confirmation_deadline, partner_confirmed_at, partner_confirmed_by,
    cancellation_reason, cancelled_at,
    picked_up_at, picked_up_confirmed_by, user_confirmed_at, auto_completed_at,
    created_at, updated_at
) VALUES
    -- Активные/текущие заказы
    ('80000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'PK100001', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100001',  390.00, NOW() + INTERVAL '2 hours', NOW() + INTERVAL '4 hours', 'paid',      NOW() + INTERVAL '90 minutes', NULL,                          NULL,                                   NULL,                        NULL,                       NULL,                       NULL,                                   NULL,                       NULL,                       NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '10 minutes'),
    ('80000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', 'PK100002', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100002',  430.00, NOW() + INTERVAL '3 hours', NOW() + INTERVAL '5 hours', 'confirmed', NOW() + INTERVAL '1 hour',     NOW() - INTERVAL '20 minutes', '40000000-0000-0000-0000-000000000002', NULL,                        NULL,                       NULL,                       NULL,                                   NULL,                       NULL,                       NOW() - INTERVAL '50 minutes', NOW() - INTERVAL '20 minutes'),
    ('80000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'PK100003', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100003',  790.00, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '1 hour',  'completed', NOW() - INTERVAL '6 hours',    NOW() - INTERVAL '5 hours',    '40000000-0000-0000-0000-000000000003', NULL,                        NULL,                       NOW() - INTERVAL '70 minutes',  '40000000-0000-0000-0000-000000000003', NOW() - INTERVAL '65 minutes',  NOW() - INTERVAL '60 minutes',  NOW() - INTERVAL '8 hours',    NOW() - INTERVAL '60 minutes'),
    ('80000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000004', 'PK100004', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100004',  520.00, NOW() - INTERVAL '1 day' + INTERVAL '2 hours', NOW() - INTERVAL '1 day' + INTERVAL '4 hours', 'completed', NOW() - INTERVAL '1 day' + INTERVAL '90 minutes', NOW() - INTERVAL '1 day' + INTERVAL '70 minutes', '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '20 hours', '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '19 hours', NOW() - INTERVAL '18 hours', NOW() - INTERVAL '26 hours', NOW() - INTERVAL '18 hours'),
    ('80000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000005', '60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'PK100005', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100005',  390.00, NOW() - INTERVAL '2 days' + INTERVAL '1 hour', NOW() - INTERVAL '2 days' + INTERVAL '3 hours', 'cancelled', NOW() - INTERVAL '2 days' + INTERVAL '30 minutes', NULL, NULL, 'Пользователь отменил заказ.', NOW() - INTERVAL '47 hours', NULL, NULL, NULL, NULL, NOW() - INTERVAL '49 hours', NOW() - INTERVAL '47 hours'),
    ('80000000-0000-0000-0000-000000000006', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000004', 'PK100006', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100006',  990.00, NOW() - INTERVAL '3 days' + INTERVAL '4 hours', NOW() - INTERVAL '3 days' + INTERVAL '6 hours', 'disputed',  NOW() - INTERVAL '3 days' + INTERVAL '2 hours', NOW() - INTERVAL '3 days' + INTERVAL '90 minutes', '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '66 hours', '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '65 hours', NULL, NOW() - INTERVAL '76 hours', NOW() - INTERVAL '65 hours'),
    ('80000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'PK100007', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100007',  790.00, NOW() - INTERVAL '5 days' + INTERVAL '2 hours', NOW() - INTERVAL '5 days' + INTERVAL '4 hours', 'completed', NOW() - INTERVAL '5 days' + INTERVAL '1 hour', NOW() - INTERVAL '5 days' + INTERVAL '50 minutes', '40000000-0000-0000-0000-000000000003', NULL, NULL, NOW() - INTERVAL '116 hours', '40000000-0000-0000-0000-000000000003', NOW() - INTERVAL '115 hours', NOW() - INTERVAL '114 hours', NOW() - INTERVAL '118 hours', NOW() - INTERVAL '114 hours'),
    ('80000000-0000-0000-0000-000000000008', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000004', 'PK100008', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100008',  520.00, NOW() - INTERVAL '6 days' + INTERVAL '2 hours', NOW() - INTERVAL '6 days' + INTERVAL '3 hours', 'completed', NOW() - INTERVAL '6 days' + INTERVAL '1 hour', NOW() - INTERVAL '6 days' + INTERVAL '55 minutes', '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '141 hours', '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '140 hours', NOW() - INTERVAL '139 hours', NOW() - INTERVAL '143 hours', NOW() - INTERVAL '139 hours'),
    ('80000000-0000-0000-0000-000000000009', '70000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', 'PK100009', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100009',  430.00, NOW() + INTERVAL '1 day' + INTERVAL '2 hours', NOW() + INTERVAL '1 day' + INTERVAL '4 hours', 'pending',   NOW() + INTERVAL '1 day' + INTERVAL '30 minutes', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),
    -- Исторические завершённые заказы (для расчёта выплат, 55–32 дня назад)
    ('80000000-0000-0000-0000-000000000010', '70000000-0000-0000-0000-000000000006', '60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'PK100010', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100010',  390.00, NOW() - INTERVAL '55 days' + INTERVAL '20 hours', NOW() - INTERVAL '55 days' + INTERVAL '22 hours', 'completed', NOW() - INTERVAL '55 days' + INTERVAL '19 hours', NOW() - INTERVAL '55 days' + INTERVAL '18 hours', '40000000-0000-0000-0000-000000000001', NULL, NULL, NOW() - INTERVAL '55 days' + INTERVAL '21 hours', '40000000-0000-0000-0000-000000000001', NOW() - INTERVAL '55 days' + INTERVAL '21 hours', NULL, NOW() - INTERVAL '55 days' + INTERVAL '10 hours', NOW() - INTERVAL '55 days' + INTERVAL '21 hours'),
    ('80000000-0000-0000-0000-000000000011', '70000000-0000-0000-0000-000000000007', '60000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 'PK100011', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100011',  320.00, NOW() - INTERVAL '50 days' + INTERVAL '20 hours', NOW() - INTERVAL '50 days' + INTERVAL '21 hours', 'completed', NOW() - INTERVAL '50 days' + INTERVAL '19 hours', NOW() - INTERVAL '50 days' + INTERVAL '18 hours', '40000000-0000-0000-0000-000000000001', NULL, NULL, NOW() - INTERVAL '50 days' + INTERVAL '20 hours', '40000000-0000-0000-0000-000000000001', NOW() - INTERVAL '50 days' + INTERVAL '20 hours', NULL, NOW() - INTERVAL '50 days' + INTERVAL '10 hours', NOW() - INTERVAL '50 days' + INTERVAL '20 hours'),
    ('80000000-0000-0000-0000-000000000012', '70000000-0000-0000-0000-000000000008', '60000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', 'PK100012', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100012',  430.00, NOW() - INTERVAL '47 days' + INTERVAL '20 hours', NOW() - INTERVAL '47 days' + INTERVAL '22 hours', 'completed', NOW() - INTERVAL '47 days' + INTERVAL '19 hours', NOW() - INTERVAL '47 days' + INTERVAL '18 hours', '40000000-0000-0000-0000-000000000002', NULL, NULL, NOW() - INTERVAL '47 days' + INTERVAL '21 hours', '40000000-0000-0000-0000-000000000002', NOW() - INTERVAL '47 days' + INTERVAL '21 hours', NULL, NOW() - INTERVAL '47 days' + INTERVAL '10 hours', NOW() - INTERVAL '47 days' + INTERVAL '21 hours'),
    ('80000000-0000-0000-0000-000000000013', '70000000-0000-0000-0000-000000000009', '60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'PK100013', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100013',  790.00, NOW() - INTERVAL '45 days' + INTERVAL '22 hours', NOW() - INTERVAL '45 days' + INTERVAL '23 hours', 'completed', NOW() - INTERVAL '45 days' + INTERVAL '21 hours', NOW() - INTERVAL '45 days' + INTERVAL '20 hours', '40000000-0000-0000-0000-000000000003', NULL, NULL, NOW() - INTERVAL '45 days' + INTERVAL '22 hours', '40000000-0000-0000-0000-000000000003', NOW() - INTERVAL '45 days' + INTERVAL '22 hours', NULL, NOW() - INTERVAL '45 days' + INTERVAL '11 hours', NOW() - INTERVAL '45 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000014', '70000000-0000-0000-0000-000000000010', '60000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000004', 'PK100014', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100014',  520.00, NOW() - INTERVAL '42 days' + INTERVAL '21 hours', NOW() - INTERVAL '42 days' + INTERVAL '23 hours', 'completed', NOW() - INTERVAL '42 days' + INTERVAL '20 hours', NOW() - INTERVAL '42 days' + INTERVAL '19 hours', '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '42 days' + INTERVAL '22 hours', '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '42 days' + INTERVAL '22 hours', NULL, NOW() - INTERVAL '42 days' + INTERVAL '10 hours', NOW() - INTERVAL '42 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000015', '70000000-0000-0000-0000-000000000011', '60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'PK100015', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100015',  790.00, NOW() - INTERVAL '40 days' + INTERVAL '22 hours', NOW() - INTERVAL '40 days' + INTERVAL '23 hours', 'completed', NOW() - INTERVAL '40 days' + INTERVAL '21 hours', NOW() - INTERVAL '40 days' + INTERVAL '20 hours', '40000000-0000-0000-0000-000000000003', NULL, NULL, NOW() - INTERVAL '40 days' + INTERVAL '22 hours', '40000000-0000-0000-0000-000000000003', NOW() - INTERVAL '40 days' + INTERVAL '22 hours', NULL, NOW() - INTERVAL '40 days' + INTERVAL '11 hours', NOW() - INTERVAL '40 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000016', '70000000-0000-0000-0000-000000000012', '60000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000004', 'PK100016', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100016',  990.00, NOW() - INTERVAL '38 days' + INTERVAL '22 hours', NOW() - INTERVAL '38 days' + INTERVAL '23 hours', 'completed', NOW() - INTERVAL '38 days' + INTERVAL '21 hours', NOW() - INTERVAL '38 days' + INTERVAL '20 hours', '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '38 days' + INTERVAL '22 hours', '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '38 days' + INTERVAL '22 hours', NULL, NOW() - INTERVAL '38 days' + INTERVAL '11 hours', NOW() - INTERVAL '38 days' + INTERVAL '22 hours'),
    -- Недавние заказы (последние 25 дней — для аналитики, pending payouts)
    ('80000000-0000-0000-0000-000000000017', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'PK100017', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100017',  390.00, NOW() - INTERVAL '20 days' + INTERVAL '20 hours', NOW() - INTERVAL '20 days' + INTERVAL '22 hours', 'completed', NOW() - INTERVAL '20 days' + INTERVAL '19 hours', NOW() - INTERVAL '20 days' + INTERVAL '18 hours', '40000000-0000-0000-0000-000000000001', NULL, NULL, NOW() - INTERVAL '20 days' + INTERVAL '21 hours', '40000000-0000-0000-0000-000000000001', NOW() - INTERVAL '20 days' + INTERVAL '21 hours', NULL, NOW() - INTERVAL '20 days' + INTERVAL '10 hours', NOW() - INTERVAL '20 days' + INTERVAL '21 hours'),
    ('80000000-0000-0000-0000-000000000018', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'PK100018', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100018',  790.00, NOW() - INTERVAL '15 days' + INTERVAL '22 hours', NOW() - INTERVAL '15 days' + INTERVAL '23 hours', 'completed', NOW() - INTERVAL '15 days' + INTERVAL '21 hours', NOW() - INTERVAL '15 days' + INTERVAL '20 hours', '40000000-0000-0000-0000-000000000003', NULL, NULL, NOW() - INTERVAL '15 days' + INTERVAL '22 hours', '40000000-0000-0000-0000-000000000003', NOW() - INTERVAL '15 days' + INTERVAL '22 hours', NULL, NOW() - INTERVAL '15 days' + INTERVAL '11 hours', NOW() - INTERVAL '15 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000019', '70000000-0000-0000-0000-000000000006', '60000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000004', 'PK100019', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100019',  520.00, NOW() - INTERVAL '12 days' + INTERVAL '21 hours', NOW() - INTERVAL '12 days' + INTERVAL '23 hours', 'completed', NOW() - INTERVAL '12 days' + INTERVAL '20 hours', NOW() - INTERVAL '12 days' + INTERVAL '19 hours', '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '12 days' + INTERVAL '22 hours', '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '12 days' + INTERVAL '22 hours', NULL, NOW() - INTERVAL '12 days' + INTERVAL '10 hours', NOW() - INTERVAL '12 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000020', '70000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', 'PK100020', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100020',  430.00, NOW() - INTERVAL '10 days' + INTERVAL '20 hours', NOW() - INTERVAL '10 days' + INTERVAL '22 hours', 'completed', NOW() - INTERVAL '10 days' + INTERVAL '19 hours', NOW() - INTERVAL '10 days' + INTERVAL '18 hours', '40000000-0000-0000-0000-000000000002', NULL, NULL, NOW() - INTERVAL '10 days' + INTERVAL '21 hours', '40000000-0000-0000-0000-000000000002', NOW() - INTERVAL '10 days' + INTERVAL '21 hours', NULL, NOW() - INTERVAL '10 days' + INTERVAL '10 hours', NOW() - INTERVAL '10 days' + INTERVAL '21 hours'),
    ('80000000-0000-0000-0000-000000000021', '70000000-0000-0000-0000-000000000005', '60000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', 'PK100021', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100021',  790.00, NOW() - INTERVAL '8 days' + INTERVAL '22 hours',  NOW() - INTERVAL '8 days' + INTERVAL '23 hours',  'completed', NOW() - INTERVAL '8 days' + INTERVAL '21 hours',  NOW() - INTERVAL '8 days' + INTERVAL '20 hours',  '40000000-0000-0000-0000-000000000003', NULL, NULL, NOW() - INTERVAL '8 days' + INTERVAL '22 hours',  '40000000-0000-0000-0000-000000000003', NOW() - INTERVAL '8 days' + INTERVAL '22 hours',  NULL, NOW() - INTERVAL '8 days' + INTERVAL '11 hours',  NOW() - INTERVAL '8 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000022', '70000000-0000-0000-0000-000000000007', '60000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000004', 'PK100022', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100022',  990.00, NOW() - INTERVAL '6 days' + INTERVAL '22 hours',  NOW() - INTERVAL '6 days' + INTERVAL '23 hours',  'completed', NOW() - INTERVAL '6 days' + INTERVAL '21 hours',  NOW() - INTERVAL '6 days' + INTERVAL '20 hours',  '40000000-0000-0000-0000-000000000004', NULL, NULL, NOW() - INTERVAL '6 days' + INTERVAL '22 hours',  '40000000-0000-0000-0000-000000000004', NOW() - INTERVAL '6 days' + INTERVAL '22 hours',  NULL, NOW() - INTERVAL '6 days' + INTERVAL '11 hours',  NOW() - INTERVAL '6 days' + INTERVAL '22 hours'),
    ('80000000-0000-0000-0000-000000000023', '70000000-0000-0000-0000-000000000008', '60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'PK100023', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100023',  390.00, NOW() - INTERVAL '4 days' + INTERVAL '20 hours',  NOW() - INTERVAL '4 days' + INTERVAL '22 hours',  'completed', NOW() - INTERVAL '4 days' + INTERVAL '19 hours',  NOW() - INTERVAL '4 days' + INTERVAL '18 hours',  '40000000-0000-0000-0000-000000000001', NULL, NULL, NOW() - INTERVAL '4 days' + INTERVAL '21 hours',  '40000000-0000-0000-0000-000000000001', NOW() - INTERVAL '4 days' + INTERVAL '21 hours',  NULL, NOW() - INTERVAL '4 days' + INTERVAL '10 hours',  NOW() - INTERVAL '4 days' + INTERVAL '21 hours'),
    -- Последние заказы (отменённый, суши, pending)
    ('80000000-0000-0000-0000-000000000024', '70000000-0000-0000-0000-000000000010', '60000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 'PK100024', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100024',  320.00, NOW() - INTERVAL '3 days' + INTERVAL '20 hours',  NOW() - INTERVAL '3 days' + INTERVAL '21 hours',  'cancelled', NOW() - INTERVAL '3 days' + INTERVAL '19 hours',  NULL,                                          NULL,                                   'Передумал, нашёл еду.',    NOW() - INTERVAL '3 days' + INTERVAL '15 hours',  NULL,                       NULL,                                   NULL,                       NULL,                       NOW() - INTERVAL '3 days' + INTERVAL '8 hours',   NOW() - INTERVAL '3 days' + INTERVAL '15 hours'),
    ('80000000-0000-0000-0000-000000000025', '70000000-0000-0000-0000-000000000011', '60000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000006', 'PK100025', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100025',  680.00, NOW() + INTERVAL '4 hours',  NOW() + INTERVAL '6 hours',  'confirmed', NOW() + INTERVAL '2 hours',     NOW() - INTERVAL '30 minutes', '40000000-0000-0000-0000-000000000008', NULL,                        NULL,                       NULL,                       NULL,                                   NULL,                       NULL,                       NOW() - INTERVAL '1 hour',     NOW() - INTERVAL '30 minutes'),
    ('80000000-0000-0000-0000-000000000026', '70000000-0000-0000-0000-000000000012', '60000000-0000-0000-0000-000000000011', '30000000-0000-0000-0000-000000000006', 'PK100026', 'https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=PK100026',  420.00, NOW() + INTERVAL '2 hours',  NOW() + INTERVAL '4 hours',  'paid',      NOW() + INTERVAL '1 hour',     NULL,                          NULL,                                   NULL,                        NULL,                       NULL,                       NULL,                                   NULL,                       NULL,                       NOW() - INTERVAL '20 minutes', NOW() - INTERVAL '20 minutes');

-- ─── PAYMENTS ─────────────────────────────────────────────────────────────────

INSERT INTO payments (
    id, order_id, provider_payment_id, payment_url, method, provider, amount, status, paid_at, created_at, updated_at
) VALUES
    ('90000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', 'pay_seed_0001', 'https://yookassa.ru/payments/pay_seed_0001', 'bank_card', 'yookassa',  390.00, 'succeeded', NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '10 minutes'),
    ('90000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000002', 'pay_seed_0002', 'https://yookassa.ru/payments/pay_seed_0002', 'bank_card', 'yookassa',  430.00, 'succeeded', NOW() - INTERVAL '50 minutes', NOW() - INTERVAL '50 minutes', NOW() - INTERVAL '20 minutes'),
    ('90000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000003', 'pay_seed_0003', 'https://yookassa.ru/payments/pay_seed_0003', 'sbp',       'yookassa',  790.00, 'succeeded', NOW() - INTERVAL '8 hours',    NOW() - INTERVAL '8 hours',    NOW() - INTERVAL '60 minutes'),
    ('90000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000004', 'pay_seed_0004', 'https://yookassa.ru/payments/pay_seed_0004', 'bank_card', 'yookassa',  520.00, 'succeeded', NOW() - INTERVAL '26 hours',   NOW() - INTERVAL '26 hours',   NOW() - INTERVAL '18 hours'),
    ('90000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000005', 'pay_seed_0005', 'https://yookassa.ru/payments/pay_seed_0005', 'bank_card', 'yookassa',  390.00, 'cancelled', NULL,                          NOW() - INTERVAL '49 hours',   NOW() - INTERVAL '47 hours'),
    ('90000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000006', 'pay_seed_0006', 'https://yookassa.ru/payments/pay_seed_0006', 'wallet',    'yookassa',  990.00, 'succeeded', NOW() - INTERVAL '76 hours',   NOW() - INTERVAL '76 hours',   NOW() - INTERVAL '65 hours'),
    ('90000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000007', 'pay_seed_0007', 'https://yookassa.ru/payments/pay_seed_0007', 'bank_card', 'yookassa',  790.00, 'succeeded', NOW() - INTERVAL '118 hours',  NOW() - INTERVAL '118 hours',  NOW() - INTERVAL '114 hours'),
    ('90000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000008', 'pay_seed_0008', 'https://yookassa.ru/payments/pay_seed_0008', 'sbp',       'yookassa',  520.00, 'succeeded', NOW() - INTERVAL '143 hours',  NOW() - INTERVAL '143 hours',  NOW() - INTERVAL '139 hours'),
    ('90000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000009', 'pay_seed_0009', 'https://yookassa.ru/payments/pay_seed_0009', NULL,        'yookassa',  430.00, 'pending',   NULL,                          NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),
    -- Исторические платежи
    ('90000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000010', 'pay_seed_0010', 'https://yookassa.ru/payments/pay_seed_0010', 'sbp',       'yookassa',  390.00, 'succeeded', NOW() - INTERVAL '55 days', NOW() - INTERVAL '55 days', NOW() - INTERVAL '55 days'),
    ('90000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000011', 'pay_seed_0011', 'https://yookassa.ru/payments/pay_seed_0011', 'bank_card', 'yookassa',  320.00, 'succeeded', NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days'),
    ('90000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000012', 'pay_seed_0012', 'https://yookassa.ru/payments/pay_seed_0012', 'bank_card', 'yookassa',  430.00, 'succeeded', NOW() - INTERVAL '47 days', NOW() - INTERVAL '47 days', NOW() - INTERVAL '47 days'),
    ('90000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000013', 'pay_seed_0013', 'https://yookassa.ru/payments/pay_seed_0013', 'sbp',       'yookassa',  790.00, 'succeeded', NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days'),
    ('90000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000014', 'pay_seed_0014', 'https://yookassa.ru/payments/pay_seed_0014', 'wallet',    'yookassa',  520.00, 'succeeded', NOW() - INTERVAL '42 days', NOW() - INTERVAL '42 days', NOW() - INTERVAL '42 days'),
    ('90000000-0000-0000-0000-000000000015', '80000000-0000-0000-0000-000000000015', 'pay_seed_0015', 'https://yookassa.ru/payments/pay_seed_0015', 'bank_card', 'yookassa',  790.00, 'succeeded', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days'),
    ('90000000-0000-0000-0000-000000000016', '80000000-0000-0000-0000-000000000016', 'pay_seed_0016', 'https://yookassa.ru/payments/pay_seed_0016', 'sbp',       'yookassa',  990.00, 'succeeded', NOW() - INTERVAL '38 days', NOW() - INTERVAL '38 days', NOW() - INTERVAL '38 days'),
    -- Недавние платежи
    ('90000000-0000-0000-0000-000000000017', '80000000-0000-0000-0000-000000000017', 'pay_seed_0017', 'https://yookassa.ru/payments/pay_seed_0017', 'bank_card', 'yookassa',  390.00, 'succeeded', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
    ('90000000-0000-0000-0000-000000000018', '80000000-0000-0000-0000-000000000018', 'pay_seed_0018', 'https://yookassa.ru/payments/pay_seed_0018', 'sbp',       'yookassa',  790.00, 'succeeded', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'),
    ('90000000-0000-0000-0000-000000000019', '80000000-0000-0000-0000-000000000019', 'pay_seed_0019', 'https://yookassa.ru/payments/pay_seed_0019', 'bank_card', 'yookassa',  520.00, 'succeeded', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days'),
    ('90000000-0000-0000-0000-000000000020', '80000000-0000-0000-0000-000000000020', 'pay_seed_0020', 'https://yookassa.ru/payments/pay_seed_0020', 'bank_card', 'yookassa',  430.00, 'succeeded', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'),
    ('90000000-0000-0000-0000-000000000021', '80000000-0000-0000-0000-000000000021', 'pay_seed_0021', 'https://yookassa.ru/payments/pay_seed_0021', 'wallet',    'yookassa',  790.00, 'succeeded', NOW() - INTERVAL '8 days',  NOW() - INTERVAL '8 days',  NOW() - INTERVAL '8 days'),
    ('90000000-0000-0000-0000-000000000022', '80000000-0000-0000-0000-000000000022', 'pay_seed_0022', 'https://yookassa.ru/payments/pay_seed_0022', 'sbp',       'yookassa',  990.00, 'succeeded', NOW() - INTERVAL '6 days',  NOW() - INTERVAL '6 days',  NOW() - INTERVAL '6 days'),
    ('90000000-0000-0000-0000-000000000023', '80000000-0000-0000-0000-000000000023', 'pay_seed_0023', 'https://yookassa.ru/payments/pay_seed_0023', 'bank_card', 'yookassa',  390.00, 'succeeded', NOW() - INTERVAL '4 days',  NOW() - INTERVAL '4 days',  NOW() - INTERVAL '4 days'),
    ('90000000-0000-0000-0000-000000000024', '80000000-0000-0000-0000-000000000024', 'pay_seed_0024', 'https://yookassa.ru/payments/pay_seed_0024', 'bank_card', 'yookassa',  320.00, 'refunded',  NULL,                       NOW() - INTERVAL '3 days',  NOW() - INTERVAL '3 days'),
    ('90000000-0000-0000-0000-000000000025', '80000000-0000-0000-0000-000000000025', 'pay_seed_0025', 'https://yookassa.ru/payments/pay_seed_0025', 'sbp',       'yookassa',  680.00, 'succeeded', NOW() - INTERVAL '1 hour',  NOW() - INTERVAL '1 hour',  NOW() - INTERVAL '30 minutes'),
    ('90000000-0000-0000-0000-000000000026', '80000000-0000-0000-0000-000000000026', 'pay_seed_0026', 'https://yookassa.ru/payments/pay_seed_0026', 'bank_card', 'yookassa',  420.00, 'succeeded', NOW() - INTERVAL '20 minutes', NOW() - INTERVAL '20 minutes', NOW() - INTERVAL '20 minutes');

-- ─── PAYMENT EVENTS ───────────────────────────────────────────────────────────

INSERT INTO payment_events (id, payment_id, event_type, payload, created_at) VALUES
    ('91000000-0000-0000-0000-000000000001', '90000000-0000-0000-0000-000000000001', 'payment.succeeded', '{"source":"seed","order":"80000000-0000-0000-0000-000000000001"}'::jsonb, NOW() - INTERVAL '10 minutes'),
    ('91000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000005', 'payment.cancelled', '{"source":"seed","reason":"user_cancelled"}'::jsonb,                      NOW() - INTERVAL '47 hours'),
    ('91000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000003', 'payment.succeeded', '{"source":"seed","method":"sbp"}'::jsonb,                                 NOW() - INTERVAL '8 hours'),
    ('91000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000006', 'payment.succeeded', '{"source":"seed","method":"wallet"}'::jsonb,                              NOW() - INTERVAL '76 hours'),
    ('91000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000024', 'payment.refunded',  '{"source":"seed","reason":"order_cancelled"}'::jsonb,                    NOW() - INTERVAL '3 days'),
    ('91000000-0000-0000-0000-000000000006', '90000000-0000-0000-0000-000000000013', 'payment.succeeded', '{"source":"seed","method":"sbp"}'::jsonb,                                 NOW() - INTERVAL '45 days'),
    ('91000000-0000-0000-0000-000000000007', '90000000-0000-0000-0000-000000000016', 'payment.succeeded', '{"source":"seed","method":"sbp"}'::jsonb,                                 NOW() - INTERVAL '38 days'),
    ('91000000-0000-0000-0000-000000000008', '90000000-0000-0000-0000-000000000025', 'payment.succeeded', '{"source":"seed","method":"sbp"}'::jsonb,                                 NOW() - INTERVAL '30 minutes');

-- ─── REVIEWS ──────────────────────────────────────────────────────────────────

INSERT INTO reviews (id, order_id, user_id, location_id, rating, comment, created_at, updated_at) VALUES
    ('a0000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', 5, 'Очень удобно, еда была вкусной и выдача прошла быстро.',                       NOW() - INTERVAL '55 minutes',  NOW() - INTERVAL '55 minutes'),
    ('a0000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000004', 4, 'Хороший набор, но хотелось бы чуть больше горячих блюд.',                       NOW() - INTERVAL '17 hours',    NOW() - INTERVAL '17 hours'),
    ('a0000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000003', 5, 'Отличная порция и приятный персонал.',                                          NOW() - INTERVAL '113 hours',   NOW() - INTERVAL '113 hours'),
    ('a0000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000008', '70000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000004', 3, 'Нормально для теста, но ассортимент был средний.',                              NOW() - INTERVAL '138 hours',   NOW() - INTERVAL '138 hours'),
    ('a0000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000010', '70000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000001', 5, 'Первый заказ — всё понравилось! Свежий хлеб и хорошая упаковка.',              NOW() - INTERVAL '55 days',     NOW() - INTERVAL '55 days'),
    ('a0000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000011', '70000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000001', 4, 'Десерты были немного помятые, но в целом вкусно.',                             NOW() - INTERVAL '50 days',     NOW() - INTERVAL '50 days'),
    ('a0000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000013', '70000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000003', 5, 'Ресторанный уровень за копейки. Рекомендую!',                                   NOW() - INTERVAL '45 days',     NOW() - INTERVAL '45 days'),
    ('a0000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000014', '70000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000004', 3, 'Продукты свежие, но набор был скучноватый в этот раз.',                         NOW() - INTERVAL '42 days',     NOW() - INTERVAL '42 days'),
    ('a0000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000017', '70000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 5, 'Беру уже второй раз. Стабильно хорошее качество!',                             NOW() - INTERVAL '20 days',     NOW() - INTERVAL '20 days'),
    ('a0000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000018', '70000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000003', 4, 'Тверская — удобное место. Сотрудник был приветлив.',                            NOW() - INTERVAL '15 days',     NOW() - INTERVAL '15 days'),
    ('a0000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000019', '70000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000004', 2, 'Ожидал большего. Часть продуктов была с истекающим сроком.',                   NOW() - INTERVAL '12 days',     NOW() - INTERVAL '12 days'),
    ('a0000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000022', '70000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000004', 5, 'Семейный набор — идеально для ужина на двоих. Буду брать снова.',              NOW() - INTERVAL '6 days',      NOW() - INTERVAL '6 days'),
    ('a0000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000023', '70000000-0000-0000-0000-000000000008', '30000000-0000-0000-0000-000000000001', 4, 'Хлебная пекарня — супер! Только очередь бывает в конце дня.',                  NOW() - INTERVAL '4 days',      NOW() - INTERVAL '4 days'),
    ('a0000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000021', '70000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000003', 5, 'Ужин получился очень вкусным. Ресторан не разочаровал.',                        NOW() - INTERVAL '8 days',      NOW() - INTERVAL '8 days');

-- ─── PAYOUT DESTINATIONS ──────────────────────────────────────────────────────

INSERT INTO partner_payout_destinations (partner_id, type, sbp_phone, sbp_bank_id, recipient_name, created_at, updated_at) VALUES
    ('20000000-0000-0000-0000-000000000001', 'sbp', '+79991000001', '100000000001', 'ООО Хлеб и Кофе',       NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days'),
    ('20000000-0000-0000-0000-000000000002', 'sbp', '+79991000002', '100000000004', 'ООО Вечерняя Трапеза',  NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days'),
    ('20000000-0000-0000-0000-000000000003', 'sbp', '+79991000003', '100000000007', 'ООО Сеть Городской Вкус', NOW() - INTERVAL '80 days', NOW() - INTERVAL '80 days');

-- ─── PARTNER PAYOUTS ──────────────────────────────────────────────────────────
-- Завершённые выплаты за период 60–32 дня назад

INSERT INTO partner_payouts (
    id, partner_id, period_start, period_end,
    gross_amount, commission_amount, commission_rate_applied, net_amount,
    status, provider, provider_payout_id, idempotency_key, processed_at, created_at, updated_at
) VALUES
    (
        'b0000000-0000-0000-0000-000000000001',
        '20000000-0000-0000-0000-000000000001',
        NOW() - INTERVAL '60 days', NOW() - INTERVAL '32 days',
        1140.00, 205.20, 0.1800, 934.80,
        'completed', 'yookassa', 'payout_seed_001', 'idem_seed_001',
        NOW() - INTERVAL '30 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '30 days'
    ),
    (
        'b0000000-0000-0000-0000-000000000002',
        '20000000-0000-0000-0000-000000000002',
        NOW() - INTERVAL '60 days', NOW() - INTERVAL '32 days',
        1580.00, 316.00, 0.2000, 1264.00,
        'completed', 'yookassa', 'payout_seed_002', 'idem_seed_002',
        NOW() - INTERVAL '30 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '30 days'
    ),
    (
        'b0000000-0000-0000-0000-000000000003',
        '20000000-0000-0000-0000-000000000003',
        NOW() - INTERVAL '60 days', NOW() - INTERVAL '32 days',
        1510.00, 241.60, 0.1600, 1268.40,
        'completed', 'yookassa', 'payout_seed_003', 'idem_seed_003',
        NOW() - INTERVAL '30 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '30 days'
    ),
    -- Pending payouts за последний период (32 дня назад — сегодня)
    (
        'b0000000-0000-0000-0000-000000000004',
        '20000000-0000-0000-0000-000000000001',
        NOW() - INTERVAL '32 days', NOW(),
        1210.00, 217.80, 0.1800, 992.20,
        'pending', 'yookassa', NULL, 'idem_seed_004',
        NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'
    ),
    (
        'b0000000-0000-0000-0000-000000000005',
        '20000000-0000-0000-0000-000000000002',
        NOW() - INTERVAL '32 days', NOW(),
        1580.00, 316.00, 0.2000, 1264.00,
        'pending', 'yookassa', NULL, 'idem_seed_005',
        NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'
    );

-- ─── PAYOUT ORDERS ────────────────────────────────────────────────────────────
-- partner 01, payout b0001: заказы 10, 11, 12 (390 + 320 + 430 = 1140)

INSERT INTO partner_payout_orders (payout_id, order_id, order_amount, commission_part) VALUES
    ('b0000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000010', 390.00,  70.20),
    ('b0000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000011', 320.00,  57.60),
    ('b0000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000012', 430.00,  77.40),
    -- partner 02, payout b0002: заказы 13, 15 (790 + 790 = 1580)
    ('b0000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000013', 790.00, 158.00),
    ('b0000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000015', 790.00, 158.00),
    -- partner 03, payout b0003: заказы 14, 16 (520 + 990 = 1510)
    ('b0000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000014', 520.00,  83.20),
    ('b0000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000016', 990.00, 158.40),
    -- partner 01, payout b0004 (pending): заказы 17, 20, 23 (390 + 430 + 390 = 1210)
    ('b0000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000017', 390.00,  70.20),
    ('b0000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000020', 430.00,  77.40),
    ('b0000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000023', 390.00,  70.20),
    -- partner 02, payout b0005 (pending): заказы 18, 21 (790 + 790 = 1580)
    ('b0000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000018', 790.00, 158.00),
    ('b0000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000021', 790.00, 158.00);

-- ─── ДОПОЛНИТЕЛЬНЫЕ ПАРТНЁРЫ ДЛЯ КАРТЫ ──────────────────────────────────────

INSERT INTO partners (id, brand_name, logo_url, parent_partner_id, account_type, commission_rate, promo_commission_rate, promo_commission_until, status, created_at, updated_at) VALUES
    ('20000000-0000-0000-0000-000000000006', 'Пекарня Артизан',    'https://images.unsplash.com/photo-1517433670267-08bbd4be890f?w=800', NULL, 'independent', 0.1800, NULL, NULL, 'active', NOW() - INTERVAL '40 days', NOW() - INTERVAL '2 days'),
    ('20000000-0000-0000-0000-000000000007', 'Суши Дома',          'https://images.unsplash.com/photo-1559410545-0bdcd187e0a6?w=800',    NULL, 'independent', 0.2000, NULL, NULL, 'active', NOW() - INTERVAL '35 days', NOW() - INTERVAL '1 day'),
    ('20000000-0000-0000-0000-000000000008', 'Кофе Угол',          'https://images.unsplash.com/photo-1445116572660-236099ec97a0?w=800',  NULL, 'independent', 0.1800, 0.1200, CURRENT_DATE + INTERVAL '20 days', 'active', NOW() - INTERVAL '28 days', NOW() - INTERVAL '3 days'),
    ('20000000-0000-0000-0000-000000000009', 'Бистро 24',          'https://images.unsplash.com/photo-1552566626-52f8b828add9?w=800',    NULL, 'independent', 0.2000, NULL, NULL, 'active', NOW() - INTERVAL '50 days', NOW() - INTERVAL '1 day'),
    ('20000000-0000-0000-0000-000000000010', 'Зелёная лавка',      'https://images.unsplash.com/photo-1516594798947-e65505dbb29d?w=800',  NULL, 'independent', 0.1600, NULL, NULL, 'active', NOW() - INTERVAL '55 days', NOW() - INTERVAL '4 days'),
    ('20000000-0000-0000-0000-000000000011', 'Мясной двор',        'https://images.unsplash.com/photo-1544025162-d76694265947?w=800',    NULL, 'independent', 0.1800, NULL, NULL, 'active', NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 days'),
    ('20000000-0000-0000-0000-000000000012', 'Японский квартал',   'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=800', NULL, 'independent', 0.2000, 0.1400, CURRENT_DATE + INTERVAL '10 days', 'active', NOW() - INTERVAL '33 days', NOW() - INTERVAL '1 day'),
    ('20000000-0000-0000-0000-000000000013', 'Пирожковая №1',      'https://images.unsplash.com/photo-1483695028939-5bb13f8648b0?w=800', NULL, 'independent', 0.1800, NULL, NULL, 'active', NOW() - INTERVAL '60 days', NOW() - INTERVAL '3 days'),
    ('20000000-0000-0000-0000-000000000014', 'Смузи Бар',          'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=800', NULL, 'independent', 0.1800, NULL, NULL, 'active', NOW() - INTERVAL '18 days', NOW() - INTERVAL '1 day'),
    ('20000000-0000-0000-0000-000000000015', 'Восточный базар',    'https://images.unsplash.com/photo-1424847651672-bf20a4b0982b?w=800', NULL, 'independent', 0.1600, NULL, NULL, 'active', NOW() - INTERVAL '45 days', NOW() - INTERVAL '2 days');

INSERT INTO partner_legal_info (
    partner_id, inn, ogrn, kpp, legal_address, legal_name,
    verification_status, verification_comment, verified_by, verified_at, created_at, updated_at
) VALUES
    ('20000000-0000-0000-0000-000000000006', '7702000006', '1027702000006', '770206001', 'Москва, ул. Профсоюзная, 12',     'ООО Пекарня Артизан',   'pending', NULL, NULL, NULL, NOW() - INTERVAL '40 days', NOW() - INTERVAL '10 days'),
    ('20000000-0000-0000-0000-000000000007', '7702000007', '1027702000007', '770207001', 'Москва, ул. Дмитровская, 5',      'ООО Суши Дома',         'pending', NULL, NULL, NULL, NOW() - INTERVAL '35 days', NOW() - INTERVAL '8 days'),
    ('20000000-0000-0000-0000-000000000008', '7702000008', '1027702000008', '770208001', 'Москва, ул. Красноармейская, 7',  'ООО Кофе Угол',         'pending', NULL, NULL, NULL, NOW() - INTERVAL '28 days', NOW() - INTERVAL '7 days'),
    ('20000000-0000-0000-0000-000000000009', '7702000009', '1027702000009', '770209001', 'Москва, пр-т Вернадского, 88',   'ООО Бистро 24',         'verified', NULL, '01000000-0000-0000-0000-000000000001', NOW() - INTERVAL '45 days', NOW() - INTERVAL '50 days', NOW() - INTERVAL '45 days'),
    ('20000000-0000-0000-0000-000000000010', '7702000010', '1027702000010', '770210001', 'Москва, Ботаническая ул., 15',    'ООО Зелёная лавка',     'verified', NULL, '01000000-0000-0000-0000-000000000001', NOW() - INTERVAL '50 days', NOW() - INTERVAL '55 days', NOW() - INTERVAL '50 days'),
    ('20000000-0000-0000-0000-000000000011', '7702000011', '1027702000011', '770211001', 'Москва, ул. Рязанская, 3',        'ООО Мясной двор',       'pending', NULL, NULL, NULL, NOW() - INTERVAL '25 days', NOW() - INTERVAL '5 days'),
    ('20000000-0000-0000-0000-000000000012', '7702000012', '1027702000012', '770212001', 'Москва, ул. Садовническая, 20',   'ООО Японский квартал',  'pending', NULL, NULL, NULL, NOW() - INTERVAL '33 days', NOW() - INTERVAL '9 days'),
    ('20000000-0000-0000-0000-000000000013', '7702000013', '1027702000013', '770213001', 'Москва, ул. Флотская, 33',        'ООО Пирожковая',        'verified', NULL, '01000000-0000-0000-0000-000000000002', NOW() - INTERVAL '55 days', NOW() - INTERVAL '60 days', NOW() - INTERVAL '55 days'),
    ('20000000-0000-0000-0000-000000000014', '7702000014', '1027702000014', '770214001', 'Москва, ул. Арбат, 44',           'ООО Смузи Бар',         'pending', NULL, NULL, NULL, NOW() - INTERVAL '18 days', NOW() - INTERVAL '4 days'),
    ('20000000-0000-0000-0000-000000000015', '7702000015', '1027702000015', '770215001', 'Москва, ул. Братиславская, 8',    'ООО Восточный базар',   'pending', NULL, NULL, NULL, NOW() - INTERVAL '45 days', NOW() - INTERVAL '11 days');

INSERT INTO locations (
    id, partner_id, category_code, name, address, location, phone,
    logo_url, cover_image_url, gallery_urls, working_hours, status, created_at, updated_at
) VALUES
    (
        '30000000-0000-0000-0000-000000000008',
        '20000000-0000-0000-0000-000000000006', 'bakery',
        'Пекарня Артизан на Профсоюзной', 'Москва, ул. Профсоюзная, 14',
        ST_SetSRID(ST_MakePoint(37.5830, 55.7340), 4326), '+74951000008',
        'https://images.unsplash.com/photo-1517433670267-08bbd4be890f?w=800',
        'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200',
        ARRAY['https://images.unsplash.com/photo-1483695028939-5bb13f8648b0?w=1200'],
        '{"mon":{"open":"08:00","close":"21:00"},"tue":{"open":"08:00","close":"21:00"},"wed":{"open":"08:00","close":"21:00"},"thu":{"open":"08:00","close":"21:00"},"fri":{"open":"08:00","close":"22:00"},"sat":{"open":"09:00","close":"22:00"},"sun":{"open":"09:00","close":"20:00"}}'::jsonb,
        'active', NOW() - INTERVAL '38 days', NOW() - INTERVAL '2 days'
    ),
    (
        '30000000-0000-0000-0000-000000000009',
        '20000000-0000-0000-0000-000000000007', 'restaurant',
        'Суши Дома на Дмитровской', 'Москва, ул. Дмитровская, 9',
        ST_SetSRID(ST_MakePoint(37.6480, 55.8010), 4326), '+74951000009',
        'https://images.unsplash.com/photo-1559410545-0bdcd187e0a6?w=800',
        'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200',
        ARRAY['https://images.unsplash.com/photo-1553621042-f6e147245754?w=1200'],
        '{"mon":{"open":"11:00","close":"23:00"},"tue":{"open":"11:00","close":"23:00"},"wed":{"open":"11:00","close":"23:00"},"thu":{"open":"11:00","close":"23:00"},"fri":{"open":"11:00","close":"00:00"},"sat":{"open":"12:00","close":"00:00"},"sun":{"open":"12:00","close":"22:00"}}'::jsonb,
        'active', NOW() - INTERVAL '33 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000010',
        '20000000-0000-0000-0000-000000000008', 'cafe',
        'Кофе Угол на Красноармейской', 'Москва, ул. Красноармейская, 9',
        ST_SetSRID(ST_MakePoint(37.6980, 55.7630), 4326), '+74951000010',
        'https://images.unsplash.com/photo-1445116572660-236099ec97a0?w=800',
        'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=1200',
        ARRAY['https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=1200'],
        '{"mon":{"open":"08:00","close":"21:00"},"tue":{"open":"08:00","close":"21:00"},"wed":{"open":"08:00","close":"21:00"},"thu":{"open":"08:00","close":"21:00"},"fri":{"open":"08:00","close":"22:00"},"sat":{"open":"09:00","close":"22:00"},"sun":{"open":"09:00","close":"21:00"}}'::jsonb,
        'active', NOW() - INTERVAL '26 days', NOW() - INTERVAL '3 days'
    ),
    (
        '30000000-0000-0000-0000-000000000011',
        '20000000-0000-0000-0000-000000000009', 'restaurant',
        'Бистро 24 на Вернадского', 'Москва, пр-т Вернадского, 90',
        ST_SetSRID(ST_MakePoint(37.5540, 55.7100), 4326), '+74951000011',
        'https://images.unsplash.com/photo-1552566626-52f8b828add9?w=800',
        'https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=1200',
        ARRAY['https://images.unsplash.com/photo-1424847651672-bf20a4b0982b?w=1200'],
        '{"mon":{"open":"10:00","close":"23:00"},"tue":{"open":"10:00","close":"23:00"},"wed":{"open":"10:00","close":"23:00"},"thu":{"open":"10:00","close":"23:00"},"fri":{"open":"10:00","close":"00:00"},"sat":{"open":"10:00","close":"00:00"},"sun":{"open":"11:00","close":"22:00"}}'::jsonb,
        'active', NOW() - INTERVAL '48 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000012',
        '20000000-0000-0000-0000-000000000010', 'grocery',
        'Зелёная лавка на Ботанической', 'Москва, Ботаническая ул., 17',
        ST_SetSRID(ST_MakePoint(37.6270, 55.8340), 4326), '+74951000012',
        'https://images.unsplash.com/photo-1516594798947-e65505dbb29d?w=800',
        'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200',
        ARRAY['https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200'],
        '{"mon":{"open":"08:00","close":"22:00"},"tue":{"open":"08:00","close":"22:00"},"wed":{"open":"08:00","close":"22:00"},"thu":{"open":"08:00","close":"22:00"},"fri":{"open":"08:00","close":"23:00"},"sat":{"open":"09:00","close":"23:00"},"sun":{"open":"09:00","close":"21:00"}}'::jsonb,
        'active', NOW() - INTERVAL '53 days', NOW() - INTERVAL '4 days'
    ),
    (
        '30000000-0000-0000-0000-000000000013',
        '20000000-0000-0000-0000-000000000011', 'restaurant',
        'Мясной двор на Рязанской', 'Москва, ул. Рязанская, 5',
        ST_SetSRID(ST_MakePoint(37.6640, 55.7490), 4326), '+74951000013',
        'https://images.unsplash.com/photo-1544025162-d76694265947?w=800',
        'https://images.unsplash.com/photo-1544025162-d76694265947?w=1200',
        ARRAY['https://images.unsplash.com/photo-1552566626-52f8b828add9?w=1200'],
        '{"mon":{"open":"12:00","close":"23:00"},"tue":{"open":"12:00","close":"23:00"},"wed":{"open":"12:00","close":"23:00"},"thu":{"open":"12:00","close":"23:00"},"fri":{"open":"12:00","close":"00:00"},"sat":{"open":"12:00","close":"00:00"},"sun":{"open":"13:00","close":"22:00"}}'::jsonb,
        'active', NOW() - INTERVAL '23 days', NOW() - INTERVAL '2 days'
    ),
    (
        '30000000-0000-0000-0000-000000000014',
        '20000000-0000-0000-0000-000000000012', 'restaurant',
        'Японский квартал на Садовнической', 'Москва, ул. Садовническая, 22',
        ST_SetSRID(ST_MakePoint(37.6340, 55.7440), 4326), '+74951000014',
        'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=800',
        'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200',
        ARRAY['https://images.unsplash.com/photo-1559410545-0bdcd187e0a6?w=1200'],
        '{"mon":{"open":"11:00","close":"23:00"},"tue":{"open":"11:00","close":"23:00"},"wed":{"open":"11:00","close":"23:00"},"thu":{"open":"11:00","close":"23:00"},"fri":{"open":"11:00","close":"00:00"},"sat":{"open":"12:00","close":"00:00"},"sun":{"open":"12:00","close":"22:00"}}'::jsonb,
        'active', NOW() - INTERVAL '31 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000015',
        '20000000-0000-0000-0000-000000000013', 'bakery',
        'Пирожковая №1 на Флотской', 'Москва, ул. Флотская, 35',
        ST_SetSRID(ST_MakePoint(37.5020, 55.8370), 4326), '+74951000015',
        'https://images.unsplash.com/photo-1483695028939-5bb13f8648b0?w=800',
        'https://images.unsplash.com/photo-1517686469429-8bdb88b9f907?w=1200',
        ARRAY['https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200'],
        '{"mon":{"open":"07:00","close":"21:00"},"tue":{"open":"07:00","close":"21:00"},"wed":{"open":"07:00","close":"21:00"},"thu":{"open":"07:00","close":"21:00"},"fri":{"open":"07:00","close":"22:00"},"sat":{"open":"08:00","close":"22:00"},"sun":{"open":"08:00","close":"20:00"}}'::jsonb,
        'active', NOW() - INTERVAL '58 days', NOW() - INTERVAL '3 days'
    ),
    (
        '30000000-0000-0000-0000-000000000016',
        '20000000-0000-0000-0000-000000000014', 'cafe',
        'Смузи Бар на Арбате', 'Москва, ул. Арбат, 46',
        ST_SetSRID(ST_MakePoint(37.5910, 55.7500), 4326), '+74951000016',
        'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=800',
        'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=1200',
        ARRAY['https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=1200'],
        '{"mon":{"open":"09:00","close":"21:00"},"tue":{"open":"09:00","close":"21:00"},"wed":{"open":"09:00","close":"21:00"},"thu":{"open":"09:00","close":"21:00"},"fri":{"open":"09:00","close":"22:00"},"sat":{"open":"10:00","close":"22:00"},"sun":{"open":"10:00","close":"20:00"}}'::jsonb,
        'active', NOW() - INTERVAL '16 days', NOW() - INTERVAL '1 day'
    ),
    (
        '30000000-0000-0000-0000-000000000017',
        '20000000-0000-0000-0000-000000000015', 'grocery',
        'Восточный базар на Братиславской', 'Москва, ул. Братиславская, 10',
        ST_SetSRID(ST_MakePoint(37.7740, 55.6820), 4326), '+74951000017',
        'https://images.unsplash.com/photo-1424847651672-bf20a4b0982b?w=800',
        'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200',
        ARRAY['https://images.unsplash.com/photo-1516594798947-e65505dbb29d?w=1200'],
        '{"mon":{"open":"09:00","close":"22:00"},"tue":{"open":"09:00","close":"22:00"},"wed":{"open":"09:00","close":"22:00"},"thu":{"open":"09:00","close":"22:00"},"fri":{"open":"09:00","close":"23:00"},"sat":{"open":"10:00","close":"23:00"},"sun":{"open":"10:00","close":"21:00"}}'::jsonb,
        'active', NOW() - INTERVAL '43 days', NOW() - INTERVAL '2 days'
    );

INSERT INTO partner_employees (
    id, partner_id, location_id, email, password_hash, role, name,
    must_change_password, last_login_at, created_at, updated_at
) VALUES
    ('40000000-0000-0000-0000-000000000010', '20000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000008', 'artizan.owner@berezhok.local',   '__PARTNER_PASSWORD_HASH__', 'owner', 'Наталья Белкина',   FALSE, NOW() - INTERVAL '4 hours',  NOW() - INTERVAL '38 days', NOW() - INTERVAL '4 hours'),
    ('40000000-0000-0000-0000-000000000011', '20000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000009', 'sushidoma.owner@berezhok.local', '__PARTNER_PASSWORD_HASH__', 'owner', 'Игорь Ямада',       FALSE, NOW() - INTERVAL '2 hours',  NOW() - INTERVAL '33 days', NOW() - INTERVAL '2 hours'),
    ('40000000-0000-0000-0000-000000000012', '20000000-0000-0000-0000-000000000008', '30000000-0000-0000-0000-000000000010', 'coffeeugol.owner@berezhok.local','__PARTNER_PASSWORD_HASH__', 'owner', 'Мария Лазарева',    FALSE, NOW() - INTERVAL '7 hours',  NOW() - INTERVAL '26 days', NOW() - INTERVAL '7 hours'),
    ('40000000-0000-0000-0000-000000000013', '20000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000011', 'bistro24.owner@berezhok.local',  '__PARTNER_PASSWORD_HASH__', 'owner', 'Олег Стрелков',     FALSE, NOW() - INTERVAL '1 day',   NOW() - INTERVAL '48 days', NOW() - INTERVAL '1 day'),
    ('40000000-0000-0000-0000-000000000014', '20000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000012', 'zelenaya.owner@berezhok.local',  '__PARTNER_PASSWORD_HASH__', 'owner', 'Анастасия Гурова',  FALSE, NOW() - INTERVAL '5 hours',  NOW() - INTERVAL '53 days', NOW() - INTERVAL '5 hours'),
    ('40000000-0000-0000-0000-000000000015', '20000000-0000-0000-0000-000000000011', '30000000-0000-0000-0000-000000000013', 'myasnoy.owner@berezhok.local',   '__PARTNER_PASSWORD_HASH__', 'owner', 'Виктор Пашков',     FALSE, NOW() - INTERVAL '3 hours',  NOW() - INTERVAL '23 days', NOW() - INTERVAL '3 hours'),
    ('40000000-0000-0000-0000-000000000016', '20000000-0000-0000-0000-000000000012', '30000000-0000-0000-0000-000000000014', 'yaponka.owner@berezhok.local',   '__PARTNER_PASSWORD_HASH__', 'owner', 'Дина Карасёва',     FALSE, NOW() - INTERVAL '6 hours',  NOW() - INTERVAL '31 days', NOW() - INTERVAL '6 hours'),
    ('40000000-0000-0000-0000-000000000017', '20000000-0000-0000-0000-000000000013', '30000000-0000-0000-0000-000000000015', 'pirozhki.owner@berezhok.local',  '__PARTNER_PASSWORD_HASH__', 'owner', 'Людмила Фролова',   FALSE, NOW() - INTERVAL '8 hours',  NOW() - INTERVAL '58 days', NOW() - INTERVAL '8 hours'),
    ('40000000-0000-0000-0000-000000000018', '20000000-0000-0000-0000-000000000014', '30000000-0000-0000-0000-000000000016', 'smoothie.owner@berezhok.local',  '__PARTNER_PASSWORD_HASH__', 'owner', 'Кирилл Агеев',      FALSE, NOW() - INTERVAL '2 hours',  NOW() - INTERVAL '16 days', NOW() - INTERVAL '2 hours'),
    ('40000000-0000-0000-0000-000000000019', '20000000-0000-0000-0000-000000000015', '30000000-0000-0000-0000-000000000017', 'vostok.owner@berezhok.local',    '__PARTNER_PASSWORD_HASH__', 'owner', 'Зоя Хасанова',      FALSE, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '43 days', NOW() - INTERVAL '12 hours');

INSERT INTO surprise_boxes (
    id, location_id, name, description, original_price, discount_price,
    quantity_available, pickup_time_start, pickup_time_end, image_url, status, created_at, updated_at
) VALUES
    ('60000000-0000-0000-0000-000000000013', '30000000-0000-0000-0000-000000000008', 'Хлебный микс Артизан',       'Ремесленный хлеб, круассаны и булочки конца дня.',          750.00,  330.00, 5, '19:30', '21:00', 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200', 'active', NOW() - INTERVAL '35 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000014', '30000000-0000-0000-0000-000000000008', 'Сладкий бокс пекарни',       'Торты-нарезки, пирожные и маффины.',                        680.00,  290.00, 3, '20:30', '21:30', 'https://images.unsplash.com/photo-1517433670267-08bbd4be890f?w=1200', 'active', NOW() - INTERVAL '30 days', NOW() - INTERVAL '2 days'),
    ('60000000-0000-0000-0000-000000000015', '30000000-0000-0000-0000-000000000009', 'Суши-сет на вечер',          'Ассорти роллов и суши от шефа, остатки сервиса.',           1400.00, 590.00, 4, '21:30', '22:45', 'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200', 'active', NOW() - INTERVAL '28 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000016', '30000000-0000-0000-0000-000000000009', 'Лёгкий ланч суши',           'Роллы и мисо-суп — компактный набор на одного.',            850.00,  370.00, 6, '14:00', '15:30', 'https://images.unsplash.com/photo-1559410545-0bdcd187e0a6?w=1200',    'active', NOW() - INTERVAL '25 days', NOW() - INTERVAL '3 hours'),
    ('60000000-0000-0000-0000-000000000017', '30000000-0000-0000-0000-000000000010', 'Кофейный сет вечера',        'Кофе, чай и десерты из витрины — за час до закрытия.',      620.00,  250.00, 8, '20:00', '21:00', 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=1200', 'active', NOW() - INTERVAL '24 days', NOW() - INTERVAL '2 hours'),
    ('60000000-0000-0000-0000-000000000018', '30000000-0000-0000-0000-000000000011', 'Бистро-бокс ужин',           'Горячее блюдо дня, салат и хлеб.',                          1100.00, 480.00, 5, '21:00', '22:30', 'https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=1200', 'active', NOW() - INTERVAL '45 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000019', '30000000-0000-0000-0000-000000000011', 'Обеденный остаток',          'Два блюда из меню бизнес-ланча — конец дня.',               870.00,  360.00, 4, '15:30', '17:00', 'https://images.unsplash.com/photo-1552566626-52f8b828add9?w=1200',    'active', NOW() - INTERVAL '40 days', NOW() - INTERVAL '4 hours'),
    ('60000000-0000-0000-0000-000000000020', '30000000-0000-0000-0000-000000000012', 'Эко-корзина с фермы',        'Овощи, зелень, молочные продукты с коротким сроком.',       1300.00, 540.00, 6, '19:00', '21:00', 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200',    'active', NOW() - INTERVAL '50 days', NOW() - INTERVAL '2 days'),
    ('60000000-0000-0000-0000-000000000021', '30000000-0000-0000-0000-000000000013', 'Стейк вечером',              'Стейк + гарнир из остатков заготовок шефа.',                2200.00, 950.00, 2, '22:00', '23:00', 'https://images.unsplash.com/photo-1544025162-d76694265947?w=1200',    'active', NOW() - INTERVAL '20 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000022', '30000000-0000-0000-0000-000000000014', 'Суши-сет Японский',          'Авторские роллы шеф-повара, остатки вечернего сервиса.',    1700.00, 720.00, 3, '22:00', '23:15', 'https://images.unsplash.com/photo-1617196034183-421b4040ed20?w=1200', 'active', NOW() - INTERVAL '28 days', NOW() - INTERVAL '2 hours'),
    ('60000000-0000-0000-0000-000000000023', '30000000-0000-0000-0000-000000000015', 'Пирожки ассорти',            'Печёные и жареные пирожки с разными начинками.',            520.00,  210.00, 10, '19:00', '20:30', 'https://images.unsplash.com/photo-1483695028939-5bb13f8648b0?w=1200', 'active', NOW() - INTERVAL '55 days', NOW() - INTERVAL '1 day'),
    ('60000000-0000-0000-0000-000000000024', '30000000-0000-0000-0000-000000000015', 'Домашняя выпечка',           'Пироги, кулебяки и сдобные булки по рецептам бабушки.',     780.00,  330.00, 5, '20:30', '21:45', 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=1200', 'active', NOW() - INTERVAL '48 days', NOW() - INTERVAL '3 hours'),
    ('60000000-0000-0000-0000-000000000025', '30000000-0000-0000-0000-000000000016', 'Смузи & снеки набор',        'Смузи, ореховые батончики и фруктовые нарезки.',            590.00,  240.00, 7, '20:00', '21:00', 'https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=1200', 'active', NOW() - INTERVAL '14 days', NOW() - INTERVAL '1 hour'),
    ('60000000-0000-0000-0000-000000000026', '30000000-0000-0000-0000-000000000017', 'Восточный микс',             'Хумус, фалафель, лаваш и маринованные овощи.',              980.00,  420.00, 6, '20:30', '22:00', 'https://images.unsplash.com/photo-1424847651672-bf20a4b0982b?w=1200', 'active', NOW() - INTERVAL '40 days', NOW() - INTERVAL '2 days'),
    ('60000000-0000-0000-0000-000000000027', '30000000-0000-0000-0000-000000000017', 'Специи и сухофрукты',        'Набор орехов, специй и сухофруктов с истекающим сроком.',   760.00,  310.00, 4, '21:00', '22:30', 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=1200',    'active', NOW() - INTERVAL '35 days', NOW() - INTERVAL '5 hours');
