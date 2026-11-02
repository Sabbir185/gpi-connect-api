-- Write your DOWN migration SQL here
ALTER TABLE users
    DROP COLUMN IF EXISTS designation_status,
    DROP COLUMN IF EXISTS designation_title,
    DROP COLUMN IF EXISTS designation_id,
    DROP COLUMN IF EXISTS country_code;

DROP TABLE IF EXISTS designations CASCADE;
