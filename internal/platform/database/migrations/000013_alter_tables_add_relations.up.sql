ALTER TABLE experiences
ADD COLUMN project_id UUID REFERENCES projects (id) ON DELETE SET NULL;

CREATE INDEX idx_experiences_project_id ON experiences (project_id);

CREATE TABLE
    project_skills (
        project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
        skill_id UUID NOT NULL REFERENCES skills (id) ON DELETE CASCADE,
        PRIMARY KEY (project_id, skill_id)
    );

CREATE INDEX idx_project_skills_skill_id ON project_skills (skill_id);