CREATE TABLE IF NOT EXISTS
    permissions (
        id SERIAL PRIMARY KEY,
        module VARCHAR(50) NOT NULL,
        action TEXT NOT NULL,
        name VARCHAR(100) NOT NULL UNIQUE,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions (name);

CREATE INDEX IF NOT EXISTS idx_permissions_module ON permissions (module);