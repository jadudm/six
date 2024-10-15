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

func run_searcher(vcap_services *vcap.VcapServices, r *chi.Mux) {

	ch_a := make(chan gsts.JSON)
	ch_b := make(chan gsts.JSON)
	ch_c := make(chan []string)

	buckets := obj.InitBuckets(vcap_services)
	//tlp.SetupSearchRoutes(r)

	var wg sync.WaitGroup
	wg.Add(1)

	go tlp.CheckQueue(vcap_services, "SEARCH", "@every 1m", ch_a)
	// HeadCheck eats anything that doesn't return a 200
	go tlp.HeadCheck(ch_a, ch_b)
	go tlp.CopyDatabases(vcap_services, buckets, ch_b, ch_c)
	// ch_c is the list of databases copied in.
	go tlp.ServeSearch(vcap_services, buckets, ch_c, r)

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
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// Serve up the search page
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	go tlp.HealthCheck("searcher", vcap_services, r)

	log.Println("running searcher")
	go run_searcher(vcap_services, r)
	wg.Wait()

	log.Println("we will never see this")
}
