-- name: CreateHtmlEntry :one
INSERT INTO html_index (host, path, title, text) 
    VALUES (?, ?, ?, ?)
    RETURNING *;

-- This could become a AS (SELECT ...) 
-- so the table has the data at creation time?
-- name: SetVersion :exec
INSERT INTO metadata (version, last_updated) 
    VALUES (?, DateTime('now'));

-- name: SearchHtmlIndex :many
SELECT * FROM html_index 
    WHERE text MATCH ?
    ORDER BY rank
    LIMIT ?;

-- https://www.sqlitetutorial.net/sqlite-full-text-search/
-- name: SearchHtmlIndexSnippets :many
SELECT path, snippet(html_index, 3, '<b>', '</b>', '...', 16)
    FROM html_index 
    WHERE text MATCH ?
    ORDER BY rank
    LIMIT ?;

