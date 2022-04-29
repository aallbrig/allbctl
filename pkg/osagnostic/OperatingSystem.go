package osagnostic

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"runtime"
)

func NewOperatingSystem() *OperatingSystem {
	operatingSystem := &OperatingSystem{}
	// Load up the data
	operatingSystem.setName()
	operatingSystem.setHomeDirectory()
	operatingSystem.setCurrentWorkingDirectory()
	return operatingSystem
}

type OperatingSystem struct {
	Name                    string
	HomeDirectoryPath       string
	CurrentWorkingDirectory string
}

func (o *OperatingSystem) setName() {
	o.Name = runtime.GOOS
}

func (o *OperatingSystem) setHomeDirectory() {
	if path, err := homedir.Dir(); err == nil {
		o.HomeDirectoryPath = path
	}
}

func (o *OperatingSystem) setCurrentWorkingDirectory() {
	if path, err := os.Getwd(); err == nil {
		o.CurrentWorkingDirectory = path
	}
}

func (o *OperatingSystem) CreateDirectory(path string) {
	// Does the directory already exist?
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		// TODO: what happens if this fails?
	}
}
