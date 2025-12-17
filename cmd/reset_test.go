package cmd

import (
	"bytes"
	"testing"
)

func TestResetCommand_Executes(t *testing.T) {
	buf := new(bytes.Buffer)
	ResetCmd.SetOut(buf)
	ResetCmd.SetArgs([]string{})
	err := ResetCmd.Execute()
	if err != nil {
		t.Errorf("Reset command failed: %v", err)
	}
	// The command should execute without error
	// Output validation is handled by the underlying components
}
