-- name: CreateHtmlEntry :one
INSERT INTO html_index (host, path, title, text) 
    VALUES (?, ?, ?, ?)
    RETURNING *;

-- This could become a AS (SELECT ...) 
-- so the table has the data at creation time?
-- name: SetVersion :exec
INSERT INTO metadata (version, last_updated) 
    VALUES (?, DateTime('now'));