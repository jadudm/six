package queueing

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	gtst "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
	"github.com/go-resty/resty/v2"
)

type QueueServer struct {
	Scheme string
	Host   string
	Port   int
	client *resty.Client
}

func get_param_string(vc *vcap.VcapServices, service_name string, field string) string {
	query := fmt.Sprintf("user-provided.#(instance_name==\"%s\").credentials.%s", service_name, field)
	s := vc.VCAP.Get(query).String()
	return s
}

func get_param_int(vc *vcap.VcapServices, service_name string, field string) int {
	query := fmt.Sprintf("user-provided.#(instance_name==\"%s\").credentials.%s", service_name, field)
	s := vc.VCAP.Get(query).Int()
	return int(s)
}

func NewQueueServer(queue_server_name string, vcap_services *vcap.VcapServices) *QueueServer {
	// Looks for a queue server by name, and grabs the info we need.

	return &QueueServer{
		Scheme: get_param_string(vcap_services, queue_server_name, "scheme"),
		Host:   get_param_string(vcap_services, queue_server_name, "host"),
		Port:   get_param_int(vcap_services, queue_server_name, "port"),
	}
}

func buildUrl(qs *QueueServer, path string) url.URL {
	var u url.URL
	u.Scheme = qs.Scheme
	u.Host = qs.Host + ":" + strconv.Itoa(qs.Port)
	u.Path = path
	return u
}

func (qs *QueueServer) Enqueue(queue string, msg gtst.JSON) {
	u := buildUrl(qs, "enqueue/"+queue)

	if qs.client == nil {
		qs.client = resty.New()
	}

	_, err := qs.client.R().
		EnableTrace().
		SetBody(msg).
		Put(u.String())

	if err != nil {
		log.Println("COULD NOT ENQUEUE", u.String())
	}
}
