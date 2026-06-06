DROP INDEX IF EXISTS idx_experiences_project_id;

ALTER TABLE experiences
DROP COLUMN IF EXISTS project_id;

DROP TABLE IF EXISTS project_skills;