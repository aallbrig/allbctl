package computersetup

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testCase struct {
	dirName        string
	functionToTest func(string) (err error)
	testExisting   bool
}

func handleCreateDirectoriesTestCase(tc testCase, t *testing.T) {
	// Arrange
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	if tc.testExisting {
		err = os.Mkdir(filepath.Join(tempDir, tc.dirName), 0755)
		if err != nil {
			t.Error("[Testing] Error creating preexisting directory")
		}
	}

	// Execute
	err = tc.functionToTest(tempDir)
	if err != nil {
		t.Error("[Testing] Error executing test function")
	}

	// Assert
	srcDir, err := os.Stat(filepath.Join(tempDir, tc.dirName))
	if err != nil && os.IsNotExist(err) {
		t.Error("[Testing] Error getting stat on directory")
	}

	if srcDir == nil {
		t.Error("[Testing] Error getting stat on directory")
	}

	assert.Equal(t, srcDir.Name(), tc.dirName)
	assert.Equal(t, srcDir.IsDir(), true)
	assert.Equal(t, srcDir.Mode().String(), "drwxr-xr-x")
}

func TestCreateDirectoriesInHome(t *testing.T) {
	testCases := []testCase{
		{"src", DirectoryForSourceCode, false},
		{"src", DirectoryForSourceCode, true},
		{"bin", DirectoryForUserBin, false},
		{"bin", DirectoryForUserBin, true},
	}

	for _, tc := range testCases {
		handleCreateDirectoriesTestCase(tc, t)
	}
}
