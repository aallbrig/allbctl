package computersetup

import (
	"bytes"
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

func (s SpyMachineConfiguration) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	return s.OnValidate(), out
}

func (s SpyMachineConfiguration) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	return s.OnInstall(), out
}

func (s SpyMachineConfiguration) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	return s.OnUninstall(), out
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

	_, _ = sut.ConfigurationStatus()

	assert.True(t, called)
}

// func TestTweaker_CanApplyValidMachineConfiguration(t *testing.T) {
