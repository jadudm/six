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

	// FIXME: these should be in VCAP
	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	sleep_duration := 1 * time.Second
	min_sleep := 100 * time.Millisecond
	max_sleep := 30 * time.Second

	log.Printf("queue-server: %s", endpoint)

	retries := 0
	for {
		resp, err := client.R().
			EnableTrace().
			Get(fmt.Sprintf("%s/dequeue/%s", endpoint, queue))
		//log.Println(resp)
		if err == nil {
			rmap := gtst.BytesToMap(resp.Body())
			if rmap["result"] == nil || rmap["result"].(string) == "error" {
				retries += 1
				sleep_duration += CalculateBackoff(retries, initialBackoff, maxBackoff)
				if sleep_duration > time.Duration(max_sleep) {
					sleep_duration = max_sleep
				}

			} else {
				retries = 0
				sleep_duration = min_sleep
				body := resp.Body()
				ch_msg <- body

			}
		}

		log.Printf("queue checker sleeping [%s]", sleep_duration)
		time.Sleep(sleep_duration)
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
