FROM golang:1.14-alpine

RUN apk add --no-cache git
# https://github.com/kubernetes/kubernetes/issues/39583
RUN git config --global http.https://gopkg.in.followRedirects true

WORKDIR $GOPATH/src/github.com/iwataka/mybot

COPY go.* ./
RUN go mod download

COPY . .
RUN go build && cp mybot /usr/local/bin

CMD mybot serve -H 0.0.0.0
EXPOSE 8080
