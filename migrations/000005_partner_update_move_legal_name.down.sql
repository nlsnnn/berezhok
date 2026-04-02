ALTER TABLE partners
    ADD COLUMN legal_name VARCHAR(200);

ALTER TABLE partners
    ALTER COLUMN brand_name DROP NOT NULL;

ALTER TABLE partner_legal_info
    DROP COLUMN legal_name;

ALTER TABLE locations
    DROP CONSTRAINT IF EXISTS locations_status_check;

ALTER TABLE locations
    ADD CONSTRAINT locations_status_check
    CHECK (status IN ('active', 'inactive', 'closed'));
