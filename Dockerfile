# [Usage]
# cd /path/to/mybot
# docker build -t mybot .
# docker run -d -v `pwd`:/mybot -v ~/.config/gcloud:/root/.config/gcloud -v ~/.cache/mybot:/root/.cache/mybot --name mybot -p 8080:8080 mybot
FROM golang:1.8-alpine

WORKDIR /mybot

RUN apk add --no-cache git
RUN go get -u github.com/jteeuwen/go-bindata/...

CMD go get -d ./... && go generate && go build && ./mybot serve --log .mybot-debug.log > .mybot-dump.log 2>&1

EXPOSE 8080
