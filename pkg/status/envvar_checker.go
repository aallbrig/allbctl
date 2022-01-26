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
	return fmt.Sprintf("Name: %s\tExists: %t", result.Name, result.Exists)
}

type EnvironmentVariableChecker struct{}

func NewEnvironmentVariableChecker() *EnvironmentVariableChecker {
	return &EnvironmentVariableChecker{}
}

func (checker EnvironmentVariableChecker) Check(envVarName string) (err error, result *CheckEnvVarResult) {
	_, exists := os.LookupEnv(envVarName)
	return err, &CheckEnvVarResult{
		Name:   envVarName,
		Exists: exists,
	}
}
