package handlers

import (
	"io"
	"net/http"
	"time"

	. "com.jadud.search.six/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Broken out for testing.
func putEnqueueHandler(queue string, msg JSON) {
	//log.Println("Queueing", domain)
	// It is polite to ask for a new queue.
	// The library will protect us if we don't.
	Q.NewQueue(queue)
	Q.Enqueue(queue, msg)
}

func PutEnqueueHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	queue := chi.URLParam(r, "queue")
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		panic("COULD NOT READ BODY")
	}
	putEnqueueHandler(queue, msg)
	duration := time.Since(start)
	render.DefaultResponder(w, r, render.M{"result": "ok", "elapsed": duration, "msg": string(msg)})
}
