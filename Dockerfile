FROM golang:1.8-alpine

WORKDIR /mybot

RUN apk add --no-cache git
RUN go get github.com/iwataka/mybot

CMD mybot serve -H 0.0.0.0

EXPOSE 3256
