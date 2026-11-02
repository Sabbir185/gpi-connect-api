-- Write your DOWN migration SQL here
DELETE FROM otps;
ALTER TABLE otps
    DROP COLUMN IF EXISTS consumed_at,
    DROP COLUMN IF EXISTS request_id,
    ALTER COLUMN code_hash TYPE VARCHAR(6);
ALTER TABLE otps RENAME COLUMN code_hash TO otp;

DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS devices CASCADE;

ALTER TABLE users
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS profile_completed_at,
    DROP COLUMN IF EXISTS email_verified_at,
    DROP COLUMN IF EXISTS phone_verified_at,
    DROP COLUMN IF EXISTS firebase_uid,
    ADD COLUMN IF NOT EXISTS fcm_token TEXT;
