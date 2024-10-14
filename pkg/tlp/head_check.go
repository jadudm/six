package tlp

import (
	"net/url"

	gtst "com.jadud.search.six/pkg/types"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

func HeadCheck(ch_in <-chan gtst.JSON, ch_out chan<- gtst.JSON) {
	client := resty.New()

	for {
		msg := <-ch_in

		var u url.URL
		u.Scheme = "https"
		u.Host = gjson.GetBytes(msg, "host").String()
		u.Path = gjson.GetBytes(msg, "path").String()

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/HEAD
		resp, err := client.R().Head(u.String())
		if err != nil {
			// log.Println("HEAD ERR", err)
		} else {
			if resp.StatusCode() == 200 {
				j := gtst.BytesToMap(msg)
				j["content-type"] = resp.Header().Get("Content-Type")
				j["content-length"] = resp.Header().Get("Content-Length")
				j["url"] = u.String()
				ch_out <- gtst.MapToBytes(j)
			} else {
				// SKIP
			}
		}
	}
}
