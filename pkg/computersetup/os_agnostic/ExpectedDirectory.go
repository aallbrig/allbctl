package os_agnostic

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

func (e ExpectedDirectory) Validate() (err error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")
	out.WriteString(fmt.Sprintf("%s ", e.Path))
	if stat, statErr := os.Stat(e.Path); statErr != nil && !os.IsNotExist(statErr) {
		_, _ = color.New(color.FgRed).Fprint(out, "stat error")
		err = statErr
	} else if stat != nil && !stat.IsDir() {
		_, _ = color.New(color.FgRed).Fprint(out, "expected directory is file")
		err = fmt.Errorf("directory %s cannot be created due to conflict", e.Path)
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "PRESENT")
	}

	return
}

func (e ExpectedDirectory) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, validateOut := e.Validate()
	out.WriteString(validateOut.String() + "\n")
	if err != nil {
		return err, out
	}

	err = os.Mkdir(e.Path, e.Permission)
	if err != nil {
		_, _ = color.New(color.FgRed).Fprint(out, fmt.Sprintf("Fail to create %s", e.Path))
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, fmt.Sprintf("Create success %s", e.Path))
	}

	return err, out
}

func (e ExpectedDirectory) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, _ := e.Validate()
	if err != nil {
		return nil, out
	}

	err = os.RemoveAll(e.Path)

	return err, out
}
