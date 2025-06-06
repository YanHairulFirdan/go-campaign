CREATE TABLE IF NOT EXISTS donations (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    donatur_id INT NOT NULL,
    campaign_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    note TEXT NULL,
    payment_status INT NOT NULL DEFAULT 0, -- 0: pending, 1: processed, 2: paid, 3: failed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(donatur_id) REFERENCES donaturs(id) ON DELETE CASCADE,
    FOREIGN KEY(campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);
-- add index for payment_status
CREATE INDEX idx_donations_payment_status ON donations (payment_status);
-- add index for created_at
CREATE INDEX idx_donations_created_at ON donations (created_at);
-- add index for updated_at
CREATE INDEX idx_donations_updated_at ON donations (updated_at);
-- add index foreign key donatur_id
CREATE INDEX idx_donations_donatur_id ON donations (donatur_id);
-- add index foreign key campaign_id
CREATE INDEX idx_donations_campaign_id ON donations (campaign_id);
-- -- end of donations table
