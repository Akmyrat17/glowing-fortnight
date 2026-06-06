CREATE TABLE
    experiences (
        id UUID PRIMARY KEY,
        company TEXT NOT NULL,
        position TEXT NOT NULL,
        description TEXT,
        start_date TIMESTAMPTZ NOT NULL,
        end_date TIMESTAMPTZ,
        is_current BOOLEAN NOT NULL DEFAULT FALSE,
        location TEXT,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
    );