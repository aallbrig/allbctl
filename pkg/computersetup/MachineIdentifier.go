package computersetup

import (
	"runtime"

	"github.com/aallbrig/allbctl/pkg/computersetup/providers"
	"github.com/aallbrig/allbctl/pkg/model"
)

type MachineIdentifier struct{}

func (m MachineIdentifier) ConfigurationProviderForOperatingSystem(os string) model.IMachineConfigurationProvider {
	switch os {
	case "windows":
		return providers.WindowsConfigurationProvider{}
	case "darwin":
		return providers.MacbookConfigurationProvider{}
	case "linux":
		return providers.LinuxConfigurationProvider{}
	default:
		return nil
	}
}

func (m MachineIdentifier) GetCurrentOperatingSystem() string {
	return runtime.GOOS
}
