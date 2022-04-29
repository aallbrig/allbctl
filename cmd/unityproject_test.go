package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func Test_ProjectCmd_CreatesSourceCodeDirectory(t *testing.T) {
	sut := NewUnityProjectCommand()
	sut.SetArgs([]string{
		"--project", "test-project",
		"--ignore-unity-commands",
	})

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	operatingSystem.CurrentWorkingDirectory = tempDir

	err = sut.Execute()

	if _, err := os.Stat(filepath.Join(tempDir, "test-project")); os.IsNotExist(err) {
		t.Fail()
	}
}

func Test_ProjectCmd_InitializesSourceCodeRepo(t *testing.T) {
	sut := NewUnityProjectCommand()
	sut.SetArgs([]string{
		"--project", "test-project",
		"--ignore-unity-commands",
	})

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	operatingSystem.CurrentWorkingDirectory = tempDir

	err = sut.Execute()

	if _, err := os.Stat(filepath.Join(tempDir, "test-project", ".git")); os.IsNotExist(err) {
		t.Log("Expected a .git directory to exist")
		t.Fail()
	}
}

func Test_ProjectCmd_CreatesDirectoryForUnity(t *testing.T) {
	sut := NewUnityProjectCommand()
	sut.SetArgs([]string{
		"--project", "test-project",
		"--ignore-unity-commands",
	})

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	operatingSystem.CurrentWorkingDirectory = tempDir

	err = sut.Execute()

	if _, err := os.Stat(filepath.Join(tempDir, "test-project", "unity", "test-project")); os.IsNotExist(err) {
		t.Log("Expected a .git directory to exist")
		t.Fail()
	}
}

func Test_ProjectCmd_CreatesNewUnityProject(t *testing.T) {
	sut := NewUnityProjectCommand()
	sut.SetArgs([]string{"--project", "test-project"})

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	operatingSystem.CurrentWorkingDirectory = tempDir

	err = sut.Execute()

	if _, err := os.Stat(filepath.Join(tempDir, "test-project", "unity", "test-project", "Library")); os.IsNotExist(err) {
		t.Log("Expected a unity project directory to exist")
		t.Fail()
	}
}

func Test_ProjectCmd_GitIgnoresCommonUnityFiles(t *testing.T) {
	sut := NewUnityProjectCommand()
	sut.SetArgs([]string{
		"--project", "test-project",
		"--ignore-unity-commands",
	})

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	operatingSystem.CurrentWorkingDirectory = tempDir

	err = sut.Execute()

	gitIgnoreFileLocation := filepath.Join(tempDir, "test-project", "unity", "test-project", ".gitignore")
	if _, err := os.Stat(gitIgnoreFileLocation); os.IsNotExist(err) {
		t.Log("Expected a .gitignore file to exist in the unity project directory")
		t.Fail()
	}
}
