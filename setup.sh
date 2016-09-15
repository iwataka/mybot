#!/bin/bash

set -e

dir=$(cd $(dirname $0) && pwd)

docker build -t mybot "$dir"

if [[ -z $MYBOT_CONTAINER_NAME ]]; then
    container_name="mybot"
else
    container_name=$MYBOT_CONTAINER_NAME
fi

if [[ -z $MYBOT_PORT ]]; then
    port="8080"
else
    port=$MYBOT_PORT
fi

docker run -d -v "$dir":/mybot \
    --name "$container_name" \
    -p "$port":8080 \
    mybot sh -c "cd /mybot && go get -d ./... && go generate && go build && ./mybot s --log .mybot-debug.log > .mybot-dump.log 2>&1"
