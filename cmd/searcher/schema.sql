-- Use the FTS5 module to create a full-text searchable table.
CREATE VIRTUAL TABLE html_index 
    USING fts5(host, path, title, text);

CREATE VIRTUAL TABLE pdf_index
    USING fts5(host, path, page_number, text);

CREATE TABLE metadata
    (
        version TEXT NOT NULL,
        last_updated DATE NOT NULL
    )

