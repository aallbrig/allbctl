package externalapi

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

var githubAuthTokenEnvVar = "GH_AUTH_TOKEN"

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func newGithubClient(githubAPIAuthToken string) (client *github.Client) {
	tokenSource := &tokenSource{
		AccessToken: githubAPIAuthToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client = github.NewClient(oauthClient)
	return
}

type githubAuthTokenProvider interface {
	GetAuthToken() (string, error)
}

// GithubAuthTokenProvider providers a way to get the GH auth token
type GithubAuthTokenProvider struct{}

// GetAuthToken gets github auth token
func (authTokenProvider *GithubAuthTokenProvider) GetAuthToken() (authToken string, err error) {
	authToken = os.Getenv(githubAuthTokenEnvVar)

	if authToken == "" {
		err = fmt.Errorf("%s envvar not set", githubAuthTokenEnvVar)
	}

	return
}

// GetMyDotfiles gets dotfile repos from my github
func GetMyDotfiles(authTokenProvider githubAuthTokenProvider) (repositories []github.Repository, err error) {
	ctx := context.TODO()
	githubAuthToken, err := authTokenProvider.GetAuthToken()

	if err != nil {
		return
	}

	client := newGithubClient(githubAuthToken)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return
	}

	reposResponse, _, err := client.Search.Repositories(ctx, fmt.Sprintf("dotfiles user:%s", *user.Login), nil)
	if err != nil {
		return
	}

	repositories = reposResponse.Repositories
	return
}
