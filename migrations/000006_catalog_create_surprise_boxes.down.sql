DROP TABLE IF EXISTS surprise_boxes;

DROP INDEX IF EXISTS idx_boxes_location_id;
DROP INDEX IF EXISTS idx_boxes_status;

DROP TRIGGER IF EXISTS set_surprise_boxes_updated_at ON surprise_boxes;
