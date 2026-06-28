CREATE TABLE links (
    code           TEXT        PRIMARY KEY,
    original_url   TEXT        NOT NULL,
    normalized_url TEXT        NOT NULL UNIQUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
