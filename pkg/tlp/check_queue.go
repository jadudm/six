package tlp

import (
	"fmt"
	"log"
	"time"

	gtst "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
	"github.com/go-resty/resty/v2"
)

func CheckQueue(vcap_services *vcap.VcapServices, queue string, crontab string, ch_msg chan<- gtst.JSON) {
	client := resty.New()
	endpoint := GetQueueServerEndpoint(vcap_services)

	log.Printf("queue-server: %s", endpoint)

	for {
		resp, err := client.R().
			EnableTrace().
			Get(fmt.Sprintf("%s/dequeue/%s", endpoint, queue))
		if err == nil {
			body := resp.Body()
			ch_msg <- body
		}
		time.Sleep(3 * time.Second)
	}
}

// c := cron.New()
// c.AddFunc(crontab, func() {
// 	client := resty.New()
// 	resp, err := client.R().
// 		EnableTrace().
// 		//FIXME variable/vcap
// 		Get("http://localhost:6000/dequeue/head")
// 	if err != nil {
// 		panic("DEQUEUE FAILED")
// 	}
// 	ch_msg <- resp.Body()
// })
// c.Start()
