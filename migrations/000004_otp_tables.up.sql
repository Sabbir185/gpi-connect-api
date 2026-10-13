-- Write your UP migration SQL here
CREATE TABLE IF NOT EXISTS otps(
    id BIGSERIAL PRIMARY KEY,
    identity VARCHAR(100) NOT NULL,
    otp VARCHAR(6) NOT NULL,
    type VARCHAR(30) NOT NULL,
    channel VARCHAR(30) NOT NULL,
    attempt INT NOT NULL DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_otps_identity_type ON otps(identity, type);
