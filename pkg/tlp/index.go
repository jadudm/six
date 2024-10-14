package tlp

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	queueing "com.jadud.search.six/cmd/queue-server/pkg/queueing"
	con "com.jadud.search.six/pkg/content"
	obj "com.jadud.search.six/pkg/object-storage"
	gtst "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

func index_sel(sel *goquery.Selection) string {
	content := ""
	txt := sel.Text()
	if len(txt) > 0 {
		repl := strings.ToLower(txt)
		repl = con.RemoveStopwords(repl)
		repl += " "
		if len(repl) > 2 {
			content += repl
		}
	}
	return content
}

// move to content pkg
func fetch_html_content(url string) string {
	client := resty.New()
	log.Println("INDEXER Fetching", url)
	resp, err := client.R().Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.RawBody().Close()
	reader := bytes.NewReader(resp.Body())
	content := ""

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("p").Each(func(ndx int, sel *goquery.Selection) {
		content += index_sel(sel)
	})
	doc.Find("li").Each(func(ndx int, sel *goquery.Selection) {
		content += index_sel(sel)
	})
	doc.Find("td").Each(func(ndx int, sel *goquery.Selection) {
		content += index_sel(sel)
	})

	return content
}

// mv to object storage
func store_to_s3(b *obj.Bucket, m map[string]interface{}) {
	log.Println("INDEXER Storing to S3", m["host"])
	log.Println(m["content"])
	// Base the filename on the current timestamp.
	// Not UUID : uuid.NewString() + ".json"
	path := []string{"indexed", m["host"].(string), m["indexed-on"].(string) + "-" + uuid.NewString() + ".json"}
	b.PutObject(path, gtst.MapToBytes(m))
}

// mv to types
// also split from pkg to internal -- only things to share ext. should be in pkg
func get(msg gtst.JSON, key string) string {
	return gjson.GetBytes(msg, key).String()
}

func Index(vcap_services *vcap.VcapServices, buckets *obj.Buckets, ch_msg <-chan gtst.JSON) {
	qs := queueing.NewQueueServer("queue-server", vcap_services)
	for {
		msg := <-ch_msg
		m := gtst.BytesToMap(msg)
		content_type := get(msg, "content-type")
		if strings.Contains(content_type, "html") {
			cb := make(map[string]string, 0)
			cb["domain"] = get(msg, "domain")
			cb["type"] = "pack_full"
			qs.Enqueue("TIMER", gtst.MapToBytes(m))

			content := fetch_html_content(get(msg, "url"))

			m["content"] = content
			// FIXME We truncate to the hour.
			// https://gigi.nullneuron.net/gigilabs/golang-timestamps-a-cross-platform-nightmare/
			m["indexed-on"] = fmt.Sprint(time.Now().Truncate(time.Hour).Unix())
			store_to_s3(&buckets.Ephemeral, m)

		} else if strings.Contains(content_type, "pdf") {
			log.Println("found pdf")
		}
	}
}
