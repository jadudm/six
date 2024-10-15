# six

Exploration number six.

# to run/build quickly/incrementally

This does a `run *.go` on the directories, instead of full container builds.

Build the base image

```
docker build -t six/dev -f Dockerfile.builder .
```

```
docker compose up
```


# run/build statically

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
http POST http://localhost:8484/search/example.com search-terms="something or other"
```


(base) jadudm@jade:six$ docker run --rm -v $PWD:/tmp aldanial/cloc .
github.com/AlDanial/cloc v 2.02  T=0.01 s (5261.0 files/s, 284802.8 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              32            317            177           1620
JSON                             5              0              0            563
Markdown                         4             83              0            238
YAML                             3             10             24            136
Text                             1              0              0            127
make                            10             21              2             74
Dockerfile                       4             18             16             30
SQL                              3              7              8             20
Ruby                             1              2              0             17
Bourne Shell                     4             13             91             13
-------------------------------------------------------------------------------
SUM:                            67            471            318           2838
-------------------------------------------------------------------------------