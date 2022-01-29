package computersetup

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/providers"
	"github.com/aallbrig/allbctl/pkg/model"
	"runtime"
)

type MachineIdentifier struct{}

func (m MachineIdentifier) ConfigurationProviderForOperatingSystem(os string) model.IMachineConfigurationProvider {
	switch os {
	case "windows":
		return nil
	case "darwin":
		return providers.MacbookConfigurationProvider{}
	case "linux":
		return nil
	default:
		return nil
	}
}

func (m MachineIdentifier) GetCurrentOperatingSystem() string {
	return runtime.GOOS
}
