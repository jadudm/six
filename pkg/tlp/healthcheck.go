package tlp

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"com.jadud.search.six/pkg/vcap"
	"github.com/go-chi/chi"
)

func HealthCheck(instance_name string, vcap_services *vcap.VcapServices, r *chi.Mux) {
	SetupSearchRoutes(r)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	var wg sync.WaitGroup
	wg.Add(1)

	port := vcap_services.VCAP.Get(fmt.Sprintf(`user-provided.#(name=="%s").credentials.port`, instance_name)).String()
	log.Printf("starting %s on port %s\n", instance_name, port)
	http.ListenAndServe(":"+port, r)

	wg.Wait()
}
