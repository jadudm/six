package main

import (
	"log"
	"net/http"
	"os"

	handle "com.jadud.search.six/cmd/queue-server/handlers"
	vcap "com.jadud.search.six/pkg/vcap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Getenv("VCAP_SERVICES")) < 30 {
		log.Println("export VCAP_SERVICES=$(cat /app/vcap.json)")
		log.Fatal("Set VCAP_SERVICES to run the queue-server. Exiting.")
	}
	vcap_services := vcap.ReadVCAPConfig()

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

	port := vcap_services.VCAP.Get(`user-provided.#(name=="queue-server").credentials.port`).String()
	log.Printf("starting queue-server on port %s\n", port)
	http.ListenAndServe(":"+port, r)
}
