package computersetup

type IMachineConfiguration interface {
	Name() string
	Validate() error
	Install() error
	Uninstall() error
}

type IMachineConfigurationProvider interface {
	GetConfiguration() []IMachineConfiguration
}
