#!/bin/bash

set -e

dir=$(cd $(dirname $0) && pwd)

docker build -t mybot "$dir"

if [[ -z $MYBOT_CONTAINER_NAME ]]; then
    container_name="mybot"
else
    container_name=$MYBOT_CONTAINER_NAME
fi

docker run -d -v "$dir":/mybot \
    --name "$container_name" \
    -p 8080:8080 \
    mybot sh -c "cd /mybot && go get -d ./... && go generate && go build && ./mybot s --log .mybot-debug.log &>> .mybot-debug.log"
