package os_agnostic

import (
	"fmt"
	"os"
)

type ExpectedDirectory struct {
	Path       string
	Permission os.FileMode
}

func NewExpectedDirectory(path string) *ExpectedDirectory {
	return &ExpectedDirectory{Path: path, Permission: 0755}
}

func (e ExpectedDirectory) Name() string {
	return fmt.Sprintf("Expected Directory %s", e.Path)
}

func (e ExpectedDirectory) Validate() (err error) {
	if stat, statErr := os.Stat(e.Path); statErr != nil && !os.IsNotExist(statErr) {
		err = statErr
	} else if stat != nil && !stat.IsDir() {
		err = fmt.Errorf("directory %s cannot be created due to conflict", e.Path)
	}

	return
}

func (e ExpectedDirectory) Install() error {
	err := e.Validate()
	if err != nil {
		return err
	}

	err = os.Mkdir(e.Path, e.Permission)

	return err
}

func (e ExpectedDirectory) Uninstall() error {
	err := e.Validate()
	if err != nil {
		return nil
	}

	err = os.RemoveAll(e.Path)

	return err
}
