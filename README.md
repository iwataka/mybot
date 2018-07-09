# Mybot

[![Build Status](https://travis-ci.org/iwataka/mybot.svg?branch=master)](https://travis-ci.org/iwataka/mybot)
[![Coverage Status](https://img.shields.io/coveralls/github/iwataka/mybot/master.svg)](https://coveralls.io/github/iwataka/mybot?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwataka/mybot)](https://goreportcard.com/report/github.com/iwataka/mybot)

## Introduction

Mybot is a bot server to collect and broadcast information via social network.

## Getting Started

Firstly, make sure you've already created the following things:

- Twitter account
- Slack account

To get started, just run the below command:

```
go get -u github.com/iwataka/mybot
mybot s(erve)
```

or by using Docker:

```
docker run -d --name mybot -p 8080:8080 iwataka/mybot
```

or by using docker-compose:

```
go get -d github.com/iwataka/mybot
cd $GOPATH/src/github.com/iwataka/mybot
make deploy_app
# To clean app, run `make clean_app`
```
