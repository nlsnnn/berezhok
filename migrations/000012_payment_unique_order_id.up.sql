DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM payments
        GROUP BY order_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot create unique payments(order_id) index: duplicate order_id values exist in payments';
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_order_id_unique
    ON payments(order_id);
