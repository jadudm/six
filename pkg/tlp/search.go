package tlp

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	sqlitezstd "github.com/jtarchie/sqlitezstd"

	queueing "com.jadud.search.six/cmd/queue-server/pkg/queueing"
	"com.jadud.search.six/pkg/dyndb/mdb"
	sack "com.jadud.search.six/pkg/knapsack"
	obj "com.jadud.search.six/pkg/object-storage"
	gtst "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tidwall/gjson"
)

func copy_in_databases(buckets *obj.Buckets, limit int64) []string {
	// Limit is in MB; conver to KB for calculation.
	knappy := sack.NewKnapsack(limit * 1000)

	objects := buckets.Databases.ListObjects("sqlite/")
	for _, o := range objects {
		size_b := float64(*o.Size)
		size_k := math.Ceil(size_b / 1000)
		sf := &gtst.SqliteFile{
			Size: int64(size_k),
			Name: *o.Key,
		}
		knappy.Add(sf)
	}
	best_knappy := knappy.Solve()

	copied_in := make([]string, 0)
	for _, item := range best_knappy.Items {
		log.Println(item)
		log.Println(item.Id())
		zfs := strings.Split(item.Id(), "/")
		zstd_filename := zfs[len(zfs)-1]
		buckets.Databases.DownloadFile([]string{item.Id()}, zstd_filename)
		copied_in = append(copied_in, zstd_filename)
	}
	return copied_in
}

func search_handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	host := chi.URLParam(r, "host")
	log.Println("SEARCH HANDLER", host)
	sqlite_file := host + ".sqlite"
	//compressed_file := host + ".sqlite.zstd?vfs=zstd"

	// can_handle := false
	// for _, filename := range csf {
	// 	log.Println(filename, host)
	// 	if strings.HasPrefix(filename, host) {
	// 		can_handle = true
	// 		sqlite_file = filename + "?cache=shared&mode=ro"
	// 	}
	// }

	if true {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(msg))
		search_terms := gjson.GetBytes(msg, "search-terms").String()
		log.Println(search_terms)

		// Search DB
		ctx := context.Background()

		log.Println("OPENING", sqlite_file)
		db, err := sql.Open("sqlite3", sqlite_file)
		if err != nil {
			log.Fatal("CANNOT OPEN DBFILE", sqlite_file)
		}
		log.Println("WE HAVE A DATABASE")

		// log.Println("OPENING", compressed_file)
		// cdb, err := sql.Open("sqlite3", compressed_file)
		// if err != nil {
		// 	log.Fatal("CANNOT OPEN DBFILE", compressed_file)
		// }
		// log.Println("WE HAVE A DATABASE")

		// TESTING
		// _, err = db.Query("SELECT * FROM html_index WHERE text MATCH ? ORDER BY rank", search_terms)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// Set PRAGMA for each connection
		// db.SetConnMaxLifetime(0) // Disable connection pooling
		// log.Println("A")
		// db.SetMaxOpenConns(1) // Allow only one open connection
		// log.Println("B")
		// conn, err := db.Conn(ctx)
		// if err != nil {
		// 	panic(fmt.Sprintf("Failed to get connection: %s", err))
		// }
		// defer conn.Close()

		// log.Println("ABOUT TO PRAGMA")
		// _, err = db.Exec(`PRAGMA temp_store = file; PRAGMA temp_store_directory = '/tmp';`)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		queries := mdb.New(db)
		log.Println("QUERIES")

		terms := mdb.SearchHtmlIndexSnippetsParams{
			Text:  search_terms,
			Limit: 10,
		}

		res, err := queries.SearchHtmlIndexSnippets(ctx, terms)
		if err != nil {
			log.Fatal(err)
		}

		duration := time.Since(start)
		render.DefaultResponder(w, r, render.M{"result": "ok", "elapsed": duration, "results": res})
	} else {
		duration := time.Since(start)
		render.DefaultResponder(w, r, render.M{"result": "error", "elapsed": duration})

	}
}

func SetupSearchRoutes(r *chi.Mux) {
	r.Post("/search/{host}", search_handler)
}

func ServeSearch(vcap_services *vcap.VcapServices, buckets *obj.Buckets, ch_in <-chan []string, r *chi.Mux) {
	initErr := sqlitezstd.Init()
	if initErr != nil {
		panic(fmt.Sprintf("Failed to initialize SQLiteZSTD: %s", initErr))
	}
	for {
		csf := <-ch_in
		log.Println("SEARCH now serving", csf)
	}
}

func CopyDatabases(vcap_services *vcap.VcapServices,
	buckets *obj.Buckets,
	ch_in <-chan gtst.JSON,
	ch_out chan<- []string) {

	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	sleep_duration := 1 * time.Second
	min_sleep := 1 * time.Second
	max_sleep := 600 * time.Second
	retries := 0

	qs := queueing.NewQueueServer("queue-server", vcap_services)
	//FIXME: don't use instance names in all the VCAP searches
	knapsack_limit := vcap_services.VCAP.Get("user-provided.#(instance_name==\"searcher\").credentials.local-disk-limit").Int()
	my_instance_name := vcap_services.VCAP.Get("user-provided.#(instance_name==\"searcher\").instance_name").String()
	if len(my_instance_name) < 1 {
		log.Fatal("I HAVE NO INSTANCE NAME")
	}

	for {
		msg := <-ch_in

		m := gtst.BytesToMap(msg)

		// We will get a domain.
		// Walk that domain key in the S3 bucket, and stuff everything into an SQLite DB.
		// Put the DB back into S3
		if m["result"] != nil && m["result"] != "error" {
			log.Println(m)
			if m["type"] == "search" && m["search-id"] == my_instance_name {
				retries = 0
				sleep_duration = min_sleep
				copied_in := copy_in_databases(buckets, knapsack_limit)
				ch_out <- copied_in
			} else {
				// It wasn't for us.
				// Backoff and try again
				retries += 1
				sleep_duration += CalculateBackoff(retries, initialBackoff, maxBackoff)
				if sleep_duration > time.Duration(max_sleep) {
					// If we get here, we've slept for a long time, and no one has picked up
					// this message. This suggests either a wrong ID, or something bad has happened.
					// But, we shouldn't keep propagating the message. Dropping.
					log.Println("SEARCH dropping", m)
				} else {
					go func() {
						log.Println("SEARCH requeing", sleep_duration, m)
						time.Sleep(sleep_duration)
						qs.Enqueue("SEARCH", gtst.MapToBytes(m))
					}()
				}
			}
		}
	}
}
