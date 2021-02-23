package status

import (
	"bytes"
	"os"
	"path/filepath"
)

func CheckForDirectory(w *bytes.Buffer, dir string, dirToFind string) (err error) {
	targetDirectory := filepath.Join(dir, dirToFind)
	w.WriteString(targetDirectory)

	dirStat, err := os.Stat(targetDirectory)
	if os.IsNotExist(err) {
		w.WriteString("Missing")
	} else if dirStat.IsDir() {
		w.WriteString("Present")
	} else {
		w.WriteString("Not Directory")
	}

	return
}
