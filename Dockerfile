FROM golang:1.8-alpine

WORKDIR /mybot

RUN apk add --no-cache git
RUN go get -u github.com/jteeuwen/go-bindata/...

CMD go get -d ./... && go generate ./src && go build ./cmd/mybot && ./mybot serve > .mybot-dump.log 2>&1

EXPOSE 8080
