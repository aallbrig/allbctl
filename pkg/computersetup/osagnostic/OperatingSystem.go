package osagnostic

import (
	"github.com/mitchellh/go-homedir"
	"runtime"
)

type OperatingSystem struct{}

func (o OperatingSystem) GetName() (os string, err error) {
	os = runtime.GOOS
	return
}

func (o OperatingSystem) HomeDir() (path string, err error) {
	path, err = homedir.Dir()
	return
}
