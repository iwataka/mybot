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
    -e MYBOT_TWITTER_CONSUMER_KEY="$twitter_consumer_key" \
    -e MYBOT_TWITTER_CONSUMER_SECRET="$twitter_consumer_secret" \
    -e MYBOT_TWITTER_ACCESS_TOKEN="$twitter_access_token" \
    -e MYBOT_TWITTER_ACCESS_TOKEN_SECRET="$twitter_access_token_secret" \
    -e MYBOT_SLACK_TOKEN="$slack_token" \
    --name "$container_name" \
    -p 8080:8080 \
    mybot sh -c "cd /mybot && go get -d ./... && go build && ./mybot s"
