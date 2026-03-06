DROP TABLE IF EXISTS partner_employees;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS location_categories;

DROP FUNCTION IF EXISTS set_employees_updated_at();
DROP FUNCTION IF EXISTS set_locations_updated_at();

DROP INDEX IF EXISTS idx_employees_partner;
DROP INDEX IF EXISTS idx_employees_location;
DROP INDEX IF EXISTS idx_employees_email;
DROP INDEX IF EXISTS idx_locations_partner;
DROP INDEX IF EXISTS idx_locations_category;
DROP INDEX IF EXISTS idx_locations_geography;
DROP INDEX IF EXISTS idx_locations_status;

DROP EXTENSION IF EXISTS postgis;