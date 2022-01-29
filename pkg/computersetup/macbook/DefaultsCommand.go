package macbook

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"strings"
)

type DefaultsCommand struct {
	Domain        string
	Key           string
	ExpectedValue string
	DefaultValue  string
}

func NewDefaultsCommand(domain string, key string) *DefaultsCommand {
	return &DefaultsCommand{Domain: domain, Key: key}
}

func (d DefaultsCommand) ReadCurrentValue() (error, string) {
	var stdOut bytes.Buffer

	cmd := exec.Command("defaults", "read", d.Domain, d.Key)
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		return err, ""
	}

	return nil, strings.TrimSuffix(stdOut.String(), "\n")
}

func (d DefaultsCommand) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString(fmt.Sprintf("Defaults setting %s %s ", d.Domain, d.Key))
	err, currentValue := d.ReadCurrentValue()
	if err != nil {
		return err, out
	}

	if currentValue != d.ExpectedValue {
		_, _ = color.New(color.FgRed).Fprint(out, "INCORRECT")
		err = errors.New("current value is not expected value")
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "CORRECT")
	}

	return err, out
}

func (d DefaultsCommand) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	validationError, _ := d.Validate()

	// already valid -- no need to reinstall
	if validationError == nil {
		return nil, out
	}

	cmd := exec.Command("defaults", "write", d.Domain, d.Key, d.ExpectedValue)
	err := cmd.Run()
	if err != nil {
		return err, out
	}

	return nil, nil
}

func (d DefaultsCommand) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")

	cmd := exec.Command("defaults", "write", d.Domain, d.Key, d.DefaultValue)
	err := cmd.Run()
	if err != nil {
		return err, out
	}

	return nil, out
}
