-- Write your UP migration SQL here
CREATE TABLE IF NOT EXISTS designations(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(60) UNIQUE NOT NULL,
    hint VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_designations_title_lower ON designations(LOWER(title));

INSERT INTO designations(title, hint) VALUES
    ('Lecturer', 'Most used in universities'),
    ('Assistant Professor', 'Most used in universities'),
    ('Associate Professor', NULL),
    ('Professor', NULL),
    ('Teacher', 'Most used in schools'),
    ('Head Teacher', NULL),
    ('Principal', NULL),
    ('Student', NULL),
    ('Research Assistant', NULL),
    ('Pharmacist', NULL),
    ('Pharmacy Technician', NULL),
    ('Trainee Pharmacist', 'Most used in pharmaceuticals'),
    ('Medical Representative', 'Most used in pharmaceuticals'),
    ('Medical Officer', NULL),
    ('Staff Nurse', NULL),
    ('Software Engineer', NULL),
    ('Civil Engineer', NULL),
    ('Accountant', NULL),
    ('HR Manager', NULL),
    ('Manager', NULL),
    ('Assistant Manager', NULL),
    ('Executive', NULL),
    ('Officer', 'Most used in government'),
    ('Field Officer', 'Most used in NGOs'),
    ('Sales Representative', NULL),
    ('Program Coordinator', NULL)
ON CONFLICT (title) DO NOTHING;

-- country_code is the ISO 3166-1 alpha-2 country the user picked; country_id stays for supported-country settings.
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS country_code VARCHAR(2),
    ADD COLUMN IF NOT EXISTS designation_id BIGINT REFERENCES designations(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS designation_title VARCHAR(60),
    ADD COLUMN IF NOT EXISTS designation_status VARCHAR(20);
