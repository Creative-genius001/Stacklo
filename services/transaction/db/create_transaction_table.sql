CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS transactions  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    wallet_id UUID NOT NULL,
    currency currency_enum NOT NULL,
    amount DECIMAL(18, 8) NOT NULL DEFAULT 0.00,
    reason VARCHAR(255),
    entry_type entry_type_enum NOT NULL,
    status status_enum NOT NULL,
    transaction_type transaction_type_enum NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fiat_transaction (
    id UUID PRIMARY KEY REFERENCES transactions(fiat_details_id),
    reference_id VARCHAR(255) NOT NULL UNIQUE,
    bank_name VARCHAR(255),
    account_name VARCHAR(255),
    account_number VARCHAR(50),
    sender_name VARCHAR(255),   
    recipient_name VARCHAR(255),
    fee DECIMAL(18, 8),         
    net_amount DECIMAL(18, 8)
);


CREATE TABLE IF NOT EXISTS crypto_transaction (
    id UUID PRIMARY KEY REFERENCES transactions(crypto_details_id),
    exchange_order_id VARCHAR(255),
    blockchain_transaction_hash VARCHAR(255) UNIQUE,
    sender_address VARCHAR(255),  
    receiver_address VARCHAR(255),
    network VARCHAR(50),       
    network_fee DECIMAL(18, 8),          
    price_at_transaction DECIMAL(18, 8),
    quote_currency_amount DECIMAL(18, 8) 
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_crypto_transaction_blockchain_hash ON crypto_transaction (blockchain_transaction_hash);

CREATE TYPE transaction_type_enum ENUM ('CRYPTO', 'FIAT');
CREATE TYPE entry_type_enum ENUM ('CREDIT', 'DEBIT');
CREATE TYPE currency_enum AS ENUM ('BTC', 'ETH', 'USDT', 'NGN');
CREATE TYPE status_enum AS ENUM ('FAILED', 'SUCCESS', 'PENDING', 'PROCESSING', 'REVERSED');

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_type ON transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_fiat_details_gateway_ref ON fiat_transaction (reference_id);