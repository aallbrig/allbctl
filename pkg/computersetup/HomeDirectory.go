package computersetup

import (
	"errors"
	"os"
	"path/filepath"
)

// TODO: .allbctl config should drive these
var sourceCodeDirectoryName = "src"
var userBinDirectoryName = "bin"

func createDirectory(homeDir string, desiredDirectory string) (err error) {
	dirFilepath := filepath.Join(homeDir, desiredDirectory)

	// Does the directory already exist?
	stat, err := os.Stat(dirFilepath)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if stat != nil && stat.IsDir() {
		// Source code directory already exists, no need for work
		return
	} else if stat != nil && !stat.IsDir() {
		err = errors.New("desired source code directory cannot be created due to conflicting file")
		return
	}

	err = os.Mkdir(dirFilepath, os.ModeDir)
	return
}

// DirectoryForSourceCode used to create directory for personal source code
func DirectoryForSourceCode(homeDir string) (err error) {
	err = createDirectory(homeDir, sourceCodeDirectoryName)
	return
}

// DirectoryForUserBin used for user specific binaries (e.g. jetbrains IDE commands)
func DirectoryForUserBin(homeDir string) (err error) {
	err = createDirectory(homeDir, userBinDirectoryName)
	return
}
