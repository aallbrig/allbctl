package os_agnostic

import (
	"github.com/mitchellh/go-homedir"
	"runtime"
)

type OperatingSystem struct{}

func (o OperatingSystem) GetName() (err error, os string) {
	os = runtime.GOOS
	return
}

func (o OperatingSystem) HomeDir() (err error, path string) {
	path, err = homedir.Dir()
	return
}
