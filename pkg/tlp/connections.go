package tlp

import (
	"fmt"
	"log"
	"os"

	"com.jadud.search.six/pkg/vcap"
)

func GetQueueServerEndpoint(vcap_services *vcap.VcapServices) string {
	var endpoint string
	env := os.Getenv("ENV")
	switch env {
	case "LOCAL":
		port := vcap_services.VCAP.Get(`user-provided.#(name=="queue-server").credentials.port`).String()
		endpoint = fmt.Sprintf("http://localhost:%s", port)
	case "DOCKER":
		endpoint = vcap_services.VCAP.Get(`user-provided.#(name=="queue-server").credentials.endpoint`).String()
	default:
		log.Fatal("FAIL IN CLOUD ENV")
	}
	return endpoint
}
