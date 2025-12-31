package osagnostic

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpectedDirectory_Validate_NonExistentDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist")

	ed := NewExpectedDirectory(nonExistentPath)
	out, err := ed.Validate()

	if err != nil {
		t.Errorf("Validate() should not return error for non-existent directory, got: %v", err)
	}

	output := out.String()
	if strings.Contains(output, "PRESENT") {
		t.Errorf("Validate() should not report PRESENT for non-existent directory, got: %s", output)
	}
	if !strings.Contains(output, "NOT FOUND") && !strings.Contains(output, "MISSING") {
		t.Errorf("Validate() should report NOT FOUND or MISSING for non-existent directory, got: %s", output)
	}
}

func TestExpectedDirectory_Validate_ExistingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	existingDir := filepath.Join(tmpDir, "exists")
	if err := os.Mkdir(existingDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	ed := NewExpectedDirectory(existingDir)
	out, err := ed.Validate()

	if err != nil {
		t.Errorf("Validate() should not return error for existing directory, got: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "PRESENT") {
		t.Errorf("Validate() should report PRESENT for existing directory, got: %s", output)
	}
}

func TestExpectedDirectory_Validate_FileInsteadOfDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ed := NewExpectedDirectory(filePath)
	out, err := ed.Validate()

	if err == nil {
		t.Error("Validate() should return error when path is a file, not a directory")
	}

	output := out.String()
	if !strings.Contains(output, "expected directory is file") {
		t.Errorf("Validate() should report 'expected directory is file', got: %s", output)
	}
}
