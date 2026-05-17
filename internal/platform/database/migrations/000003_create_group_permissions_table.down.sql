ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_group_permission;

DROP INDEX IF EXISTS idx_group_permissions_name;

DROP TABLE IF EXISTS group_permissions;