CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE transaction_type_enum AS ENUM ('CRYPTO', 'FIAT');
CREATE TYPE entry_type_enum AS ENUM ('CREDIT', 'DEBIT');
CREATE TYPE status_enum AS ENUM ('FAILED', 'SUCCESS', 'PENDING', 'PROCESSING', 'REVERSED'); 

CREATE TABLE IF NOT EXISTS transactions  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount DECIMAL(18, 8) NOT NULL DEFAULT 0.00,
    reason VARCHAR(255),
    entry_type entry_type_enum NOT NULL,
    status status_enum NOT NULL,
    transaction_type transaction_type_enum NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fiat_transactions (
    id UUID PRIMARY KEY REFERENCES transactions(id),
    reference_id VARCHAR(255) NOT NULL UNIQUE,
    transaction_number VARCHAR(255) NOT NULL,
    bank_name VARCHAR(255),
    account_name VARCHAR(255),
    account_number VARCHAR(50),
    fee DECIMAL(18, 8),         
    net_amount DECIMAL(18, 8)
);


CREATE TABLE IF NOT EXISTS crypto_transactions (
    id UUID PRIMARY KEY REFERENCES transactions(id),
    exchange_order_id VARCHAR(255),
    network VARCHAR(50),       
    network_fee DECIMAL(18, 8),          
    price_at_transaction DECIMAL(18, 8),
    quote_currency_amount DECIMAL(18, 8) 
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_wallet_id ON transactions(wallet_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_currency ON transactions(currency, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_type ON transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_fiat_details_gateway_ref ON fiat_transaction (reference_id);