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

	slogger.NewLogger(f).LogEvent("warn", "Something happened", slogger.Arguments("one", "two", "Another", false, "four", true, "bob"))
	slogger.NewLogger(f).LogEvent("warn", "Something happened", slogger.Arguments("one", "two", "Another", false, "four", true))

	slogger.NewLogger(f).LogEvent("warn", "Something happened", slogger.Arguments("hic", "leaf", "Another"))

	slogger.NewLogger(os.Stdout).LogEvent("info", "Something happened", slogger.Arguments("key", "value", "AnotherKey", false, "four", 123123))

	slogger.NewLogger(os.Stdout).LogEvent("debug", "Something happened", map[string]any{"one": "two", "Another": false, "four": true})

	slogger.NewLogger(os.Stdout).LogEvent("info", "Hi mom!")
	slogger.NewLogger(os.Stdout).LogEvent("debug", "Hi mom!")

	slogger.NewLogger(os.Stdout).LogError("error", fmt.Errorf("something died"))

	slogger.NewLogger(os.Stdout).LogError("error", fmt.Errorf("something died"), slogger.Arguments(1, 2, 3, "masdasd"))
}
