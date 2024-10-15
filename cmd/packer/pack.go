package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"

	mdb "com.jadud.search.six/cmd/packer/mdb"
	obj "com.jadud.search.six/pkg/object-storage"
	gtst "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
	"github.com/klauspost/compress/zstd"
	_ "github.com/mattn/go-sqlite3"
)

func safe_string(o interface{}) string {
	if o == nil {
		return "SAFE-NIL"
	} else {
		return o.(string)
	}
}

func pack_full(buckets *obj.Buckets, m map[string]interface{}) {
	ephemeral_b := buckets.Ephemeral

	// Create the DB file, deleting if it exists.
	host := m["host"].(string)
	ctx := context.Background()
	ddl_b, err := os.ReadFile("schema.sql")
	ddl := string(ddl_b)
	if err != nil {
		log.Fatal(err)
	}

	sqlite_filename := host + ".sqlite"
	_ = os.Remove(sqlite_filename)

	db, err := sql.Open("sqlite3", sqlite_filename)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}
	queries := mdb.New(db)

	// FIXME: versioning tables...
	// Put it in the SQL directly?
	err = queries.SetVersion(ctx, "1.0")
	if err != nil {
		log.Fatal(err)
	}

	objects := ephemeral_b.ListObjects("indexed/" + host)
	for _, o := range objects {
		json_b := ephemeral_b.GetObject(*o.Key)
		jsm := gtst.BytesToMap(json_b)
		log.Println(jsm)

		// host, path, title, text)
		qp := mdb.CreateHtmlEntryParams{
			Host:  safe_string(jsm["host"]),
			Path:  safe_string(jsm["path"]),
			Title: "",
			Text:  safe_string(jsm["content"]),
		}

		inserted, err := queries.CreateHtmlEntry(ctx, qp)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("PACKER INSERTED", inserted)
	}

}

// https://pkg.go.dev/github.com/klauspost/compress/zstd#section-readme
func compress(in io.Reader, out io.Writer) error {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		enc.Close()
		log.Fatal(err)
	}
	return enc.Close()
}

func copy_db_to_s3(buckets *obj.Buckets, m map[string]interface{}) {
	databases_b := buckets.Databases
	host := m["host"].(string)
	sqlite_filename := host + ".sqlite"
	compressed_filename := sqlite_filename + ".zstd"
	//FIXME catch these
	in_reader, _ := os.Open(sqlite_filename)
	out_writer, _ := os.Create(compressed_filename)
	compress(in_reader, out_writer)
	in_reader.Close()
	out_writer.Close()
	key_path := []string{"sqlite", sqlite_filename}
	databases_b.UploadFile(key_path, sqlite_filename)
	key_path = []string{"zstd", compressed_filename}
	databases_b.UploadFile(key_path, compressed_filename)
	os.Remove(sqlite_filename)
	os.Remove(compressed_filename)
}

func Pack(vcap_services *vcap.VcapServices, buckets *obj.Buckets, ch_msg <-chan gtst.JSON) {
	//qs := queueing.NewQueueServer("queue-server", vcap_services)

	for {
		msg := <-ch_msg
		m := gtst.BytesToMap(msg)

		// We will get a domain.
		// Walk that domain key in the S3 bucket, and stuff everything into an SQLite DB.
		// Put the DB back into S3
		if m["result"] != nil && m["result"] != "error" {
			if m["type"] == "pack_full" {
				pack_full(buckets, m)
				copy_db_to_s3(buckets, m)
			}
		}
	}
}
