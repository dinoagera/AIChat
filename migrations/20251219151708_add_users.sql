-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    pass_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- 1. Бригады
CREATE TABLE brigades (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    lat DECIMAL(9,6) NOT NULL,
    lon DECIMAL(9,6) NOT NULL,
    status VARCHAR(20) DEFAULT 'free',  -- free, busy
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Заявки (история + текущие)
CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    lat DECIMAL(9,6) NOT NULL,
    lon DECIMAL(9,6) NOT NULL,
    address TEXT,
    priority VARCHAR(20) DEFAULT 'normal',
    assigned_brigade_id INTEGER REFERENCES brigades(id),
    eta_minutes INTEGER,  -- расчётное время прибытия
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS refresh_tokens (
    token TEXT PRIMARY KEY,          
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users
DROP TABLE IF EXISTS brigades
DROP TABLE IF EXISTS requests
