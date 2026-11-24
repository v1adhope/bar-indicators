CREATE TABLE IF NOT EXISTS weather_runs (
    id         INTEGER PRIMARY KEY,
    created_at TEXT NOT NULL,
    city       TEXT NOT NULL,
    status     TEXT NOT NULL,
    result     TEXT NOT NULL,
    CONSTRAINT ck_weather_runs_status
        CHECK (status IN ('success', 'error')),
    CONSTRAINT ck_weather_runs_result
        CHECK (status != 'success' OR json_valid(result))
);

CREATE INDEX IF NOT EXISTS weather_runs_latest
    ON weather_runs (city, status, id DESC);
