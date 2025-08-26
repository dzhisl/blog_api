CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    password_hash TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_role TEXT NOT NULL,
    user_status TEXT NOT NULL
);