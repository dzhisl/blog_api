CREATE TABLE IF NOT EXISTS reactions (
    user_id INTEGER NOT NULL REFERENCES users(id),
    comment_id INTEGER NOT NULL REFERENCES comments(id),
    reaction TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);