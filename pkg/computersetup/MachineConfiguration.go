package computersetup

type ValidateResult struct {
	Name         string
	AlreadySetup bool
}
type IMachineConfiguration interface {
	Validate() ValidateResult
	Install() error
	Uninstall() error
}
