CREATE TABLE users (
    id UUID PRIMARY KEY, // uuidv4? //
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    token_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    UNIQUE (user_id, token_hash)
);