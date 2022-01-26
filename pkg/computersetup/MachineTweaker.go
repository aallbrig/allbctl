package computersetup

type MachineTweaker struct {
	MachineConfiguration []IMachineConfiguration
}

func (t MachineTweaker) CheckCurrentMachine() []ValidateResult {
	var report []ValidateResult
	for _, machineConfig := range t.MachineConfiguration {
		result := machineConfig.Validate()
		report = append(report, result)
	}
	return report
}

func NewMachineTweaker(configs []IMachineConfiguration) *MachineTweaker {
	return &MachineTweaker{
		MachineConfiguration: configs,
	}
}
