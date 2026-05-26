-- 26.05.2026 12:00:00

ALTER TABLE surprise_boxes DROP CONSTRAINT IF EXISTS valid_price;
ALTER TABLE surprise_boxes ADD CONSTRAINT valid_price CHECK (
    discount_price >= 0
    AND (original_price IS NULL OR original_price > discount_price)
);
