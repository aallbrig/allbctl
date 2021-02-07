package computersetup

import (
	"fmt"
	"os"
	"path/filepath"
)

// TODO: .allbctl config should drive these
var sourceCodeDirectoryName = "src"
var userBinDirectoryName = "bin"

func createDirectory(homeDir string, desiredDirectory string) (err error) {
	dirFilepath := filepath.Join(homeDir, desiredDirectory)

	if stat, statErr := os.Stat(dirFilepath); statErr != nil && !os.IsNotExist(statErr) {
		err = statErr
	} else if stat != nil && !stat.IsDir() {
		err = fmt.Errorf("directory %s cannot be created due to conflict", dirFilepath)
	} else if stat == nil {
		err = os.Mkdir(dirFilepath, 0755)
	}

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
