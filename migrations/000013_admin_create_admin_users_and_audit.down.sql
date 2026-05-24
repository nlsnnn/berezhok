ALTER TABLE partner_applications
    DROP COLUMN IF EXISTS reviewed_by;

DROP TABLE IF EXISTS admin_audit_log;
DROP TABLE IF EXISTS admin_users;
