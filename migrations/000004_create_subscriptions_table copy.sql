CREATE TABLE IF NOT EXISTS subscriptions (
    follower_id INTEGER NOT NULL REFERENCES users(id),
    following_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL
);