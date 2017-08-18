# Mybot

Super cool Twitter bot written in Go

## Why `Mybot`

Twitter is one of the most famous `SNS`s and the largest information sources in
the world, but we still can't receive its benefit enough because collecting and
broadcasting information relies on your hands and you don't have enough time to
do it.
`Mybot` does this kind of things automatically, only what you have to do is
configuring it as you want.
`Mybot` does all things exactly as you ordered.

## How to collect and broadcast information

`Mybot` defines this task by these 3 components:

+ Source: from which this retrieves information

    + Timeline
    + Favorite
    + Search

+ Filter: by which this filters out information

    + Text Patterns (Regular Expression supported)
    + Check whether it has media or not
    + Check whether it has URL or not
    + Retweet Count Threshold
    + Favorite Count Threshold
    + Google Vision API (You need your own credential to use it)

+ Action

    + Retweet
    + Favorite
    + Follow
    + Add to the specified collection (See [Collection API](https://dev.twitter.com/rest/collections) for more details)

## Interact with you via Direct Message

These commands are now available:

+ `collections` (`cols`): returns your own collection list
+ `configuration` (`config`, `conf`): returns your `Mybot` configuration

## Built-in Web Interface

`Mybot` has built-in Web interface and you can monitor and configure this App.

+ Error logs
+ Twitter collection list
+ Google Vision API's result.

## Other useful features

More features availble:

+ Sends error messages via Direct Message
+ Notifies the place of tweets via Direct Message if available

## Usage

Run `mybot serve` and access `localhost:8080`.

### Use with Docker

Run the following commands and that's all.

```
docker run -d -v ~/.config/gcloud:/root/.config/gcloud -v ~/.config/mybot:/root/.config/mybot -v ~/.cache/mybot:/root/.cache/mybot --name mybot -p 8080:8080 iwataka/mybot
```

### MySQL on Docker cotainer

Run the following commands and make configuration for DB.

```
cd /path/to/mybot/scripts
docker run -v `pwd`:/docker-entrypoint-initdb.d -d --name mysql -e MYSQL_ROOT_PASSWORD=mysql mysql
```

And edit config.toml as you want.

## Features to be implemented

+ Plug-in architecture
+ Apply this to other SNSs (like Facebook, Google+, Tumblr,...)
