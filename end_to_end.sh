#!/bin/bash

all () {
    local host="$1"
    http PUT http://localhost:8080/enqueue/CRAWL host=${host} path=/
    sleep 10m
    http PUT http://localhost:8080/enqueue/PACK host=${host} type=pack_full
    sleep 10m
    http PUT http://localhost:8080/enqueue/SEARCH host=${host} type=search search-id=searcher
    sleep 10m
}

all jadud.com
all www.search.gov
all www.fac.gov
all www.cloud.gov

#read -p "Press <enter> when crawling is complete." bah
#read -p "Press <enter> when packing is complete." bah
#echo "Run some queries!"
#echo http POST http://localhost:8484/search/ terms="something or other"
