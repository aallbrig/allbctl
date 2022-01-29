package computersetup

import (
	"bytes"
	"fmt"
	"github.com/aallbrig/allbctl/pkg/model"
)

type MachineTweaker struct {
	MachineConfiguration []model.IMachineConfiguration
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
		err, _ := machineConfig.Validate()
		if err == nil {
			result.Valid = true
		}
		report = append(report, *result)
	}
	return report
}

func (t MachineTweaker) ApplyConfiguration() ([]error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	var errs []error
	for _, configuration := range t.MachineConfiguration {
		out.WriteString(fmt.Sprintf("Applying configuration: %s\n", configuration.Name()))
		err, validateOut := configuration.Validate()
		out.WriteString(validateOut.String() + "\n")

		if err != nil {
			err, installOut := configuration.Install()
			out.WriteString(installOut.String() + "\n")

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs, out
}

func NewMachineTweaker(configs []model.IMachineConfiguration) *MachineTweaker {
	return &MachineTweaker{
		MachineConfiguration: configs,
	}
}
