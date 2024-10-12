package tlp

import (
	"fmt"
	"time"

	GTST "com.jadud.search.six/pkg/types"
	"github.com/go-resty/resty/v2"
)

func checkQueue() {

}

// func check_queue() string {
// 	client := resty.New()
// 	var counter int64 = 0

// 	for EVER {
// 		time.Sleep(time.Duration(QUEUE_CHECK_FREQUENCY) * time.Second)
// 		qj := &GTST.JobResponse{}
// 		resp, err := client.R().
// 			EnableTrace().
// 			SetResult(&qj).
// 			Get("http://localhost:6000/dequeue")
// 		if err != nil {
// 			panic("DEQUEUE FAILED")
// 		}

// 		if qj.Result == "ok" {
// 			return qj.Domain
// 		} else {
// 			log.Println("GET Response:", resp.Status())
// 			log.Println("GET Body:", string(resp.Body()))
// 			counter += 1
// 			log.Printf("Empty job. Counting sheep: %d\n", counter)
// 		}
// 	}

// 	return "check_queue() NEVER GETS HERE"
// }

func CheckQueue(queue string, crontab string, ch_msg chan<- GTST.JSON) {
	client := resty.New()

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

	for {
		resp, err := client.R().
			EnableTrace().
			//FIXME variable/vcap
			Get(fmt.Sprintf("http://localhost:6000/dequeue/%s", queue))
		if err != nil {
			panic("DEQUEUE FAILED")
		}
		ch_msg <- resp.Body()
		time.Sleep(5 * time.Second)
	}
}
