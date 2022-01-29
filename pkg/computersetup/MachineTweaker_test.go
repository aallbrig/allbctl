package computersetup

import (
	"github.com/aallbrig/allbctl/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SpyMachineConfiguration struct {
	OnValidate  func() error
	OnInstall   func() error
	OnUninstall func() error
}

func (s SpyMachineConfiguration) Name() string {
	return "Spy Machine Configuration"
}

func (s SpyMachineConfiguration) Validate() error {
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
	spy := new(SpyMachineConfiguration)
	spy.OnValidate = func() error {
		called = true
		return nil
	}
	machineConfigs := []model.IMachineConfiguration{spy}
	sut := NewMachineTweaker(machineConfigs)

	_ = sut.CheckCurrentMachine()[0]

	assert.True(t, called)
}

// func TestTweaker_CanApplyValidMachineConfiguration(t *testing.T) {