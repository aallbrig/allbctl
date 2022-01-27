package macbook

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type DefaultsCommand struct {
	Domain        string
	Key           string
	ExpectedValue string
	CurrentValue  string
	DefaultValue  string
}

func NewDefaultsCommand(domain string, key string) *DefaultsCommand {
	d := &DefaultsCommand{Domain: domain, Key: key}
	_ = d.SyncCurrentValue()
	return d
}

func (d DefaultsCommand) SyncCurrentValue() error {
	var stdOut bytes.Buffer

	cmd := exec.Command("defaults", "read", d.Domain, d.Key)
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		return err
	}

	d.CurrentValue = strings.TrimSuffix(stdOut.String(), "\n")
	return nil
}

func (d DefaultsCommand) Validate() error {
	err := d.SyncCurrentValue()

	if d.CurrentValue != d.ExpectedValue {
		err = errors.New("current value is not expected value")
	}

	return err
}

func (d DefaultsCommand) Install() error {
	validationError := d.Validate()

	// already valid -- no need to reinstall
	if validationError == nil {
		return nil
	}

	cmd := exec.Command("defaults", "write", d.Domain, d.Key, d.ExpectedValue)
	err := cmd.Run()
	if err != nil {
		return err
	}

	d.CurrentValue = d.ExpectedValue
	return nil
}

func (d DefaultsCommand) Uninstall() error {
	if d.CurrentValue == d.DefaultValue {
		return nil
	}

	cmd := exec.Command("defaults", "write", d.Domain, d.Key, d.DefaultValue)
	err := cmd.Run()
	if err != nil {
		return err
	}

	d.CurrentValue = d.ExpectedValue
	return nil
}
