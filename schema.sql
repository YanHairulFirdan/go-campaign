-- users table
CREATE TABLE IF NOT EXISTS users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    -- id is the primary key and will auto-increment
    -- name is a unique field for user identification
    -- email is a unique field for user contact

);
-- end of users table


-- campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    target_amount DECIMAL(10, 2) NOT NULL,
    current_amount DECIMAL(10, 2) DEFAULT 0.00,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status INT NOT NULL DEFAULT 0, -- 0: draft, 1: active, 2: completed, 3: cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- add unique constraint for user_id, slug and deleted_at combination (partial index)
CREATE UNIQUE INDEX idx_campaigns_user_slug_deleted_at
ON campaigns (user_id, slug)
WHERE deleted_at IS NULL;

-- add index for status
CREATE INDEX idx_campaigns_status ON campaigns (status);

-- add index for start_date
CREATE INDEX idx_campaigns_start_date ON campaigns (start_date);

-- add index for end_date
CREATE INDEX idx_campaigns_end_date ON campaigns (end_date);
-- end of campaigns table


-- start of donaturs table
CREATE TABLE IF NOT EXISTS donaturs (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    campaign_id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);
-- end of donaturs table

-- donations table
CREATE TABLE IF NOT EXISTS donations (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    donatur_id INT NOT NULL,
    campaign_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    note TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(donatur_id) REFERENCES donaturs(id) ON DELETE CASCADE,
    FOREIGN KEY(campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

-- add index for created_at
CREATE INDEX idx_donations_created_at ON donations (created_at);
-- add index for updated_at
CREATE INDEX idx_donations_updated_at ON donations (updated_at);
-- add index foreign key donatur_id
CREATE INDEX idx_donations_donatur_id ON donations (donatur_id);
-- add index foreign key campaign_id
CREATE INDEX idx_donations_campaign_id ON donations (campaign_id);
-- end of donations table

-- payments table
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
CREATE INDEX idx_payments_status ON payments (status);
-- add index for created_at
CREATE INDEX idx_payments_created_at ON payments (created_at);
-- add index for updated_at
CREATE INDEX idx_payments_updated_at ON payments (updated_at);
-- add index foreign key donatur_id
CREATE INDEX idx_payments_donatur_id ON payments (donatur_id);
-- add index foreign key campaign_id
CREATE INDEX idx_payments_campaign_id ON payments (campaign_id);
-- -- end of payments table
