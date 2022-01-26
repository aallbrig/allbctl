package computersetup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type SpyMachineConfiguration struct {
	OnValidate  func() ValidateResult
	OnInstall   func() error
	OnUninstall func() error
}

func (s SpyMachineConfiguration) Validate() ValidateResult {
	return s.OnValidate()
}

func (s SpyMachineConfiguration) Install() error {
	return s.OnInstall()
}

func (s SpyMachineConfiguration) Uninstall() error {
	return s.OnUninstall()
}

func TestTweaker_CanReport(t *testing.T) {
	called := false
	name := "test_name"
	exists := true
	spy := new(SpyMachineConfiguration)
	spy.OnValidate = func() ValidateResult {
		called = true
		return ValidateResult{name, exists}
	}
	machineConfigs := []IMachineConfiguration{spy}
	sut := NewMachineTweaker(machineConfigs)

	result := sut.CheckCurrentMachine()[0]

	assert.True(t, called)
	assert.Equal(t, name, result.Name)
}

// func TestTweaker_CanApplyValidMachineConfiguration(t *testing.T) {
