package handlers

import (
	"errors"
	"net/http"
	"time"

	. "com.jadud.search.six/pkg/types"

	"github.com/go-chi/render"
)

// Broken out for testing.
func getDequeueHandler() (Job, error) {
	if TheMultiqueue.Length(TheQueue) > 0 {
		return TheMultiqueue.Dequeue(TheQueue), nil
	} else {
		return Job{}, errors.New("EMPTY")
	}
}

func GetDequeueHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("Dequeueing")
	start := time.Now()
	job, err := getDequeueHandler()
	duration := time.Since(start)
	if err == nil {
		render.DefaultResponder(w, r, render.M{
			"result":  "ok",
			"domain":  job.Domain,
			"elapsed": duration,
		})
	} else {
		render.DefaultResponder(w, r, render.M{
			"result":  "error",
			"domain":  "empty queue",
			"elapsed": duration,
		})

	}
}
