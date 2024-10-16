package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	queueing "com.jadud.search.six/cmd/queue-server/pkg/queueing"
	obj "com.jadud.search.six/pkg/object-storage"
	gtst "com.jadud.search.six/pkg/types"
	"com.jadud.search.six/pkg/vcap"
	"github.com/PuerkitoBio/goquery"
	cache "github.com/go-pkgz/expirable-cache/v3"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
		return s
	} else {
		return s
	}
}

func is_crawlable(host *url.URL, link string) (string, error) {
	base := host
	// FIXME: we should have the scheme in the host?
	base.Scheme = "https"
	base.Path = ""

	// Is the URL at least length 1?
	if len(link) < 1 {
		return "", errors.New("crawler: URL is too short to crawl")
	}

	// Is it a relative URL? Then it is OK.
	if string([]rune(link)[0]) == "/" {
		u, _ := url.Parse(link)
		u = base.ResolveReference(u)
		return u.String(), nil
	}

	lu, err := url.Parse(link)
	if err != nil {
		return "", errors.New("crawler: link does not parse")
	}

	// Does it end in .gov?
	if bytes.HasSuffix([]byte(lu.Host), []byte("gov")) {
		return "", errors.New("crawler: URL does not end in .gov")
	}

	pieces := strings.Split(base.Host, ".")
	if len(pieces) < 2 {
		return "", errors.New("crawler: link host has too few pieces")
	} else {
		tld := pieces[len(pieces)-2] + "." + pieces[len(pieces)-1]

		// Does the link contain our TLD?
		if !strings.Contains(lu.Host, tld) {
			return "", errors.New("crawler: link does not contain the TLD")
		}
	}

	// FIXME: There seem to be whitespace URLs coming through. I don't know why.
	// This could be revisited, as it is expensive.
	// Do we still have garbage?
	if !bytes.HasSuffix([]byte(lu.String()), []byte("https")) {
		return "", errors.New("crawler: link does not start with https")
	}
	// Is it pure whitespace?
	if len(strings.Replace(lu.String(), " ", "", -1)) < 5 {
		return "", errors.New("crawler: link too short")
	}
	return lu.String(), nil
}

// func remove_trailing_slash(u *url.URL) *url.URL {
// 	// Always strip trailing `/`, as it causes confusion.
// 	if bytes.HasSuffix([]byte(u.Path), []byte("/")) {
// 		after := trimSuffix(u.Path, "/")
// 		log.Printf("TRIMMED SLASH %s -> %s\n", u.Path, after)
// 		u.Path = after
// 	} else {
// 		log.Printf("DOES NOT END IN SLASH %s", u.Path)
// 	}
// 	return u
// }

func fetch_links(host *url.URL, url string) []string {
	client := resty.New()
	log.Println("Fetching", url)
	resp, err := client.R().Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.RawBody().Close()
	reader := bytes.NewReader(resp.Body())

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	links := make([]string, 0)

	doc.Find("a[href]").Each(func(ndx int, sel *goquery.Selection) {
		link, exists := sel.Attr("href")

		if exists {
			link_to_crawl, err := is_crawlable(host, link)
			if err != nil {
				log.Println(err)
			}
			links = append(links, link_to_crawl)
			// log.Println("YES", host, link_to_crawl)
		} else {
			// log.Println("NO", link, host)
		}
	})

	// Remove all trailing slashes.
	for ndx, link := range links {
		links[ndx] = trimSuffix(link, "/")
	}
	return links
}

// mv to types
// also split from pkg to internal -- only things to share ext. should be in pkg
func get(msg gtst.JSON, key string) string {
	return gjson.GetBytes(msg, key).String()
}

func get_ttl(vcap_services *vcap.VcapServices) int64 {
	qm := fmt.Sprintf("user-provided.#(instance_name==\"crawler\").credentials.%s", "cache-ttl-minutes")
	qs := fmt.Sprintf("user-provided.#(instance_name==\"crawler\").credentials.%s", "cache-ttl-seconds")
	minutes := vcap_services.VCAP.Get(qm).Int()
	seconds := vcap_services.VCAP.Get(qs).Int()
	return (minutes * 60) + seconds
}

func Crawl(vcap_services *vcap.VcapServices, buckets *obj.Buckets, ch_msg <-chan gtst.JSON) {
	qs := queueing.NewQueueServer("queue-server", vcap_services)

	ttl := get_ttl(vcap_services)
	log.Printf("CRAWLER ttl seconds [%d]\n", ttl)
	c := cache.NewCache[string, int]().WithTTL(time.Second * time.Duration(ttl))

	for {
		msg := <-ch_msg
		m := gtst.BytesToMap(msg)
		content_type := get(msg, "content-type")

		if strings.Contains(content_type, "html") {
			// log.Println("crawler found html")
			u, _ := url.Parse("https://" + get(msg, "host") + get(msg, "path"))

			// Add ourselves to the cache.
			// log.Println("Adding to cache", u.String())
			c.Set(u.String(), 0, 0)

			new_links := fetch_links(u, get(msg, "url"))

			for ndx, link := range new_links {
				if link != "" {
					_, ok := c.Get(link)
					// If it is not in the  cache, crawl it.
					if !ok {
						c.Set(link, ndx, 0)
						log.Printf("CRAWL crawing [%s]\n", link)
						u, _ := url.Parse(link)
						m["host"] = u.Host
						m["path"] = u.Path
						// The crawler both crawls and queues pages for scraping.
						qs.Enqueue("CRAWL", gtst.MapToBytes(m))
						qs.Enqueue("SCRAPE", gtst.MapToBytes(m))
					} else {
						log.Printf("CRAWLER cache hit[%s]\n", link)
					}

				}
			}
			// Check cache
			// Post to queue

		}
	}
}
