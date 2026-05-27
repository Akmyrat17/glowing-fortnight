CREATE TABLE IF NOT EXISTS
    projects (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR(100) NOT NULL,
        url VARCHAR(150) UNIQUE,
        repo_url VARCHAR(150),
        description VARCHAR(500),
        start_date TIMESTAMPTZ,
        end_date TIMESTAMPTZ,
        tags JSON,
        status VARCHAR CHECK (status IN ('draft', 'published', 'archived')) DEFAULT 'draft',
        project_type VARCHAR CHECK (project_type IN ('pet', 'production')) DEFAULT 'pet'
    );

CREATE UNIQUE INDEX unique_name_url_repo_url ON projects (name, url, repo_url)
WHERE
    url IS NOT NULL
    AND repo_url IS NOT NULL;