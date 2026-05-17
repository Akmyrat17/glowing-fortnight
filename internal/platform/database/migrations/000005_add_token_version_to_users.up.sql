ALTER TABLE users
ADD COLUMN IF NOT EXISTS token_version INTEGER NOT NULL DEFAULT 1;

-- Keep existing tokens valid until logout by using per-user signing keys derived from the current token version.