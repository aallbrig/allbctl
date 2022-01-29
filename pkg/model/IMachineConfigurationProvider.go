package model

type IMachineConfigurationProvider interface {
	GetConfiguration() []IMachineConfiguration
}
