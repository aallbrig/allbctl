package osagnostic

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
)

type ExpectedCommand struct {
	CommandName string
}

func NewExpectedCommand(commandName string) *ExpectedCommand {
	return &ExpectedCommand{CommandName: commandName}
}

func (e ExpectedCommand) Name() string {
	return fmt.Sprintf("Expected Command: %s", e.CommandName)
}

func (e ExpectedCommand) Validate() (out *bytes.Buffer, err error) {
	out = bytes.NewBufferString("")

	_, err = exec.LookPath(e.CommandName)
	if err != nil {
		_, _ = color.New(color.FgRed).Fprint(out, "NOT FOUND")
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "INSTALLED")
	}
	out.WriteString(fmt.Sprintf(" %s", e.CommandName))

	return
}

func (e ExpectedCommand) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	out.WriteString(fmt.Sprintf("❌ Cannot auto-install %s - please install manually\n", e.CommandName))
	out.WriteString("Suggested installation methods:\n")
	out.WriteString("  - Linux: apt install / yum install / pacman -S\n")
	out.WriteString("  - macOS: brew install\n")
	out.WriteString("  - Windows: choco install / winget install\n")

	return out, fmt.Errorf("manual installation required for %s", e.CommandName)
}

func (e ExpectedCommand) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	out.WriteString(fmt.Sprintf("❌ Cannot auto-uninstall %s - please uninstall manually\n", e.CommandName))
	return out, nil
}
