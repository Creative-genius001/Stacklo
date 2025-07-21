CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE kyc_status_enum AS ENUM ('not_started', 'pending', 'approved', 'rejected', 'resubmit_required');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone_number VARCHAR(50) UNIQUE,
    country VARCHAR(100) NOT NULL,
    isVerified BOOLEAN DEFAULT false,
    kyc_status kyc_status_enum NOT NULL DEFAULT 'not_started', 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT chk_kyc_status CHECK (kyc_status IN ('not_started', 'pending', 'approved', 'rejected', 'resubmit_required'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_isVerified ON users (isVerified);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_number ON users (phone_number);
CREATE INDEX IF NOT EXISTS idx_users_kyc_status ON users (kyc_status);
CREATE INDEX IF NOT EXISTS idx_users_country ON users (country);