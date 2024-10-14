package tlp

import (
	"context"
	"database/sql"
	"log"
	"os"

	"com.jadud.search.six/pkg/dyndb/mdb"
	obj "com.jadud.search.six/pkg/object-storage"
	. "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
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
		jsm := BytesToMap(json_b)
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

func copy_db(buckets *obj.Buckets, m map[string]interface{}) {
	databases_b := buckets.Databases
	host := m["host"].(string)
	sqlite_filename := host + ".sqlite"
	key_path := []string{"packed", sqlite_filename}
	databases_b.UploadFile(key_path, sqlite_filename)
}

func Pack(vcap_services *vcap.VcapServices, buckets *obj.Buckets, ch_msg <-chan JSON) {
	//qs := queueing.NewQueueServer("queue-server", vcap_services)

	for {
		msg := <-ch_msg
		m := BytesToMap(msg)

		// We will get a domain.
		// Walk that domain key in the S3 bucket, and stuff everything into an SQLite DB.
		// Put the DB back into S3
		if m["result"] != nil && m["result"] != "error" {
			if m["type"] == "pack_full" {
				pack_full(buckets, m)
				copy_db(buckets, m)
			}
		}
	}
}
