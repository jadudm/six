#!/bin/bash

export SERVICE=queue-server
export VCAP_SERVICES=$(cat /home/vcap/app/vcap.json)


pushd /home/vcap/app/cmd/${SERVICE}
    echo Building
    make container_build
    echo Running the $SERVICE
    make run
popd