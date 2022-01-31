package macbook

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"strings"
)

type DefaultsType int

const (
	DefaultsBool DefaultsType = iota
	DefaultsString
	DefaultsInt
)

func (t DefaultsType) String() string {
	switch t {
	case DefaultsBool:
		return "-bool"
	case DefaultsString:
		return "-string"
	case DefaultsInt:
		return "-int"
	}
	return ""
}

func (t DefaultsType) EnumIndex() int {
	return int(t)
}

type UninstallDefaults func(d DefaultsCommand) (error, *bytes.Buffer)
type DefaultsCommand struct {
	Domain        string
	Key           string
	ExpectedValue string
	ValueType     DefaultsType
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
	out := bytes.NewBufferString("")
	err, currentValue := d.ReadCurrentValue()

	if err != nil || currentValue != d.ExpectedValue {
		_, _ = color.New(color.FgRed).Fprint(out, "INCORRECT")
		err = errors.New("current value is not expected value")
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "CORRECT")
	}

	out.WriteString(fmt.Sprintf(" defaults setting %s %s ", d.Domain, d.Key))

	return err, out
}

func (d DefaultsCommand) WriteExpectedValue() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	validationError, validateOut := d.Validate()

	// already valid -- no need to reinstall
	if validationError == nil {
		out.WriteString(validateOut.String())
		return nil, out
	}

	cmd := exec.Command(
		"defaults",
		"write",
		d.Domain,
		fmt.Sprintf("\"%s\"", d.Key),
		d.ValueType.String(),
		fmt.Sprintf("\"%s\"", d.ExpectedValue),
	)
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	if err != nil {
		return err, out
	}

	return nil, out
}

func (d DefaultsCommand) Delete() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")

	cmd := exec.Command("defaults", "delete", d.Domain, d.Key)
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	if err != nil {
		return err, out
	}

	return nil, out
}
