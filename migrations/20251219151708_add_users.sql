-- -- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(50) NOT NULL,
    pass_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
ALTER TABLE users 
ADD CONSTRAINT uniq_email UNIQUE (email);
-- -- +goose Down
DROP TABLE IF EXISTS users;