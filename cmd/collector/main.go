package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	cp "github.com/otiai10/copy"
)

const EVER bool = true

const QUEUE_CHECK_FREQUENCY int = 3

func load_dotenv() {
	// https://github.com/joho/godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func index_pdf() {
	// domain, path, page_number, text
}

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

type PdfIndex struct {
	r          *colly.Response `json:"-"`
	Domain     string          `json:"domain"`
	Path       string          `json:"path,omitempty"`
	PageNumber int             `json:"page_number,omitempty"`
	Text       string          `json:"text,omitempty"`
}

func (PdfIndex) TableName() string {
	return "pdf_index"
}

type HtmlIndex struct {
	r             *colly.Response `json:"-"`
	bytes_crawled string          `json:"-"\`
	Host          string          `json:"domain"`
	Path          string          `json:"path,omitempty"`
	Title         string          `json:"title"`
	Text          string          `json:"text,omitempty"`
}

func (qj *HtmlIndex) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func load_html_db(hi *HtmlIndex) *sql.DB {
	pool := make(map[string]*sql.DB, 0)

	sqldb, ok := pool[hi.Host]

	if ok {
		return sqldb
	} else {
		db_name := strings.ReplaceAll(hi.Host, ".", "_")
		template_path := "db"
		dyn_db_path := filepath.Join("dynamic", db_name)
		// Filename is relative to the dynpath...
		db_filename := filepath.Join(db_name + ".sqlite")
		dbmate_path := "sqlite:" + db_filename
		gorm_path := filepath.Join(dyn_db_path, db_filename)

		if _, err := os.Stat(dyn_db_path); errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(filepath.Join(dyn_db_path, "db"), 0755)
			err := cp.Copy(template_path, filepath.Join(dyn_db_path, "db"))
			if err != nil {
				log.Print(err)
				panic("COULD NOT COPY DB DIRECTORY")
			}
			// "CREATE VIRTUAL TABLE html_index USING fts5(host, path, title, text);"
			log.Println("DBMATE PATH", dbmate_path)
			cmd := exec.Command("dbmate", "--url", dbmate_path, "up")
			cmd.Dir = dyn_db_path
			out, err := cmd.Output()
			if err != nil {
				log.Println(err)
				log.Println(string(out))
				panic("DBMATE COMMAND FAILED")
			}
		}

		log.Println("OPENING", gorm_path)

		db, err := sql.Open("sqlite3", gorm_path)
		if err != nil {
			log.Println(err)
			panic("COULD NOT OPEN DATABASE")
		}

		pool[hi.Host] = db
		return pool[hi.Host]
	}

}

func index_html(ch_content <-chan *HtmlIndex) {
	for EVER {
		// The crawler sends us things to store.
		// Block, and then write. This makes sure
		// we only ever have one thing going to the DB at a time.
		// Many crawlers, one writer.
		hi := <-ch_content
		db := load_html_db(hi)
		// db.Create(&hi)
		//     USING fts5(host, path, title, text)
		_, err := db.Exec("INSERT INTO html_index(host, path, title, text) VALUES (?, ?, ?, ?)",
			hi.Host, hi.Path, hi.Title, hi.Text)
		if err != nil {
			log.Println(err)
			panic("COULD NOT INSERT INTO HTML INDEX")
		}
	}
}

// crawl -> parse -> store
func parse_html(ch_in <-chan *HtmlIndex, ch_out chan<- *HtmlIndex) {
	for EVER {
		hi := <-ch_in
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(hi.r.Body))
		// hi.Text = doc.Text()

		title := doc.Find("h1").Text()
		if len(title) > 0 {
			hi.Title = title
		}

		// Try grabbing only the paragraphs?
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			para_content := s.Text()
			hi.Text += " " + para_content
		})
		ch_out <- hi
	}
}

func crawl(host string, to_parse chan<- *HtmlIndex) {

	var allowed_domains = [...]string{host}
	var total_bytes int64 = 0

	// Instantiate default collector
	c := colly.NewCollector()
	c.AllowedDomains = allowed_domains[:]

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// fmt.Printf("LINK %s: %q -> %s\n", e.Request.URL.Host, e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		// fmt.Printf("Finished scraping %s\n", r.Request.URL)
		if strings.Contains(r.Headers.Get("Content-Type"), "text/html") {
			total_bytes += int64(len(r.Body))
			the_url := r.Request.URL
			hi := HtmlIndex{}
			hi.r = r
			hi.Host = the_url.Host
			hi.Path = the_url.Path
			hi.bytes_crawled = ByteCountSI(total_bytes)
			to_parse <- &hi
		}

	})

	c.Visit(fmt.Sprintf("https://%s", host))
}

type QueueJob struct {
	Action  string            `json:"action"`
	UUID    string            `json:"uuid"`
	Payload map[string]string `json:"payload"`
}

func check_queue() string {
	client := resty.New()
	var counter int64 = 0

	for EVER {
		time.Sleep(time.Duration(QUEUE_CHECK_FREQUENCY) * time.Second)
		qj := &QueueJob{}
		resp, err := client.R().
			EnableTrace().
			SetResult(&qj).
			Get("http://localhost:6000/dequeue")
		if err != nil {
			panic("RESTY COULD NOT GET A JOB")
		}

		// If we got a UUID, we must be ready to do some work.
		// Break the loop and go!
		if qj.UUID != "" {
			log.Println("QJ", qj)
			// Return the host from the payload
			// More checking would be nice.
			return qj.Payload["crawl"]
		}

		log.Println("GET Response:", resp.Status())
		log.Println("GET Body:", string(resp.Body()))
		counter += 1
		log.Printf("Empty job. Counting sheep: %d\n", counter)
	}

	return "check_queue() NEVER GETS HERE"
}

func main() {
	load_dotenv()
	ch_c2p := make(chan *HtmlIndex)
	ch_p2i := make(chan *HtmlIndex)

	// var wg sync.WaitGroup
	for EVER {
		host := check_queue()
		go func() {
			log.Println("CRAWLING", host)
			go crawl(host, ch_c2p)
			go parse_html(ch_c2p, ch_p2i)
			go index_html(ch_p2i)
		}()
	}
}

/* Cutting room floor */
// For parsing
// Maybe we want to pull the meta out, and use that for weighting.
// var desc string
// doc.Find("meta").Each(func(i int, s *goquery.Selection) {
// 	if val, _ := s.Attr("name"); strings.Contains(val, "description") {
// 		desc, _ = s.Attr("content")
// 	}
// 	if len(desc) == 0 {
// 		if val, _ := s.Attr("property"); strings.Contains(val, "description") {
// 			desc, _ = s.Attr("content")
// 		}
// 	}
// })
// if len(desc) > 0 {
// 	content.Desc = desc
// }

// Scraping
// if strings.Contains(r.Headers.Get("Content-Type"), "application/pdf") {
// 	// https://stackoverflow.com/questions/29746123/convert-byte-slice-to-io-reader
// 	url := r.Request.URL.RequestURI()
// 	fmt.Printf("Processing PDF: %s\n", url)
// 	total_bytes += int64(len(r.Body))
// 	process_pdf_bytes(db, url, r.Body)
// }
