package os_agnostic

import (
	"fmt"
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
	return fmt.Sprintf("Expected envvar %s", e.Key)
}

func (e ExpectedEnvVar) Validate() (err error) {
	_, exists := os.LookupEnv(e.Key)
	if !exists {
		err = fmt.Errorf("envvar does not exist")
	}
	return err
}

func (e ExpectedEnvVar) Install() error {
	err := e.Validate()
	if err == nil {
		return err
	}
	if e.OnInstall != nil {
		return e.OnInstall()
	}

	return fmt.Errorf("no install lambda defined for envvar %s", e.Key)
}

func (e ExpectedEnvVar) Uninstall() error {
	err := e.Validate()
	if err == nil {
		return err
	}

	if e.OnUninstall != nil {
		return e.OnUninstall()
	}

	return fmt.Errorf("no uninstall lambda defined for envvar %s", e.Key)
}
