package computersetup

import (
	"errors"
	"os"
	"path/filepath"
)

// TODO: .allbctl config should drive this
var sourceCodeDirectoryName = "src"

// DirectoryForSourceCode used to create directory for personal source code
func DirectoryForSourceCode(homeDir string) (err error) {
	srcDirFilePath := filepath.Join(homeDir, sourceCodeDirectoryName)
	srcDirStat, err := os.Stat(srcDirFilePath)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if srcDirStat != nil && srcDirStat.IsDir() {
		// Source code directory already exists, no need for work
		return
	} else if srcDirStat != nil && !srcDirStat.IsDir() {
		err = errors.New("desired source code directory cannot be created due to conflicting file")
		return
	}

	err = os.Mkdir(srcDirFilePath, os.ModeDir)
	return
}
