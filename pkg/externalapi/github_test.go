package externalapi

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type MockGHAuthTokenProviderEmpty struct{}

func (authTokenProvider *MockGHAuthTokenProviderEmpty) GetAuthToken() (authToken string, err error) {
	authToken = os.Getenv("")

	if authToken == "" {
		err = fmt.Errorf("%s envvar not set", githubAuthTokenEnvVar)
	}

	return
}

type MockGHAuthTokenProviderWrong struct{}

func (authTokenProviderWrong *MockGHAuthTokenProviderWrong) GetAuthToken() (authToken string, err error) {
	authToken = ""
	return
}

func Test_ErrorWhenNoEnvVarSet(t *testing.T) {
	// ARRANGE
	authTokenProvider := MockGHAuthTokenProviderEmpty{}

	// EXECUTE
	_, err := GetMyDotfiles(&authTokenProvider)

	// ASSERT
	assert.NotEqual(t, nil, err, "Error should NOT be nil")
}

func Test_ErrorWhenIncorrectEnvVarSet(t *testing.T) {
	// ARRANGE
	authTokenProvider := MockGHAuthTokenProviderWrong{}

	// EXECUTE
	_, err := GetMyDotfiles(&authTokenProvider)

	// ASSERT
	assert.NotEqual(t, nil, err, "Error should NOT be nil")
	assert.Equal(t, 401, err.(*github.ErrorResponse).Response.StatusCode)
}

func Test_SuccessWhenCorrectEnvVarSet(t *testing.T) {
	// ARRANGE
	// Use the real provider (for now)
	authTokenProvider := GithubAuthTokenProvider{}

	// EXECUTE
	repos, err := GetMyDotfiles(&authTokenProvider)

	// ASSERT
	assert.Equal(t, nil, err, "No error should have returned")
	assert.Equal(t, 2, len(repos))
}
