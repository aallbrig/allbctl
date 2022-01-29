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
		out.WriteString(validateOut.String())

		if err != nil {
			err, installOut := config.Install()
			out.WriteString(installOut.String())

			if err != nil {
				errs = append(errs, err)
			}
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
	for _, config := range m.Configs {
		err, validateOut := config.Validate()
		out.WriteString(validateOut.String())

		if err != nil {
			err, installOut := config.Install()
			out.WriteString(installOut.String())

			if err != nil {
				errs = append(errs, err)
			}
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
