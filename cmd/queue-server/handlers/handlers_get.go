package handlers

import (
	"errors"
	"net/http"
	"time"

	gtst "com.jadud.search.six/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tidwall/gjson"
)

// Broken out for testing.
func getDequeueHandler(queue string) (gtst.JSON, error) {
	if Q.Length(queue) > 0 {
		msg := Q.Dequeue(queue)
		return msg, nil
	} else {
		msg := gtst.EmptyObject
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
		gr := gjson.ParseBytes(msg)
		sendMap := make(render.M, 0)
		sendMap["result"] = "ok"
		sendMap["elapsed"] = duration
		for _, key := range []string{"host", "path", "type"} {
			sendMap[key] = gr.Get(key).String()
		}
		render.DefaultResponder(w, r, sendMap)
	} else {
		render.DefaultResponder(w, r, render.M{
			"result":  "error",
			"host":    nil,
			"elapsed": duration,
		})

	}
}
