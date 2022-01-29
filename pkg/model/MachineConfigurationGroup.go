package model

import (
	"fmt"
	"log"
)
import "github.com/pkg/errors"

type MachineConfigurationGroup struct {
	GroupName string
	Configs   []IMachineConfiguration
}

func (m MachineConfigurationGroup) Name() string {
	return m.GroupName
}

func (m MachineConfigurationGroup) Validate() error {
	if len(m.Configs) == 0 {
		return fmt.Errorf("%s No configuration for section", m.GroupName)
	}

	var errs []error
	for _, config := range m.Configs {
		err := config.Validate()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return WrapErrors(errors.New(m.GroupName), errs)
	}

	return nil
}

func (m MachineConfigurationGroup) Install() error {
	if len(m.Configs) == 0 {
		return fmt.Errorf("%s No configuration for section", m.GroupName)
	}

	var errs []error
	for _, config := range m.Configs {
		if err := config.Validate(); err != nil {
			log.Println(fmt.Sprintf("Installing configuration for %s", config.Name()))
			if err := config.Install(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) != 0 {
		return WrapErrors(errors.New(m.GroupName), errs)
	}

	return nil
}

func (m MachineConfigurationGroup) Uninstall() error {
	if len(m.Configs) == 0 {
		return fmt.Errorf("%s No configuration for section", m.GroupName)
	}

	var errs []error
	for _, config := range m.Configs {
		if err := config.Validate(); err == nil {
			log.Println(fmt.Sprintf("Uninstalling configuration for %s", config.Name()))
			if err := config.Uninstall(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) != 0 {
		return WrapErrors(errors.New(m.GroupName), errs)
	}

	return nil
}

func WrapErrors(err error, errs []error) error {
	for _, err2 := range errs {
		err = errors.Wrap(err, err2.Error())
	}
	return err
}
