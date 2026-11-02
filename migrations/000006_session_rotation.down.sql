-- Write your DOWN migration SQL here
DROP INDEX IF EXISTS idx_otps_identity_created_at;
DROP INDEX IF EXISTS idx_devices_push_token;
DROP INDEX IF EXISTS idx_sessions_device_id;
DROP INDEX IF EXISTS idx_sessions_previous_refresh_token_hash;

ALTER TABLE sessions
    DROP COLUMN IF EXISTS rotated_at,
    DROP COLUMN IF EXISTS previous_refresh_token_hash;
