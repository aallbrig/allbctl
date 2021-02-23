package status

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_MissingStatusForDirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	stringBuf := bytes.NewBufferString("")
	err = CheckForDirectory(stringBuf, tempDir, "src")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(stringBuf.String(), "Missing") {
		t.Fail()
	}
}

func Test_PresentStatusForDirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	err = os.Mkdir(filepath.Join(tempDir, "src"), 0755)
	if err != nil {
		t.Error("[Testing] Error creating preexisting directory")
	}

	stringBuf := bytes.NewBufferString("")
	err = CheckForDirectory(stringBuf, tempDir, "src")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(stringBuf.String(), "Present") {
		t.Fail()
	}
}
