ALTER TABLE locations DROP CONSTRAINT IF EXISTS locations_status_check;
ALTER TABLE locations ADD CONSTRAINT locations_status_check
    CHECK (status IN ('active', 'inactive', 'closed', 'draft'));
