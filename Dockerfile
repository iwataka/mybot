FROM golang:1.8-alpine

WORKDIR /mybot

RUN apk add --no-cache git

CMD go get -d ./... && go build ./cmd/mybot && ./mybot serve > .mybot-dump.log 2>&1

EXPOSE 8080
