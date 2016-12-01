#!/bin/bash

set -e

dir=$(cd $(dirname $0) && pwd)

docker build -t mybot "$dir"

if [[ -z $MYBOT_CONTAINER_NAME ]]; then
    container_name="mybot"
else
    container_name=$MYBOT_CONTAINER_NAME
fi

if [[ -z $MYBOT_NETWORK ]]; then
    network="host"
else
    network=$MYBOT_NETWORK
fi

docker run -d \
    -v "$dir":/mybot \
    -v "$HOME/.config/gcloud":/root/.config/gcloud \
    -v "$HOME/.cache/mybot":/root/.cache/mybot \
    --name "$container_name" \
    --net=$network \
    mybot sh -c "cd /mybot && go get -d ./... && go generate && go build && ./mybot s --log .mybot-debug.log > .mybot-dump.log 2>&1"
