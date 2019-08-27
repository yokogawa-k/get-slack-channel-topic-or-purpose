Slack の Channel Topic や Purpose を誤って変更した際に一つ前の Topic などを取得するツール

## 使い方

```console
Usage:
  get-slacklog [OPTIONS]

Application Options:
  -d, --debug          Show debug infomation
      --slack_token=   Set Slack API Token [$SLACK_TOKEN]
      --slack_channel= Set target Slack channel name [$SLACK_CHANNEL]
      --subtype=       Set target Slack message subtype(channel_topic, channel_purpose, channel_join etc...) (default: channel_topic) [$SLACK_MSG_SUBTYPE]

Help Options:
  -h, --help           Show this help message
```

利用例

``` console
$ ./get-slacklog --slck_channel=app-test --subtype=channel_purpose
INFO[0000] Loglevel is info
2019-08-26 01:53:06 +0900 JST: <@XXXXXXXX> set the channel purpose: Slack App のテスト用、目的追加
 url: https://xxxxxx.slack.com/archives/XXXXXXXXX/p1566751986000800
2019-08-26 01:52:04 +0900 JST: <@XXXXXXXXX> set the channel purpose: Slack App のテスト用
 url: https://xxxxxx.slack.com/archives/XXXXXXXXX/p1566751924000300
```

## Build

Makefile を利用する場合、Docker 19.03 以上が必要

```console
$ make build
```

## Slack API Token に必要な Scope

- team:read
  - 引用などで利用するための URL で利用
  - 生成する URL に利用するワークスペース名を取得するため
- channels:read
  - Channel ID を取得するため
- channels:history
  - ログを取得するため
