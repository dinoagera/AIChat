-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    pass_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
ALTER TABLE users 
-- ADD CONSTRAINT uniq_email UNIQUE (email);

-- Таблица для refresh-токенов
CREATE TABLE IF NOT EXISTS refresh_tokens (
    token TEXT PRIMARY KEY,                -- сам refresh-токен (уникальная строка)
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- -- Индекс для быстрого поиска по user_id (например, при logout всех сессий)
-- CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- -- Индекс для очистки просроченных токенов (по желанию)
-- CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users