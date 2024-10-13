package queueing

import (
	"log"
	"net/url"
	"strconv"

	gtst "com.jadud.search.six/pkg/types"
	"github.com/go-resty/resty/v2"
)

type QueueServer struct {
	Scheme string
	Host   string
	Port   int
	client *resty.Client
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
