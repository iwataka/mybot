# Mybot

## Usage

First, you should create `config.yml` in this project root or anywhere else if
you want to specify the location when executing the command.

Example configuration is like below:

```yaml
github:
        projects:
                - user: vim
                  repo: vim
                - user: golang
                  repo: go
        duration: 30m

retweet:
        accounts:
                - name: golang
                  patterns
                          - is released!
                          - "#golang"
                  opts:
                          retweeted: false
                - name: vimtips
                  opts:
                          hasMedia: false
                          hasUrl: true
                          retweeted: false
        notification:
                place:
                        allowSelf: true
                        # users:
                        #         - foo
                        #         - bar
        duration: 30m

interaction:
        duration: 30s
        # users:
        #         - foo
        #         - bar

log:
        allowSelf: true
        # users:
        #         - foo
        #         - bar

authentication:
        consumerKey: TWITTER_CONSUMER_KEY
        consumerSecret: TWITTER_CONSUMER_SECRET
        accessToken: TWITTER_ACCESS_TOKEN
        accessTokenSecret: TWITTER_ACCESS_TOKEN_SECRET

option:
        name: Iwataka
```

Then run `mybot serve` and that's all.
