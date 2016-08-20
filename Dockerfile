FROM golang:1.7-alpine

RUN apk add --no-cache git
RUN go get -u github.com/jteeuwen/go-bindata/...
EXPOSE 8080
