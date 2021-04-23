# Mybot

[![Build Status](https://travis-ci.org/iwataka/mybot.svg?branch=master)](https://travis-ci.org/iwataka/mybot)
[![Coverage Status](https://img.shields.io/coveralls/github/iwataka/mybot/master.svg)](https://coveralls.io/github/iwataka/mybot?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwataka/mybot)](https://goreportcard.com/report/github.com/iwataka/mybot)
[![GoDoc](https://godoc.org/github.com/iwataka/mybot?status.svg)](https://godoc.org/github.com/iwataka/mybot)

## Introduction

Mybot is a bot server to collect and broadcast information via social network.

## Prerequisites

Make sure you've already created the following things:

- Twitter account
- Slack account

and you've installed Go 1.13.x or later.

## Running from source

Run the below commands:

```sh
$ go get -d github.com/iwataka/mybot
$ cd $GOPATH/src/github.com/iwataka/mybot
$ go build
$ ./mybot s(erve)
```

## Running by using Docker

1. simplest way

    ```sh
    $ docker run -d -p 8080:8080 iwataka/mybot
    ```

1. docker with mounting volumes

    ```sh
    $ docker run -d -p 8080:8080 \
        -v ~/.cache/mybot:/root/.cache/mybot \
        -v ~/.config/mybot:/root/.config/mybot \
        -v ~/.config/gcloud:/root/.config/gcloud \
        iwataka/mybot
    ```

1. docker-compose

    ```sh
    $ curl -fLO https://raw.githubusercontent.com/iwataka/mybot/master/docker-compose.yml
    $ docker-compose up -d
    ```

## To use Google Cloud APIs

Mybot uses the following Google Cloud API:

- [Vision API](https://cloud.google.com/vision/docs) to analyze images attached to tweets and messages.
- [Natural Language API](https://cloud.google.com/natural-language) to analyze texts in tweets and messages.

To get authorized, run the following commands:

```sh
$ gcloud auth application-default login --scopes=https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/cloud-vision,https://www.googleapis.com/auth/cloud-language
```
