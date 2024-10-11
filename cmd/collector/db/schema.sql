CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE VIRTUAL TABLE html_index
    USING fts5(host, path, title, text)
/* html_index(host,path,title,text) */;
CREATE TABLE IF NOT EXISTS 'html_index_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'html_index_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'html_index_content'(id INTEGER PRIMARY KEY, c0, c1, c2, c3);
CREATE TABLE IF NOT EXISTS 'html_index_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'html_index_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE VIRTUAL TABLE pdf_index
    USING fts5(host, path, page_number, text)
/* pdf_index(host,path,page_number,text) */;
CREATE TABLE IF NOT EXISTS 'pdf_index_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'pdf_index_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'pdf_index_content'(id INTEGER PRIMARY KEY, c0, c1, c2, c3);
CREATE TABLE IF NOT EXISTS 'pdf_index_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'pdf_index_config'(k PRIMARY KEY, v) WITHOUT ROWID;
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20241011011755');
