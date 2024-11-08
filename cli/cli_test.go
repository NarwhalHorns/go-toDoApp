package cli

import (
	"os"
	"testing"
)

func TestCommandSwitch(t *testing.T) {
	t.Run("display help", func(t *testing.T) {
		var loop = true
		commandSwitch("notacommand", make(chan os.Signal), &loop)
		if got != got
	})
}
