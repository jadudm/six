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
http POST http://localhost:8484/search/example.com terms="something or other"
```


(base) jadudm@jade:six$ docker run --rm -v $PWD:/tmp aldanial/cloc --exclude-dir=static .
github.com/AlDanial/cloc v 2.02  T=0.01 s (5261.0 files/s, 284802.8 lines/s)

```
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