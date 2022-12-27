package slogger_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/BitlyTwiser/tinyORM/pkg/slogger"
)

func TestSloggerOutput(t *testing.T) {
	tests := []struct {
		name       string
		have       string
		want       string
		testSwitch int
	}{
		{
			name:       "Arguments Test",
			have:       "",
			want:       "",
			testSwitch: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.testSwitch {

			}
			if test.have != test.want {
				t.Errorf("Have: %v Want: %v", test.have, test.want)
			}
		})
	}
}

func TestSloggerToFile(t *testing.T) {
	f, err := os.OpenFile("/tmp/log.json", os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create("/tmp/log.json")
			if err != nil {
				t.Fatalf("could not write to log file. %v", err.Error())
			}

		}
	}

	defer f.Close()

	// fileLogger := slogger.NewLogger(f)
	// fileLogger.LogEvent("warn", "Something one", "one", "two", "Another", false, "four", true, "bob")
	// fileLogger.LogEvent("warn", "Something two", "one", "two", "Another", false, "four", true)
	// fileLogger.LogEvent("warn", "Something three", "hic", "leaf", "Another")

	stdoutLogger := slogger.NewLogger(os.Stdout)
	stdoutLogger.LogEvent("info", "Something four", "key", "value", "AnotherKey", false, "four", 123123)
	stdoutLogger.LogEvent("info", "Something five", "key")
	stdoutLogger.LogEvent("debug", "Something six", map[string]any{"one": "two", "Another": false, "four": true})
	stdoutLogger.LogEvent("info", "Something debug", map[string]any{"one": "two", "Another": false, "four": true})
	stdoutLogger.LogEvent("info", "Hi mom!")
	stdoutLogger.LogEvent("debug", "Hi mom!")
	stdoutLogger.LogError("error", fmt.Errorf("something died"))
	stdoutLogger.LogError("error", fmt.Errorf("something died"), 1, 2, 3, "masdasd")
}
