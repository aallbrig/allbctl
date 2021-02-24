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
	_ = CheckForDirectory(stringBuf, tempDir, "src")

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
	_ = CheckForDirectory(stringBuf, tempDir, "src")

	if !strings.Contains(stringBuf.String(), "Present") {
		t.Fail()
	}
}

func Test_SystemInfo(t *testing.T) {
	testCases := []struct {
		goos     string
		expected string
	}{
		{"windows", "Windows"},
		{"darwin", "MAC OS"},
		{"linux", "Linux"},
		{"random", "random"},
	}

	stringBuf := bytes.NewBufferString("")
	oldGoos := goos

	for _, testCase := range testCases {
		goos = testCase.goos

		_ = SystemInfo(stringBuf)
		if !strings.Contains(stringBuf.String(), testCase.expected) {
			t.Errorf("Expected %s to contain %s", stringBuf.String(), testCase.expected)
		}
	}
	goos = oldGoos
}
