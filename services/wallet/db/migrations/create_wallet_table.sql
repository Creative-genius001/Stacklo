CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    currency VARCHAR(25) NOT NULL,
    balance DECIMAL(18, 8) NOT NULL DEFAULT 0.00,
    virtual_account_name VARCHAR(255),
    virtual_account_number VARCHAR(50),
    virtual_bank_name VARCHAR(255),
    active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
