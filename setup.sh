#!/bin/bash

set -e

dir=$(cd $(dirname $0) && pwd)

docker build -t mybot "$dir"

if [[ -z $MYBOT_TWITTER_CONSUMER_KEY ]]; then
    echo -n "MYBOT_TWITTER_CONSUMER_KEY=" && read consumer_key
else
    consumer_key=$MYBOT_TWITTER_CONSUMER_KEY
fi
if [[ -z $MYBOT_TWITTER_CONSUMER_SECRET ]]; then
    echo -n "MYBOT_TWITTER_CONSUMER_SECRET=" && read consumer_secret
else
    consumer_secret=$MYBOT_TWITTER_CONSUMER_SECRET
fi
if [[ -z $MYBOT_TWITTER_ACCESS_TOKEN ]]; then
    echo -n "MYBOT_TWITTER_ACCESS_TOKEN=" && read access_token
else
    access_token=$MYBOT_TWITTER_ACCESS_TOKEN
fi
if [[ -z $MYBOT_TWITTER_ACCESS_TOKEN_SECRET ]]; then
    echo -n "MYBOT_TWITTER_ACCESS_TOKEN_SECRET=" && read access_token_secret
else
    access_token_secret=$MYBOT_TWITTER_ACCESS_TOKEN_SECRET
fi

docker run -d -v "$dir":/mybot \
    -e MYBOT_TWITTER_CONSUMER_KEY="$consumer_key" \
    -e MYBOT_TWITTER_CONSUMER_SECRET="$consumer_secret" \
    -e MYBOT_TWITTER_ACCESS_TOKEN="$access_token" \
    -e MYBOT_TWITTER_ACCESS_TOKEN_SECRET="$access_token_secret" \
    mybot sh -c "cd /mybot && go get -d ./... && go build && ./mybot s"
