package slacklogrus

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLevels(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	hook := &Hook{
		SlackHookURL: os.Getenv("TRAVIS_TEST_SLACK_HOOK_URL"),
		Username:     "hakaselabs-logrus",
		IconEmoji:    ":mega:",
		WithLevels:   logrus.AllLevels,
		Channel:      "#hakaselabs-dev-chan",
	}
	logrus.AddHook(hook)
	logrus.Debug("logging in Debug Mode")
	logrus.Info("Logging in Info Mode")
	if len(hook.Levels()) < 1 {
		t.Error("Error setting level, level length less than 1")
	}
}

func TestSendSlackLog(t *testing.T) {
	hook := &Hook{
		SlackHookURL: os.Getenv("TRAVIS_TEST_SLACK_HOOK_URL"),
		Username:     "hakaselabs-logrus",
		IconEmoji:    ":mega:",
		WithLevels:   logrus.AllLevels,
		Channel:      "#hakaselabs-dev-chan",
	}
	err := hook.Fire(&logrus.Entry{
		Data: map[string]interface{}{
			"tag": "testing",
		},
		Logger:  &logrus.Logger{},
		Message: "Testing Slacklogrus Log",
	})

	if err != nil {
		t.Errorf("Could not fire slacklogrus: %v", err)
	}
}
