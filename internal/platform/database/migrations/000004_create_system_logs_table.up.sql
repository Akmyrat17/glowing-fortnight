CREATE TABLE IF NOT EXISTS
    system_logs (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        level VARCHAR(10) NOT NULL,
        type VARCHAR(20) NOT NULL,
        message TEXT NOT NULL,
        context JSONB,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_system_logs_level ON system_logs (level);

CREATE INDEX IF NOT EXISTS idx_system_logs_type ON system_logs (
    type
);

CREATE INDEX IF NOT EXISTS idx_system_logs_created ON system_logs (created_at DESC);