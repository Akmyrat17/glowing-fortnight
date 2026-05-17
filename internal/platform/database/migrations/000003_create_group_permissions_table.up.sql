CREATE TABLE IF NOT EXISTS
    group_permissions (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL UNIQUE,
        permission_ids INTEGER[] NOT NULL DEFAULT '{}',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_group_permissions_name ON group_permissions (name);

ALTER TABLE users
ADD CONSTRAINT fk_users_group_permission FOREIGN KEY (group_permission_id) REFERENCES group_permissions (id) ON DELETE SET NULL;