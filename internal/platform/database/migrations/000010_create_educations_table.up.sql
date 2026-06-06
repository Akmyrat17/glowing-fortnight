CREATE TABLE
    educations (
        id UUID PRIMARY KEY,
        school TEXT NOT NULL,
        degree TEXT,
        field_of_study TEXT,
        start_date TIMESTAMPTZ NOT NULL,
        end_date TIMESTAMPTZ,
        is_current BOOLEAN NOT NULL DEFAULT FALSE,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
    );