DROP TABLE IF EXISTS reviews;

DROP INDEX IF EXISTS idx_reviews_location_id;
DROP INDEX IF EXISTS idx_reviews_user_id;
DROP INDEX IF EXISTS idx_reviews_rating;
DROP INDEX IF EXISTS idx_reviews_created;

DROP TRIGGER IF EXISTS set_reviews_updated_at ON reviews;
DROP FUNCTION IF EXISTS set_reviews_updated_at();
