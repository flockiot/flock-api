CREATE TABLE IF NOT EXISTS schema_check (
    id integer PRIMARY KEY DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT schema_check_singleton CHECK (id = 1)
);

INSERT INTO schema_check (id) VALUES (1) ON CONFLICT DO NOTHING;
