package computersetup

import (
	"bytes"
	"fmt"
	"github.com/aallbrig/allbctl/pkg/model"
)

type MachineTweaker struct {
	MachineConfiguration []model.IMachineConfiguration
}

func (t MachineTweaker) ApplyConfiguration() ([]error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	var errs []error
	for _, configuration := range t.MachineConfiguration {
		out.WriteString(fmt.Sprintf("Applying configuration: %s\n", configuration.Name()))
		validateOut, err := configuration.Validate()
		out.WriteString(validateOut.String() + "\n")

		if err != nil {
			installOut, err := configuration.Install()
			out.WriteString(installOut.String() + "\n")

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs, out
}

func (t MachineTweaker) ConfigurationStatus() (errs []error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")

	for _, configuration := range t.MachineConfiguration {
		out.WriteString(fmt.Sprintf("%s\n-----\n", configuration.Name()))
		validateOut, err := configuration.Validate()
		out.WriteString(validateOut.String() + "\n")

		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (t MachineTweaker) ResetConfiguration() (errs []error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")

	// Uninstall in reverse
	for i := len(t.MachineConfiguration) - 1; i >= 0; i-- {
		configuration := t.MachineConfiguration[i]
		uninstallOut, err := configuration.Uninstall()
		out.WriteString(uninstallOut.String() + "\n")

		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func NewMachineTweaker(configs []model.IMachineConfiguration) *MachineTweaker {
	return &MachineTweaker{
		MachineConfiguration: configs,
	}
}
