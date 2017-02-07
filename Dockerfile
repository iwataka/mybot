FROM golang:1.8-alpine

WORKDIR /mybot

RUN apk add --no-cache git
# https://github.com/kubernetes/kubernetes/issues/39583
RUN git config --global http.https://gopkg.in.followRedirects true
RUN go get github.com/iwataka/mybot

CMD mybot serve -H 0.0.0.0

EXPOSE 3256
