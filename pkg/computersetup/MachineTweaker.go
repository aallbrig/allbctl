package computersetup

import (
	"fmt"
	"log"
)

type MachineTweaker struct {
	MachineConfiguration []IMachineConfiguration
}
type ValidateResult struct {
	Name  string
	Valid bool
}

func (t MachineTweaker) CheckCurrentMachine() []ValidateResult {
	var report []ValidateResult
	for _, machineConfig := range t.MachineConfiguration {
		result := &ValidateResult{
			Name:  machineConfig.Name(),
			Valid: false,
		}
		err := machineConfig.Validate()
		if err != nil {
			result.Valid = true
		}
		report = append(report)
	}
	return report
}

func (t MachineTweaker) ApplyConfiguration() []error {
	var errs []error
	for _, configuration := range t.MachineConfiguration {
		log.Println(fmt.Sprintf("Applying configuration for %s", configuration.Name()))
		err := configuration.Validate()
		if err != nil {
			err = configuration.Install()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func NewMachineTweaker(configs []IMachineConfiguration) *MachineTweaker {
	return &MachineTweaker{
		MachineConfiguration: configs,
	}
}
