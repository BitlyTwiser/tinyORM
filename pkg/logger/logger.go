package logger

import (
	"os"

	"github.com/BitlyTwiser/slogger"
)

var Log *slogger.Logger

func init() {
	Log = slogger.NewLogger(os.Stdout)
}
