package computersetup

import (
	"runtime"
)

type MachineIdentifier struct{}

func (m MachineIdentifier) ConfigurationForMachine() IMachineConfigurationProvider {
	os := runtime.GOOS

	switch os {
	case "windows":
		return nil
	case "darwin":
		return MacbookConfigurationProvider{}
	case "linux":
		return nil
	default:
		return nil
	}
}
