ALTER TABLE partners
    ALTER COLUMN brand_name SET NOT NULL;

ALTER TABLE partners
    DROP COLUMN legal_name;

ALTER TABLE partner_legal_info
ADD COLUMN legal_name VARCHAR(200) NOT NULL;

-- "Draft" статус для локаций
ALTER TABLE locations
    DROP CONSTRAINT IF EXISTS locations_status_check;

ALTER TABLE locations
    ADD CONSTRAINT locations_status_check
    CHECK (status IN ('active', 'inactive', 'closed', 'draft'));
