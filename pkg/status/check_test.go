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

func Test_PackageManagerPresent(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	path := os.Getenv("PATH")
	_ = os.Setenv("PATH", tempDir)

	testCases := []struct {
		goos     string
		expected string
		command  string
	}{
		{"windows", "Chocolatey", "choco"},
		{"darwin", "Homebrew", "brew"},
		{"linux", "Apt-Get", "apt-get"},
	}

	oldGoos := goos

	for _, testCase := range testCases {
		stringBuf := bytes.NewBufferString("")
		err = ioutil.WriteFile(filepath.Join(tempDir, testCase.command), []byte(""), 0777)
		if err != nil {
			t.Error("[Testing] Error creating file")
		}

		goos = testCase.goos
		_ = PackageManager(stringBuf)
		if !strings.Contains(stringBuf.String(), testCase.expected) {
			t.Errorf("Expected %s to contain package manager %s", stringBuf.String(), testCase.expected)
		}
		if !strings.Contains(stringBuf.String(), "Present") {
			t.Errorf("Expected %s to contain package manager command %s", stringBuf.String(), testCase.command)
		}
	}

	_ = os.Setenv("PATH", path)
	goos = oldGoos
}

func Test_PackageManagerMissing(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	path := os.Getenv("PATH")
	_ = os.Setenv("PATH", tempDir)

	testCases := []struct {
		goos     string
		expected string
		command  string
	}{
		{"windows", "Chocolatey", "choco"},
		{"darwin", "Homebrew", "brew"},
		{"linux", "Apt-Get", "apt-get"},
	}

	oldGoos := goos

	for _, testCase := range testCases {
		stringBuf := bytes.NewBufferString("")

		goos = testCase.goos
		_ = PackageManager(stringBuf)

		if !strings.Contains(stringBuf.String(), "Missing") {
			t.Errorf("Expected %s to contain package manager command %s", stringBuf.String(), testCase.command)
		}
	}

	_ = os.Setenv("PATH", path)
	goos = oldGoos
}
