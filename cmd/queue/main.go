package main

import (
	"log"
	"net/http"
	"os"

	handle "com.jadud.search.six/cmd/queue/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func load_dotenv() {
	// https://github.com/joho/godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	load_dotenv()
	r := chi.NewRouter()
	handle.Init("/tmp/queue-backup", "*/1 0 0 0 0")

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// This eats {domain} dots. Eg. `foo.com` becomes `foo`.
	// r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("queue six"))
	})

	r.Put("/enqueue/{domain}", handle.PutEnqueueHandler)
	r.Get("/dequeue", handle.GetDequeueHandler)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
