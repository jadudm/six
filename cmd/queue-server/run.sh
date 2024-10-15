#!/bin/bash 

# This runs in the container context only.
export VCAP_SERVICES=$(cat /app/vcap.json)
echo Running the queue-server
./queue-server.exe
