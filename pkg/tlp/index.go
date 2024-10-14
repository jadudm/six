package tlp

import (
	"bytes"
	"log"
	"strings"

	qing "com.jadud.search.six/cmd/queue-server/pkg/queueing"
	con "com.jadud.search.six/pkg/content"
	obj "com.jadud.search.six/pkg/object-storage"

	. "com.jadud.search.six/pkg/types"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

func fetch_html_content(url string) string {
	client := resty.New()
	log.Println("Fetching", url)
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
		txt := sel.Text()
		if len(txt) > 0 {
			repl := strings.ToLower(txt)
			repl = con.RemoveStopwords(repl)
			repl += " "
			if len(repl) > 2 {
				content += repl
			}
		}
	})

	return content

}

func store_to_s3(b *obj.Bucket, m map[string]interface{}) {
	log.Println("Storing to S3", m["host"])
	log.Println(m["content"])
	path := []string{"indexed", m["host"].(string), uuid.NewString() + ".json"}
	b.PutObject(path, MapToBytes(m))
}

func get(msg JSON, key string) string {
	return gjson.GetBytes(msg, key).String()
}

func Index(buckets *obj.Buckets, ch_msg <-chan JSON) {
	// A good message has come in. Send it to the right home.
	qs := qing.QueueServer{
		Scheme: "http",
		Host:   "localhost",
		Port:   6000,
	}

	for {
		msg := <-ch_msg
		m := BytesToMap(msg)
		content_type := get(msg, "content-type")
		if strings.Contains(content_type, "html") {
			log.Println("found html")
			qs.Enqueue("TIMER", StructToJson(ResetTimer{
				Domain:   get(msg, "domain"),
				Callback: []byte(""),
			}))
			log.Println(m)
			content := fetch_html_content(get(msg, "url"))
			m["content"] = content
			store_to_s3(&buckets.Ephemeral, m)

		} else if strings.Contains(content_type, "pdf") {
			log.Println("found pdf")
		}
	}
}
