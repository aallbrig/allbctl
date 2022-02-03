package macbook

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"regexp"
	"strings"
)

func RestartSystemUIServer() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	cmd := exec.Command("killall", "SystemUIServer")
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	return err, out
}

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
	out.WriteString(fmt.Sprintf("macos defaults %s %s", d.Domain, d.Key))
	err, currentValue := d.ReadCurrentValue()
	unquoted := regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(currentValue, `$1`)

	out.WriteString(" .... ")
	if unquoted == d.ExpectedValue || currentValue == d.ExpectedValue {
		_, _ = color.New(color.FgGreen).Fprint(out, "VALID")
	} else {
		_, _ = color.New(color.FgRed).Fprint(out, "INCORRECT")
		err = errors.New("current value is not expected value")
	}
	out.WriteString("\n")

	return err, out
}

func (d DefaultsCommand) WriteExpectedValue() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	validationError, validateOut := d.Validate()

	// already valid -- no need to reinstall
	if validationError == nil {
		out.WriteString(validateOut.String() + "\n")
		return nil, out
	}

	out.WriteString(fmt.Sprintf("Writing %s %s %s", d.Domain, d.Key, d.ExpectedValue))
	cmd := exec.Command(
		"defaults",
		"write",
		d.Domain,
		fmt.Sprintf("%s", d.Key),
		d.ValueType.String(),
		d.ExpectedValue,
	)
	err := cmd.Run()
	if err != nil {
		out.WriteString(" ❌ failed\n")
		return err, out
	}
	out.WriteString(" ✅ success\n")

	return nil, out
}

func (d DefaultsCommand) Delete() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")

	cmd := exec.Command("defaults", "delete", d.Domain, d.Key)
	_ = cmd.Run()

	return nil, out
}
