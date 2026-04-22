CREATE TABLE password_reset_tokens (
    email TEXT,
    token_hash TEXT,
    expires_at TIMESTAMP
);