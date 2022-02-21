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

func (m MachineConfigurationGroup) Validate() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	if len(m.Configs) == 0 {
		return out, fmt.Errorf("%s No configuration for section", m.GroupName)
	}

	var errs []error
	for _, config := range m.Configs {
		validateOut, err := config.Validate()
		out.WriteString(validateOut.String() + "\n")

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return out, WrapErrors(errors.New(m.GroupName), errs)
	}

	return out, nil
}

func (m MachineConfigurationGroup) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	if len(m.Configs) == 0 {
		return out, fmt.Errorf("%s No configuration for section", m.GroupName)
	}

	var errs []error
	for _, config := range m.Configs {
		validateOut, err := config.Validate()

		if err != nil {
			out.WriteString(fmt.Sprintf("\tInstalling: %s\n", config.Name()))
			installOut, err := config.Install()
			out.WriteString(installOut.String())

			if err != nil {
				errs = append(errs, err)
			}
		} else {
			out.WriteString(validateOut.String())
		}
	}

	if len(errs) != 0 {
		return out, WrapErrors(errors.New(m.GroupName), errs)
	}

	return out, nil
}

func (m MachineConfigurationGroup) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	if len(m.Configs) == 0 {
		return out, fmt.Errorf("%s No configuration for section", m.GroupName)
	}

	var errs []error
	for i := len(m.Configs) - 1; i >= 0; i-- {
		config := m.Configs[i]
		out.WriteString(fmt.Sprintf("\tUninstalling: %s\n", config.Name()))

		uninstallOut, err := config.Uninstall()
		out.WriteString(uninstallOut.String())

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return out, WrapErrors(errors.New(m.GroupName), errs)
	}

	return out, nil
}

func WrapErrors(err error, errs []error) error {
	for _, err2 := range errs {
		err = errors.Wrap(err, err2.Error())
	}
	return err
}
