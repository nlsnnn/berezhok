CREATE TABLE partner_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_name VARCHAR(100) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(15) NOT NULL,
    business_name VARCHAR(200) NOT NULL,
    category_code VARCHAR(50),
    address TEXT,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending' NOT NULL,
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() 
);
CREATE INDEX idx_applications_status ON partner_applications(status);