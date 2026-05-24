ALTER TABLE partner_legal_info
  ADD COLUMN verification_status VARCHAR(20) DEFAULT 'pending'
    CHECK (verification_status IN ('pending', 'verified', 'failed')),
  ADD COLUMN verification_comment TEXT,
  ADD COLUMN verified_by UUID REFERENCES admin_users(id),
  ADD COLUMN verified_at TIMESTAMP;
