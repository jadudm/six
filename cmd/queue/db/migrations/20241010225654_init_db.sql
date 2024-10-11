-- migrate:up

-- GORM assumes QueueJobs as a model will be mapped to
-- queue_jobs as a table name. Convention is easy.
CREATE TABLE IF NOT EXISTS queue_jobs
(
    uuid UUID PRIMARY KEY,
    time_inserted      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payload            JSON
);

CREATE INDEX time_inserted_idx
    ON queue_jobs (time_inserted ASC);

-- migrate:down
DROP TABLE IF EXISTS queue_jobs;
DROP INDEX IF EXISTS time_inserted_idx;