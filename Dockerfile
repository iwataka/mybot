FROM golang:1.17-alpine AS golang

RUN apk add --no-cache git
# https://github.com/kubernetes/kubernetes/issues/39583
RUN git config --global http.https://gopkg.in.followRedirects true
WORKDIR /mybot
COPY . .
RUN go build

FROM node:16-alpine

COPY --from=golang /mybot /mybot
WORKDIR /mybot
RUN yarn --cwd ./web install
RUN yarn --cwd ./web build

CMD ./mybot serve -H 0.0.0.0
EXPOSE 8080
