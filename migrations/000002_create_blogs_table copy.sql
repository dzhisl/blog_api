CREATE TABLE IF NOT EXISTS blogs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    tags TEXT [] NOT NULL,
    blog_status TEXT NOT NULL
);