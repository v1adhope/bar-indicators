CREATE TABLE IF NOT EXISTS cryptocurrency_runs (
    id         INTEGER PRIMARY KEY,
    created_at TEXT NOT NULL,
    convert_id TEXT NOT NULL,
    coin_ids   TEXT NOT NULL,
    status     TEXT NOT NULL,
    result     TEXT NOT NULL,
    CONSTRAINT ck_cryptocurrency_runs_status
        CHECK (status IN ('success', 'error')),
    CONSTRAINT ck_cryptocurrency_runs_result
        CHECK (status != 'success' OR json_valid(result))
);

CREATE INDEX IF NOT EXISTS cryptocurrency_runs_latest
    ON cryptocurrency_runs (convert_id, coin_ids, status, id DESC);
