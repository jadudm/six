package handlers

import (
	"errors"
	"net/http"
	"time"

	GTST "com.jadud.search.six/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Broken out for testing.
func getDequeueHandler(queue string) (jm, error) {
	if Q.Length(queue) > 0 {
		msg := GTST.ToMap(Q.Dequeue(queue))
		return msg, nil
	} else {
		msg := GTST.ToMap([]byte("{}"))
		return msg, errors.New("EMPTY")
	}
}

func GetDequeueHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("Dequeueing")
	start := time.Now()
	queue := chi.URLParam(r, "queue")
	msg, err := getDequeueHandler(queue)

	duration := time.Since(start)
	if err == nil {
		sendMap := make(render.M, 0)
		sendMap["result"] = "ok"
		sendMap["elapsed"] = "duration"
		for k, v := range msg {
			sendMap[k] = v
		}
		render.DefaultResponder(w, r, sendMap)
	} else {
		render.DefaultResponder(w, r, render.M{
			"result":  "error",
			"domain":  "empty queue",
			"elapsed": duration,
		})

	}
}
