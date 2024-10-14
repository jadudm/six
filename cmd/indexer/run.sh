#!/bin/bash

# This runs in the container context only.
export VCAP_SERVICES=$(cat /app/vcap.json)

echo Running minio commands
# Now, make sure the indexer buckets are present
/minio/mc alias set minnie http://minio:9000 nutnutnut nutnutnut
/minio/mc mb indexed-content

./indexer.exe