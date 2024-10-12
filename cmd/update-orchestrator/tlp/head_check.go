package tlp

import (
	"log"
	"net/url"

	GTST "com.jadud.search.six/pkg/types"
	"github.com/go-resty/resty/v2"
)

func HeadCheck(ch_in <-chan GTST.JSON, ch_out chan<- GTST.JSON) {
	client := resty.New()

	for {
		msg := <-ch_in
		m := GTST.ToMap(msg)

		if v, ok := m["result"]; ok && v == "ok" {
			var u url.URL
			u.Scheme = "https"
			u.Host = m["domain"].(string)
			u.Path = m["page"].(string)
			log.Printf("HEAD %s", u.String())

			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/HEAD
			resp, _ := client.R().Head(u.String())
			if resp.StatusCode() == 200 {
				log.Println("head check 200")
				ch_out <- msg
			} else {
				log.Printf("head check fail %d", resp.StatusCode())
			}
		}

	}
}
