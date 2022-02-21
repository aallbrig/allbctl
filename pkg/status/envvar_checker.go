package status

import (
	"fmt"
	"os"
)

type CheckEnvVarResult struct {
	Name   string
	Exists bool
}

func (result CheckEnvVarResult) String() string {
	return fmt.Sprintf("ConfigName: %s\tExists: %t", result.Name, result.Exists)
}

type EnvironmentVariableChecker struct{}

func NewEnvironmentVariableChecker() *EnvironmentVariableChecker {
	return &EnvironmentVariableChecker{}
}

func (checker EnvironmentVariableChecker) Check(envVarName string) (result *CheckEnvVarResult, err error) {
	_, exists := os.LookupEnv(envVarName)
	return &CheckEnvVarResult{
		Name:   envVarName,
		Exists: exists,
	}, err
}
