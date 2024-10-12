package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/semaphore"
)

var _DB *sql.DB

func load_dotenv() {
	// https://github.com/joho/godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func load_db() {
	db, err := sql.Open("sqlite3", os.Getenv("DATABASE_PATH"))
	if err != nil {
		fmt.Println(err)
		panic("COULD NOT OPEN DATABASE")
	}
	_DB = db

}

// Handler functions
func enqueue(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	domain := chi.URLParam(r, "domain")
	sql, args, _ := squirrel.
		Insert("queue_jobs").Columns("job_id", "domain").
		Values(uuid.NewString(), domain).
		ToSql()

	_, err := _DB.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		panic("SQL")
	}
	duration := time.Since(start)
	render.DefaultResponder(w, r, render.M{"result": "ok", "elapsed": duration})
}

var sem = semaphore.NewWeighted(int64(1))

type Job struct {
	JobId  string `sql:"job_id"`
	Domain string `sql:"domain"`
	Page   string `sql:"page"`
}

func dequeue(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sem.Acquire(ctx, 1)
	sql, args, err := squirrel.Select("*").From("queue_jobs").OrderBy("time_inserted").Limit(1).ToSql()
	rows, err = _DB.Query(sql, args...)
	if err != nil {
		log.Println(sql)
		panic("COULD NOT QUERY QUEUE")

	}

	rows.Next()

	// http://go-database-sql.org/modifying.html
	// rows, err := _DB.Query("SELECT * FROM queue_jobs ORDER BY time_inserted ASC")
	// if err != nil {
	// 	panic("COULD NOT QUERY QUEUE")
	// }
	// Just grab the first row.
	// Next() must be called at least once.
	// stmt := QueueJobs.SELECT(QueueJobs.AllColumns).ORDER_BY(QueueJobs.TimeInserted.ASC()).LIMIT(1)
	// var qj QueueJob
	// var qjs []QueueJob
	// err = scan.Rows(&qjs, rows) //&qj.UUID, &js)
	// if err != nil {
	// 	log.Println(err)
	// 	panic("COULD NOT SCAN ROWS")
	// }

	// defer rows.Close()

	// qj = qjs[0]
	// if err == nil {
	// 	// Delete if we found no errors
	// 	//_DB.Where("uuid = ?", qj.UUID).Delete(&QueueJob{})
	// }
	sem.Release(1)
	// if err := render.Render(w, r, &qj); err != nil {
	// 	render.Render(w, r, ErrRender(err))
	// 	return
	// }
}

func main() {
	load_dotenv()
	load_db()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// This eats {domain} dots. Eg. `foo.com` becomes `foo`.
	// r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("search six"))
	})

	r.Put("/enqueue/{domain}", enqueue)
	r.Get("/dequeue", dequeue)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
