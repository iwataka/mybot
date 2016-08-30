# Mybot

Standalone Twitter bot written in Go

## Why `Mybot`

Twitter is one of the most famous `SNS`s and the largest information sources in the world, but we still can't receive its benefit enough because collecting and broadcasting information relies on your hands and you don't have enough time to do it.
`Mybot` does this kind of things automatically, only what you have to do is writing the configuration file for your own `Mybot`.
`Mybot` does things exactly as you ordered.

## How to collect and broadcast information

`Mybot` defines this task by these 3 components:

+ Source: from which this retrieves information

    + Timeline
    + Search

+ Filter: by which this filters out information

    + Text Patterns (Regular Expression supported)
    + Check whether having media or not
    + Check whether having URL or not
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

## Other useful features

+ Shows some useful pieces of information via built-in Web UI
+ Sends error messages via Direct Message
+ Notifies the place of tweets via Direct Message if available

## Usage

First, you should create `config.toml` in this project root or anywhere else if
you want to specify the location when executing the command.

See [config.template.toml](config.template.toml) for more details.

Then run `mybot serve` and that's all.
