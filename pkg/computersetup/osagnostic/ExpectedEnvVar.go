package osagnostic

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

func (e ExpectedEnvVar) Validate() (out *bytes.Buffer, err error) {
	out = bytes.NewBufferString("")
	_, exists := os.LookupEnv(e.Key)
	if !exists {
		_, _ = color.New(color.FgRed).Fprint(out, "MISSING")
		err = fmt.Errorf("envvar does not exist")
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "PRESENT")
	}
	out.WriteString(fmt.Sprintf(" %s", e.Name()))
	return out, err
}

func (e ExpectedEnvVar) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	_, err := e.Validate()
	if err == nil {
		return out, err
	}
	if e.OnInstall != nil {
		return out, e.OnInstall()
	}

	return out, fmt.Errorf("no install lambda defined for envvar %s", e.Key)
}

func (e ExpectedEnvVar) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	_, err := e.Validate()
	if err == nil {
		return out, err
	}

	if e.OnUninstall != nil {
		return out, e.OnUninstall()
	}

	return out, fmt.Errorf("no uninstall lambda defined for envvar %s", e.Key)
}
