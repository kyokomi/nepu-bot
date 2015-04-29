nepu-bot
=====

[![Circle CI](https://circleci.com/gh/kyokomi/nepu-bot.svg?style=svg)](https://circleci.com/gh/kyokomi/nepu-bot)

slack hubot for Go!（golang）

超次元ゲイムネプテューヌに登場した[いーすん](http://dic.pixiv.net/a/%E3%82%A4%E3%82%B9%E3%83%88%E3%83%AF%E3%83%BC%E3%83%AB%28%E3%83%8D%E3%83%97%E3%83%86%E3%83%A5%E3%83%BC%E3%83%8C%29)BOT 個人用です。

- 語尾に何かしら顔文字をつけます
- 5回に1回、会話に割って入ってきます
- `いーすん おしえて ◯◯`で質問できます


# Setup

```
$ export DOCOMO_APIKEY={docomoDeveloper api key}
$ export SLACK_TOKEN={slack bot token}
```

# Usage

```
$ go build
$ ./nepu-bot
```

![slackbot](https://qiita-image-store.s3.amazonaws.com/0/40887/704f1f94-d8ae-76a7-5417-0d841288f2b0.png "slackbot")

