package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	obj "com.jadud.search.six/pkg/object-storage"
	tlp "com.jadud.search.six/pkg/tlp"
	gsts "com.jadud.search.six/pkg/types"
	vcap "com.jadud.search.six/pkg/vcap"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func run_healthcheck(vcap_services *vcap.VcapServices) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	var wg sync.WaitGroup
	wg.Add(1)

	port := vcap_services.VCAP.Get(`user-provided.#(name=="indexer").credentials.port`).String()
	log.Printf("starting indexer on port %s\n", port)
	http.ListenAndServe(":"+port, r)

	wg.Wait()
}

func run_indexer(vcap_services *vcap.VcapServices) {

	ch_a := make(chan gsts.JSON)
	ch_b := make(chan gsts.JSON)

	buckets := obj.InitBuckets(vcap_services)

	var wg sync.WaitGroup
	wg.Add(1)

	go tlp.CheckQueue(vcap_services, "INDEX", "@every 1m", ch_a)
	// HeadCheck eats anything that doesn't return a 200
	go tlp.HeadCheck(ch_a, ch_b)
	go tlp.Index(buckets, ch_b)

	wg.Wait()

}

func main() {

	if len(os.Getenv("VCAP_SERVICES")) < 30 {
		log.Println("export VCAP_SERVICES=$(cat /app/vcap.json)")
		log.Fatal("Set VCAP_SERVICES to run the indexer. Exiting.")
	}
	vcap_services := vcap.ReadVCAPConfig()

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("running healthcheck")
	go run_healthcheck(vcap_services)

	log.Println("running indexer")
	go run_indexer(vcap_services)
	wg.Wait()

	log.Println("we will never see this")
}
