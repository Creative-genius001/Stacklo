CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS transactions  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    currency currency_enum NOT NULL,
    amount DECIMAL(18, 8) NOT NULL DEFAULT 0.00,
    status status_enum NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TYPE currency_enum AS ENUM ('BTC', 'ETH', 'USDT', 'NGN');
CREATE TYPE status_enum AS ENUM ('failed', 'success', 'pending');
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);