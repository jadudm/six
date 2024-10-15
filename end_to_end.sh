#!/bin/bash

host=$1

http PUT http://localhost:8080/enqueue/CRAWL host=${host} path=/
read -p "Press <enter> when crawling is complete." bah
http PUT http://localhost:8080/enqueue/PACK host=${host} type=pack_full
read -p "Press <enter> when packing is complete." bah
http PUT http://localhost:8080/enqueue/SEARCH host=${host} type=search search-id=searcher
echo "Run some queries!"

echo http POST http://localhost:8484/search/ terms="something or other"
