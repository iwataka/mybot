FROM golang:1.9-alpine

RUN apk add --no-cache git
# https://github.com/kubernetes/kubernetes/issues/39583
RUN git config --global http.https://gopkg.in.followRedirects true
RUN go get -u github.com/golang/dep/cmd/dep
COPY . $GOPATH/src/github.com/iwataka/mybot

RUN cd $GOPATH/src/github.com/iwataka/mybot \
            && dep ensure \
            && go build \
            && cp mybot /usr/local/bin

CMD mybot serve -H 0.0.0.0

EXPOSE 8080
