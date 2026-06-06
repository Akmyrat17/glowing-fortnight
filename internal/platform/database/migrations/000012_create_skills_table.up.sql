CREATE TABLE
    skills (
        id UUID PRIMARY KEY,
        name TEXT NOT NULL,
        level TEXT NOT NULL,
        category TEXT NOT NULL,
        icon_url TEXT,
        years_exp INTEGER,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
    );