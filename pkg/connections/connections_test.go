package connections_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/BitlyTwiser/slogger"
)

func TestPackage(t *testing.T) {
	s := slogger.NewLogger(os.Stdout)
	s.LogError("death", fmt.Errorf("error happened"))
}
