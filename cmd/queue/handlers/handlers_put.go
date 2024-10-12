package handlers

import (
	"log"
	"net/http"
	"time"

	. "com.jadud.search.six/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// Broken out for testing.
func putEnqueueHandler(domain string) {
	log.Println("Queueing", domain)

	// It is polite to ask for a new queue.
	// The library will protect us if we don't.
	TheQueue.NewQueue(domain)
	TheQueue.Enqueue(domain, Job{
		JobId:  uuid.NewString(),
		Domain: domain,
		Pages:  []string{},
	})

}

func PutEnqueueHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	domain := chi.URLParam(r, "domain")
	putEnqueueHandler(domain)
	duration := time.Since(start)
	render.DefaultResponder(w, r, render.M{"result": "ok", "elapsed": duration})

}
