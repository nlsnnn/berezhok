DROP TABLE partner_legal_info;
DROP INDEX IF EXISTS idx_legal_info_inn;
DROP TRIGGER IF EXISTS update_legal_info_updated_at ON partner_legal_info;

DROP TABLE partners;
DROP INDEX IF EXISTS idx_partners_status;
DROP INDEX IF EXISTS idx_partners_parent;
DROP TRIGGER IF EXISTS update_partners_updated_at ON partners;

DROP FUNCTION IF EXISTS set_updated_at;