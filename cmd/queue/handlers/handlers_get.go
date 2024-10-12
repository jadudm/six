package handlers

import (
	"log"
	"net/http"

	TY "com.jadud.search.six/pkg/types"
	"github.com/google/uuid"
)

// Broken out for testing.
func getDequeueHandler(domain string) {
	log.Println("Queueing", domain)

	// It is polite to ask for a new queue.
	// The library will protect us if we don't.
	TheQueue.NewQueue(domain)
	TheQueue.Enqueue(domain, TY.Job{
		JobId:  uuid.NewString(),
		Domain: domain,
		Pages:  []string{},
	})

}

func GetDequeueHandler(w http.ResponseWriter, r *http.Request) {

}
