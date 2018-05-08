# slack-logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:" /> [![Build Status](https://travis-ci.org/codehakase/slack-logrus.svg?branch=master)](https://travis-ci.org/polds/logrus-papertrail-hook)&nbsp;[![godoc reference](https://godoc.org/github.com/codehakase/slack-logrus?status.png)]https://godoc.org/github.com/codehakase/slack-logrus)

:memo::memo: An Intuitive Logrus Hook for Slack

## Installation
Using `go get`
```shell
$ go get -u github.com/codehakase/slack-logrus
```

## Usage
Log Directly to A Slack Channel

```go
package main
import (
  "github.com/codehakase/slack-logrus"
  "github.com/sirupsen/logrus"
)

func main() {

  logrus.SetOutput(os.Stderr)
  logrus.AddHook(&sl.Hook{
    SlackHookURL: "https://hooks.slack.com/services/xxxxxx/xxxxxx/xxxxxxx",
    Username:     "slack-logrus",
    IconEmoji:    ":mega:",
    WithLevels:   logrus.AllLevels,
    Channel:      "#dev-channel",
  })
  logrus.Info("Hey I'm test running my package!!!")
}
```

### Send Extra Loggin Data
To send custom fields to aid understand logging use add an `Options` field
```go
fields := map[string]interface{}{
  "hostname": "hakaselabs",
  "source":   os.Getenv("APISource"),
  "tag":      "test-tag",
}
logrus.AddHook(&slacklogrus.Hook{
  ...
  Options: fields,
})
```

