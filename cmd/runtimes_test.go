package cmd

import (
	"bytes"
	"testing"
)

func Test_RuntimesCommand_DetectsRuntimes(t *testing.T) {
	buf := new(bytes.Buffer)
	RuntimesCmd.SetOut(buf)
	RuntimesCmd.SetArgs([]string{})
	err := RuntimesCmd.Execute()
	if err != nil {
		t.Errorf("Runtimes command failed: %v", err)
	}
	// Optionally, check output for expected strings
}
