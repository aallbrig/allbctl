package computersetup

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateSrcDirWhenNotExist(t *testing.T) {
	// Arrange
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal("[Testing] Error creating a temporary directory")
		t.Fail()
	}
	defer os.RemoveAll(dir)

	// Execute
	err = DirectoryForSourceCode(dir)
	if err != nil {
		log.Fatal("[Testing] Error source code directory", err)
		t.Fail()
	}

	// Assert
	srcDir, err := os.Stat(filepath.Join(dir, "src"))
	if err != nil && os.IsNotExist(err) {
		t.Fail()
	}

	assert.Equal(t, srcDir.Name(), "src")
	assert.Equal(t, srcDir.IsDir(), true)
	assert.Equal(t, srcDir.Mode(), os.ModeDir)
}

func TestCreateSrcDirWhenExist(t *testing.T) {
	// Arrange
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(dir)

	err = os.Mkdir(filepath.Join(dir, "src"), os.ModeDir)
	if err != nil {
		t.Error("[Testing] Error creating preexisting src directory")
	}

	// Execute
	err = DirectoryForSourceCode(dir)
	if err != nil {
		t.Error("[Testing] Error source code directory")
	}

	// Assert
	srcDir, err := os.Stat(filepath.Join(dir, "src"))
	if err != nil && os.IsNotExist(err) {
		t.Error("[Testing] Error getting stat on directory")
	}

	assert.Equal(t, srcDir.Name(), "src")
	assert.Equal(t, srcDir.IsDir(), true)
	assert.Equal(t, srcDir.Mode(), os.ModeDir)
}

func TestCreateUserBinWhenNotExist(t *testing.T) {
	// Arrange
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	// Execute
	err = DirectoryForUserBin(tempDir)
	if err != nil {
		t.Error("[Testing] Error source code directory")
	}

	// Assert
	srcDir, err := os.Stat(filepath.Join(tempDir, "bin"))
	if err != nil && os.IsNotExist(err) {
		t.Error("[Testing] Error getting stat on directory")
	}
	if srcDir == nil {
		t.Error("[Testing] Error getting stat on directory")
	}
	assert.Equal(t, srcDir.Name(), "bin")
	assert.Equal(t, srcDir.IsDir(), true)
	assert.Equal(t, srcDir.Mode(), os.ModeDir)
}

func TestCreateUserBinWhenAlreadyExist(t *testing.T) {
	// Arrange
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("[Testing] Error creating a temporary directory")
	}
	defer os.RemoveAll(tempDir)

	err = os.Mkdir(filepath.Join(tempDir, "bin"), os.ModeDir)
	if err != nil {
		t.Error("[Testing] Error creating preexisting src directory")
	}

	// Execute
	err = DirectoryForUserBin(tempDir)
	if err != nil {
		t.Error("[Testing] Error source code directory")
	}

	// Assert
	srcDir, err := os.Stat(filepath.Join(tempDir, "bin"))
	if err != nil && os.IsNotExist(err) {
		t.Error("[Testing] Error getting stat on directory")
	}
	if srcDir == nil {
		t.Error("[Testing] Error getting stat on directory")
	}
	assert.Equal(t, srcDir.Name(), "bin")
	assert.Equal(t, srcDir.IsDir(), true)
	assert.Equal(t, srcDir.Mode(), os.ModeDir)
}
