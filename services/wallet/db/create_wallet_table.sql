CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE wallet_type_enum AS ENUM ('CRYPTO', 'FIAT');

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    currency VARCHAR(25) NOT NULL,
    balance DECIMAL(18, 8) NOT NULL DEFAULT 0.00,
    wallet_type wallet_type_enum NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (user_id, currency)
   
);

CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);

CREATE TABLE IF NOT EXISTS fiat_wallet_metadata (
    wallet_id UUID PRIMARY KEY REFERENCES wallets(id) ON DELETE CASCADE,
    virtual_account_name VARCHAR(255),
    virtual_account_number VARCHAR(50),
    virtual_bank_name VARCHAR(255)
);
