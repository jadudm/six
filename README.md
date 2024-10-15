# six

*Exploration number six in SQLite-based hackery.*

"What does a distributed crawler/indexer/search engine look like?"

## I just want to run the thing

To run/build quickly/incrementally

This does a `run *.go` on the directories, instead of full container builds.

Build the base image

```
docker build -t six/dev -f Dockerfile.builder .
```

Run the stack. Individual components will be `go run *go`'d.

```
docker compose up
```

# What hath thou wrought?

While it may look a bit early 2000's, I wondered what the core of a distributed/decoupled search infrastructure might look like.

I wanted a crawler that was decoupled (via FIFO queues) from the scraper, which in turn was decoupled from the system that would turn the scraped content into a database.



# Other build notes

To run/build statically (meaning, each container contains a full binary build, and executes those)...

```
docker compose -f compose.exe.yaml up --build
```

# running

To crawl a site

```
http PUT http://localhost:8080/enqueue/CRAWL host=example.com path=/
```

To pack the site that was crawled

```
http PUT http://localhost:8080/enqueue/PACK host=example.com type=pack_full
```

To move it to the searcher

```
http PUT http://localhost:8080/enqueue/SEARCH host=example.com type=search search-id=searcher
```

To run searches

```
http POST http://localhost:8484/search/example.com terms="something or other"
```

or visit

[http://localhost:8484/static/](http://localhost:8484/static/)

after you load your content.


```
six$ docker run --rm -v $PWD:/tmp aldanial/cloc --exclude-dir=static .
github.com/AlDanial/cloc v 2.02  T=0.01 s (5261.0 files/s, 284802.8 lines/s)


-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              37            407            251           2066
JSON                             5              0              0            595
YAML                             4             14             50            305
Markdown                         5            104              0            287
Text                             1              0              0            127
make                            11             27              2             96
Bourne Shell                    11             35            125             75
Dockerfile                       5             24             24             39
SQL                              3             10             11             29
Ruby                             1              6              0             24
-------------------------------------------------------------------------------
SUM:                            83            627            463           3643
-------------------------------------------------------------------------------
```