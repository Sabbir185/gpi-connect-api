-- Write your UP migration SQL here
ALTER TABLE users
    DROP COLUMN IF EXISTS fcm_token,
    ALTER COLUMN name SET DEFAULT '',
    ALTER COLUMN country_id DROP NOT NULL,
    ADD COLUMN IF NOT EXISTS firebase_uid VARCHAR(128) UNIQUE,
    ADD COLUMN IF NOT EXISTS phone_verified_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS profile_completed_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP WITH TIME ZONE;

CREATE TABLE IF NOT EXISTS devices(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    installation_id VARCHAR(100) UNIQUE NOT NULL,
    platform VARCHAR(20) NOT NULL,
    os_version VARCHAR(50),
    app_version VARCHAR(50),
    locale VARCHAR(35),
    timezone VARCHAR(64),
    push_token TEXT,
    last_seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);

CREATE TABLE IF NOT EXISTS sessions(
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id BIGINT REFERENCES devices(id) ON DELETE SET NULL,
    refresh_token_hash VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

-- Codes are stored as HMAC digests and looked up by an unguessable request id.
DELETE FROM otps;
ALTER TABLE otps RENAME COLUMN otp TO code_hash;
ALTER TABLE otps
    ALTER COLUMN code_hash TYPE VARCHAR(64),
    ADD COLUMN IF NOT EXISTS request_id UUID UNIQUE NOT NULL,
    ADD COLUMN IF NOT EXISTS consumed_at TIMESTAMP WITH TIME ZONE;
