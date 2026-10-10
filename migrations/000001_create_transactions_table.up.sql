CREATE TYPE transaction_status AS ENUM ('pending', 'success', 'failed', 'expired');

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'IDR',
    status transaction_status DEFAULT 'pending',
    payment_method VARCHAR(50) NOT NULL,
    reference_id VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);