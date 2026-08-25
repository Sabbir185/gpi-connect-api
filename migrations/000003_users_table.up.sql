-- Write your UP migration SQL here
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    profile_img TEXT,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(30) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    status VARCHAR(30) NOT NULL,
    address TEXT,
    city VARCHAR(50),
    country_id BIGINT NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    post_code VARCHAR(10),
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);


CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_phone_status ON users(phone, status);
CREATE INDEX IF NOT EXISTS idx_users_id_status ON users(id, status);
CREATE INDEX IF NOT EXISTS idx_users_country_id ON users(country_id);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);


