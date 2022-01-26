package status

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const fakeKey = "TESTING_allbctl_envvar_key"
const fakeValue = "dummy fake value"

func TestCheckerCanFindExistingEnvVar(t *testing.T) {
	err := os.Setenv(fakeKey, fakeValue)
	if err != nil {
		t.Error(err)
	}
	checker := NewEnvironmentVariableChecker()

	err, result := checker.Check(fakeKey)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, result.Exists)
	assert.True(t, result.Name == fakeKey)
	os.Unsetenv(fakeKey)
}

func TestCheckerCheckResult_CanOutputStringRepresentation(t *testing.T) {
	err := os.Setenv(fakeKey, fakeValue)
	if err != nil {
		t.Error(err)
	}
	checker := NewEnvironmentVariableChecker()

	err, result := checker.Check(fakeKey)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fmt.Sprintf("Name: %s\tExists: %t", fakeKey, true), result.String())

	os.Unsetenv(fakeKey)
}