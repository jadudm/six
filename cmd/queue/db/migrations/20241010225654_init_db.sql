-- migrate:up

CREATE TABLE IF NOT EXISTS queue_jobs
(
    job_id              UUID PRIMARY KEY NOT NULl,
    time_inserted       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    domain              TEXT NOT NULL,
    page                TEXT
);

CREATE INDEX time_inserted_idx
    ON queue_jobs (time_inserted ASC);

-- migrate:down
DROP TABLE IF EXISTS queue_jobs;
DROP INDEX IF EXISTS time_inserted_idx;