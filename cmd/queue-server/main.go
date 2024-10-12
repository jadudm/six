package main

import (
	"log"
	"net/http"
	"os"

	handle "com.jadud.search.six/cmd/queue-server/handlers"

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
	//FIXME load the path from a config/env var
	handle.Init("/tmp/queue-backup", "@every 1m")

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// This router eats {domain} dots. Eg. `foo.com` becomes `foo`.
	// r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("queue-server"))
	})

	r.Put("/enqueue/{queue}", handle.PutEnqueueHandler)
	r.Get("/dequeue/{queue}", handle.GetDequeueHandler)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// FIXME: Either env var or config/vcap
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
