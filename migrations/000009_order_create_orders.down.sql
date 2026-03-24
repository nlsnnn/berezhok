DROP TABLE IF EXISTS orders CASCADE;
DROP TYPE IF EXISTS order_status;

DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_box_id;
DROP INDEX IF EXISTS idx_orders_location_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_pickup_time;
DROP INDEX IF EXISTS idx_orders_confirmation_deadline;
DROP INDEX IF EXISTS idx_orders_customer_status;

DROP TRIGGER IF EXISTS set_orders_updated_at ON orders;
