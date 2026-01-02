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

// FilterOutSSHKeyRegistration removes SSH key registration from configuration list
func FilterOutSSHKeyRegistration(configs []model.IMachineConfiguration) []model.IMachineConfiguration {
	var filtered []model.IMachineConfiguration
	for _, config := range configs {
		// Check if it's a configuration group
		if group, ok := config.(model.MachineConfigurationGroup); ok {
			// Filter out SSH Configuration group
			if group.GroupName == "SSH Configuration" {
				continue
			}
			// Also check nested configs
			var newConfigs []model.IMachineConfiguration
			for _, nestedConfig := range group.Configs {
				if nestedConfig.Name() != "SSH Key GitHub Registration" {
					newConfigs = append(newConfigs, nestedConfig)
				}
			}
			if len(newConfigs) > 0 {
				group.Configs = newConfigs
				filtered = append(filtered, group)
			}
		} else if config.Name() != "SSH Key GitHub Registration" {
			filtered = append(filtered, config)
		}
	}
	return filtered
}
