-- Write your UP migration SQL here
ALTER TABLE sessions
    ADD COLUMN IF NOT EXISTS previous_refresh_token_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS rotated_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_sessions_previous_refresh_token_hash ON sessions(previous_refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_device_id ON sessions(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_push_token ON devices(push_token);
CREATE INDEX IF NOT EXISTS idx_otps_identity_created_at ON otps(identity, created_at);
