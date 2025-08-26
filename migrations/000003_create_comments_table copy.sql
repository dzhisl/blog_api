CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    blog_id INTEGER NOT NULL REFERENCES blogs(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    content TEXT NOT NULL,
    replied_to_comment_id INTEGER,
    comment_status TEXT NOT NULL
);