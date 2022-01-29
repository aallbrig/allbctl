package model

type IMachineConfiguration interface {
	Name() string
	Validate() error
	Install() error
	Uninstall() error
}
