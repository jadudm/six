package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	GTST "com.jadud.search.six/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type jm = map[string]interface{}

func toMap(msg GTST.JSON) jm {
	var data interface{}
	json.Unmarshal(msg, &data)
	return data.(jm)
}

// Broken out for testing.
func getDequeueHandler(queue string) (jm, error) {
	if Q.Length(queue) > 0 {
		msg := toMap(Q.Dequeue(queue))
		return msg, nil
	} else {
		msg := toMap([]byte("{}"))
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
		render.DefaultResponder(w, r, render.M{
			"result":  "ok",
			"domain":  msg["domain"],
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
