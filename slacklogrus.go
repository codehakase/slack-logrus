// Package slacklogrus provides a hook for Slack and the Logrus package.
//
// Installation
//   go get -u github.com/codehakase/slack-logrus
//
//
// Usage

//   package main
//   import (
//     "github.com/codehakase/slack-logrus"
//     "github.com/sirupsen/logrus"
//   )

//   func main() {

//     logrus.SetOutput(os.Stderr)
//     logrus.AddHook(&sl.Hook{
//       SlackHookURL: "https://hooks.slack.com/services/xxxxxx/xxxxxx/xxxxxxx",
//       Username:     "slack-logrus",
//       IconEmoji:    ":mega:",
//       WithLevels:   logrus.AllLevels,
//       Channel:      "#dev-channel",
//     })
//     logrus.Info("Hey I'm test running my package!!!")
//   }
//
// // To send custom fields to aid understand logging use add an `Options` field
//   fields := map[string]interface{}{
//     "hostname": "avivacore",
//     "source":   os.Getenv("APISource"),
//     "tag":      "test-tag",
//   }
//   logrus.AddHook(&slacklogrus.Hook{
//     ...
//     Options: fields,
//   })
package slacklogrus

import (
	"fmt"

	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/webhook"
	"github.com/sirupsen/logrus"
)

// Hook is a logrus Hook for dispatching messages to a defined slack channel
type Hook struct {
	WithLevels   []logrus.Level
	SlackHookURL string
	Channel      string
	IconEmoji    string
	IconURL      string
	Username     string
	Options      map[string]interface{}
}

// Implement Hook methods

// Fire creates a new event, and sends to Slack
func (h *Hook) Fire(entry *logrus.Entry) error {

	var wlvl_match bool
	for _, wlvl := range h.WithLevels {
		if wlvl == entry.Level {
			wlvl_match = true
			break
		}
	}

	// entry is not of desired log level
	if !wlvl_match {
		return nil
	}

	var color string
	switch entry.Level {
	case logrus.InfoLevel:
		color = "good"
	case logrus.DebugLevel:
		color = "#9B30FF"
	case logrus.FatalLevel, logrus.PanicLevel, logrus.ErrorLevel:
		color = "danger"
	default:
		color = "warning"
	}

	// setup slack connection
	c := webhook.New(h.SlackHookURL)
	text := &chat.Message{
		Username:  h.Username,
		Channel:   h.Channel,
		IconEmoji: h.IconEmoji,
		IconURL:   h.IconURL,
	}
	// add log data as slack attachment
	atx := text.NewAttachment()
	// fetch all entries and add as attachment
	allEntries := h.newEntry(entry)

	if len(allEntries.Data) > 0 {
		atx.Text = "Log fields"
		for k, v := range allEntries.Data {
			field := &chat.Field{}
			field.Title = k
			field.Value = fmt.Sprint(v)
			atx.AddField(field)
		}
		atx.PreText = allEntries.Message
	} else {
		atx.Text = allEntries.Message
	}
	atx.Color = color
	// Execute concurrently
	go text.Send(c)
	return nil
}

// newEntry adds a new entry to the Logger
func (h *Hook) newEntry(entry *logrus.Entry) *logrus.Entry {
	data := map[string]interface{}{}

	for k, v := range h.Options {
		data[k] = v
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	return &logrus.Entry{
		Logger:  entry.Logger,
		Time:    entry.Time,
		Level:   entry.Level,
		Data:    data,
		Message: entry.Message,
	}
}

// SetLevels specifies neccessary levels for a particular hook
func (h *Hook) SetLevels(level []logrus.Level) {
	h.WithLevels = level
}

// Levels sent to slack
func (h *Hook) Levels() []logrus.Level {
	if h.WithLevels == nil {
		return logrus.AllLevels
	}
	return h.WithLevels
}
