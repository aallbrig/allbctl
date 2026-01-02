package osagnostic

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
)

type ExpectedDirectory struct {
	Path       string
	Permission os.FileMode
}

func NewExpectedDirectory(path string) *ExpectedDirectory {
	return &ExpectedDirectory{Path: path, Permission: 0755}
}

func (e ExpectedDirectory) Name() string {
	return fmt.Sprintf("Expected Directory %s", e.Path)
}

func (e ExpectedDirectory) Validate() (out *bytes.Buffer, err error) {
	out = bytes.NewBufferString("")
	stat, statErr := os.Stat(e.Path)

	if statErr != nil {
		if os.IsNotExist(statErr) {
			_, _ = color.New(color.FgRed).Fprint(out, "NOT FOUND")
		} else {
			_, _ = color.New(color.FgRed).Fprint(out, "stat error")
			err = statErr
		}
	} else if !stat.IsDir() {
		_, _ = color.New(color.FgRed).Fprint(out, "expected directory is file")
		err = fmt.Errorf("directory %s cannot be created due to conflict", e.Path)
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "PRESENT")
	}
	out.WriteString(fmt.Sprintf(" %s", e.Path))

	return
}

func (e ExpectedDirectory) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")

	// Check if already exists
	validateOut, err := e.Validate()
	if err == nil {
		// Already exists
		out.WriteString(validateOut.String() + "\n")
		_, _ = color.New(color.FgGreen).Fprint(out, "✅ Directory already exists\n")
		return out, nil
	}

	out.WriteString(validateOut.String() + "\n")

	// Create directory with MkdirAll for idempotency
	err = os.MkdirAll(e.Path, e.Permission)
	if err != nil {
		_, _ = color.New(color.FgRed).Fprint(out, fmt.Sprintf("❌ Failed to create %s: %v\n", e.Path, err))
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, fmt.Sprintf("✅ Created directory %s\n", e.Path))
	}

	return out, err
}

func (e ExpectedDirectory) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	_, err := e.Validate()
	if err != nil {
		return out, nil
	}

	err = os.RemoveAll(e.Path)

	return out, err
}
