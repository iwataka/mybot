#!/bin/bash

set -e

dir=$(cd $(dirname $0) && pwd)

docker build -t mybot "$dir"

if [[ -z $MYBOT_TWITTER_CONSUMER_KEY ]]; then
    echo -n "MYBOT_TWITTER_CONSUMER_KEY=" && read twitter_consumer_key
else
    twitter_consumer_key=$MYBOT_TWITTER_CONSUMER_KEY
fi
if [[ -z $MYBOT_TWITTER_CONSUMER_SECRET ]]; then
    echo -n "MYBOT_TWITTER_CONSUMER_SECRET=" && read twitter_consumer_secret
else
    twitter_consumer_secret=$MYBOT_TWITTER_CONSUMER_SECRET
fi
if [[ -z $MYBOT_TWITTER_ACCESS_TOKEN ]]; then
    echo -n "MYBOT_TWITTER_ACCESS_TOKEN=" && read twitter_access_token
else
    twitter_access_token=$MYBOT_TWITTER_ACCESS_TOKEN
fi
if [[ -z $MYBOT_TWITTER_ACCESS_TOKEN_SECRET ]]; then
    echo -n "MYBOT_TWITTER_ACCESS_TOKEN_SECRET=" && read twitter_access_token_secret
else
    twitter_access_token_secret=$MYBOT_TWITTER_ACCESS_TOKEN_SECRET
fi

if [[ -z $MYBOT_SLACK_TOKEN ]]; then
    echo -n "MYBOT_SLACK_TOKEN=" && read slack_token
else
    slack_token=$MYBOT_SLACK_TOKEN
fi

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
    mybot sh -c "cd /mybot && go get -d ./... && go build && ./mybot s --log-file /mybot/mybot-debug.log"
