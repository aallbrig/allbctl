package model

import (
	"bytes"
	"fmt"
)
import "github.com/pkg/errors"

type MachineConfigurationGroup struct {
	GroupName string
	Configs   []IMachineConfiguration
}

func (m MachineConfigurationGroup) Name() string {
	return m.GroupName
}

func (m MachineConfigurationGroup) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	if len(m.Configs) == 0 {
		return fmt.Errorf("%s No configuration for section", m.GroupName), out
	}

	var errs []error
	for _, config := range m.Configs {
		err, validateOut := config.Validate()
		out.WriteString(validateOut.String() + "\n")

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return WrapErrors(errors.New(m.GroupName), errs), out
	}

	return nil, out
}

func (m MachineConfigurationGroup) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	if len(m.Configs) == 0 {
		return fmt.Errorf("%s No configuration for section", m.GroupName), out
	}

	var errs []error
	for _, config := range m.Configs {
		err, validateOut := config.Validate()

		if err != nil {
			out.WriteString(fmt.Sprintf("\tInstalling: %s\n", config.Name()))
			err, installOut := config.Install()
			out.WriteString(installOut.String())

			if err != nil {
				errs = append(errs, err)
			}
		} else {
			out.WriteString(validateOut.String())
		}
	}

	if len(errs) != 0 {
		return WrapErrors(errors.New(m.GroupName), errs), out
	}

	return nil, out
}

func (m MachineConfigurationGroup) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	if len(m.Configs) == 0 {
		return fmt.Errorf("%s No configuration for section", m.GroupName), out
	}

	var errs []error
	for i := len(m.Configs) - 1; i >= 0; i-- {
		config := m.Configs[i]
		out.WriteString(fmt.Sprintf("\tUninstalling: %s\n", config.Name()))

		err, uninstallOut := config.Uninstall()
		out.WriteString(uninstallOut.String())

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return WrapErrors(errors.New(m.GroupName), errs), out
	}

	return nil, out
}

func WrapErrors(err error, errs []error) error {
	for _, err2 := range errs {
		err = errors.Wrap(err, err2.Error())
	}
	return err
}
