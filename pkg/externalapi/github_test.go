package externalapi

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type MockGHAuthTokenProviderEmpty struct{}

func (provider *MockGHAuthTokenProviderEmpty) GetAuthToken() (authToken string, err error) {
	authToken = os.Getenv("")

	if authToken == "" {
		err = fmt.Errorf("%s envvar not set", githubAuthTokenEnvVar)
	}

	return
}

type MockGHAuthTokenProviderFake struct{}

func (provider *MockGHAuthTokenProviderFake) GetAuthToken() (authToken string, err error) {
	authToken = "foobar"
	return
}

type MockUsersService struct{}

func (userService *MockUsersService) Get(_ context.Context, _ string) (user *github.User,
	response *github.Response, err error) {

	login := "mock user"
	user = &github.User{
		Login: &login,
	}

	return
}

type MockSearchService struct{}

func (searchService *MockSearchService) Repositories(_ context.Context, _ string,
	_ *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	return &github.RepositoriesSearchResult{
		Repositories: nil,
	}, nil, nil
}

type MockGHProvider struct{}

func (provider *MockGHProvider) GetGithubClient(_ string) (client GithubClient, err error) {
	mockUsersService := MockUsersService{}
	mockSearchService := MockSearchService{}

	client = GithubClient{
		Users:  &mockUsersService,
		Search: &mockSearchService,
	}
	return
}

func Test_ErrorWhenNoEnvVarSet(t *testing.T) {
	// ARRANGE
	tokenProvider = &MockGHAuthTokenProviderEmpty{}
	ghProvider = &MockGHProvider{}

	// EXECUTE
	_, err := GetMyDotfiles()

	// ASSERT
	assert.NotEqual(t, nil, err, "Error should NOT be nil")
}

func Test_SuccessWhenCorrectEnvVarSet(t *testing.T) {
	// ARRANGE
	tokenProvider = &MockGHAuthTokenProviderFake{}
	ghProvider = &MockGHProvider{}

	// EXECUTE
	_, err := GetMyDotfiles()

	// ASSERT
	assert.Equal(t, nil, err, "No error should have returned")
}
