package os_agnostic

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
)

type EnvVarInstall func() error
type EnvVarUninstall func() error
type ExpectedEnvVar struct {
	Key         string
	OnInstall   EnvVarInstall
	OnUninstall EnvVarUninstall
}

func (e ExpectedEnvVar) Name() string {
	return fmt.Sprintf("Envvar %s", e.Key)
}

func (e ExpectedEnvVar) Validate() (err error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")
	_, exists := os.LookupEnv(e.Key)
	if !exists {
		_, _ = color.New(color.FgRed).Fprint(out, "MISSING")
		err = fmt.Errorf("envvar does not exist")
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "PRESENT")
	}
	out.WriteString(fmt.Sprintf(" %s", e.Name()))
	return err, out
}

func (e ExpectedEnvVar) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, _ := e.Validate()
	if err == nil {
		return err, out
	}
	if e.OnInstall != nil {
		return e.OnInstall(), out
	}

	return fmt.Errorf("no install lambda defined for envvar %s", e.Key), out
}

func (e ExpectedEnvVar) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, _ := e.Validate()
	if err == nil {
		return err, out
	}

	if e.OnUninstall != nil {
		return e.OnUninstall(), out
	}

	return fmt.Errorf("no uninstall lambda defined for envvar %s", e.Key), out
}
