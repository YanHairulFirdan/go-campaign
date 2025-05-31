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