package tlp

import (
	"log"
	"strings"

	qing "com.jadud.search.six/cmd/queue-server/pkg/queueing"
	gtst "com.jadud.search.six/pkg/types"
)

const (
	HTML_Q = iota
	PDF_Q
	WORD_Q
	queue_size
)

var queues = []string{"HTML", "PDF", "WORD"}

func QueueToIndexer(ch_msg <-chan gtst.JSON) {
	// A good message has come in. Send it to the right home.
	qs := qing.QueueServer{
		Scheme: "http",
		Host:   "localhost",
		Port:   6000,
	}

	for {
		msg := <-ch_msg
		m := gtst.BytesToMap(msg)
		if strings.Contains(m["content-type"].(string), "html") {
			log.Println("found html")
			qs.Enqueue(queues[HTML_Q], msg)
		} else if strings.Contains(m["content-type"].(string), "pdf") {
			log.Println("found pdf")
			qs.Enqueue(queues[PDF_Q], msg)
		}
	}
}
