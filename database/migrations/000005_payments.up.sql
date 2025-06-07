CREATE TABLE IF NOT EXISTS payments (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    transaction_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    donatur_id INT NOT NULL,
    campaign_id INT NOT NULL,
    vendor VARCHAR(50) NOT NULL, -- e.g., 'stripe', 'paypal', 'bank'
    method VARCHAR(50) NOT NULL, -- e.g., 'credit_card', 'bank_transfer', 'paypal'
    amount DECIMAL(10, 2) NOT NULL,
    link TEXT NULL, -- URL for payment gateway or transaction details
    note TEXT NULL,
    status INT NOT NULL DEFAULT 0, -- 0: pending, 1: processed, 2: paid, 3: failed
    response JSONB NULL, -- JSON response from the payment gateway
    payment_date TIMESTAMP NULL, -- date when the payment was made
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(donatur_id) REFERENCES donaturs(id) ON DELETE CASCADE,
    FOREIGN KEY(campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

-- add index for status
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
-- add index for created_at
CREATE INDEX idx_payments_created_at ON payments (created_at);
-- add index for updated_at
CREATE INDEX idx_payments_updated_at ON payments (updated_at);
-- add index foreign key donatur_id
CREATE INDEX idx_payments_donatur_id ON payments (donatur_id);
-- add index foreign key campaign_id
CREATE INDEX idx_payments_campaign_id ON payments (campaign_id);
-- -- end of payments table
