package osagnostic

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"runtime"
	"strings"
)

func NewOperatingSystem() *OperatingSystem {
	operatingSystem := &OperatingSystem{}
	// Load up the data
	operatingSystem.setName()
	operatingSystem.setHomeDirectory()
	operatingSystem.setCurrentWorkingDirectory("")
	operatingSystem.setEnvironmentVariables()
	return operatingSystem
}

type EnvVar struct {
	Key   string
	Value string
}

type OperatingSystem struct {
	Name                    string
	HomeDirectoryPath       string
	CurrentWorkingDirectory string
	EnvironmentVariables    []*EnvVar
}

func (o *OperatingSystem) setName() {
	o.Name = runtime.GOOS
}

func (o *OperatingSystem) setHomeDirectory() {
	if path, err := homedir.Dir(); err == nil {
		o.HomeDirectoryPath = path
	}
}

func (o *OperatingSystem) validatePathInput(path string) bool {
	return true
}

func (o *OperatingSystem) UpdateCurrentWorkingDirectory(path string) {
	o.setCurrentWorkingDirectory(path)
}

func (o *OperatingSystem) setCurrentWorkingDirectory(path string) {
	if path != "" && o.validatePathInput(path) {
		o.UpdateCurrentWorkingDirectory(path)
	} else if path == "" {
		// acceptable condition -- get current working directory by default
		if path, err := os.Getwd(); err == nil {
			o.CurrentWorkingDirectory = path
		}
	}
}

func (o *OperatingSystem) CreateDirectory(path string) {
	// Does the directory already exist?
	//nolint:errcheck // TODO: what happens if this fails?
	_ = os.MkdirAll(path, os.ModePerm)
}

func (o *OperatingSystem) setEnvironmentVariables() {
	for _, envVar := range os.Environ() {
		// each envVar in Environ() takes the form "KEY=VALUE" e.g. USER=anonymous
		pair := strings.SplitN(envVar, "=", 2)
		o.EnvironmentVariables = append(o.EnvironmentVariables, &EnvVar{
			Key:   pair[0],
			Value: pair[1],
		})
	}
}
