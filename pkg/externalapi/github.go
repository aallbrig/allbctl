package externalapi

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

// GithubAuthTokenEnvVar https://cli.github.com/manual/gh_help_environment
// default environment variables: GH_TOKEN, GITHUB_TOKEN
var GithubAuthTokenEnvVar = "GITHUB_TOKEN"
var tokenProvider IGithubTokenProvider = &GithubAuthTokenProvider{}
var ghProvider IGithubClientProvider = &GithubClientProvider{}

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type IGithubTokenProvider interface {
	GetAuthToken() (string, error)
}

// GithubAuthTokenProvider providers a way to get the GH auth token
type GithubAuthTokenProvider struct{}

// GetAuthToken gets github auth token
func (provider *GithubAuthTokenProvider) GetAuthToken() (authToken string, err error) {
	authToken = os.Getenv(GithubAuthTokenEnvVar)

	if authToken == "" {
		err = fmt.Errorf("%s envvar not set", GithubAuthTokenEnvVar)
	}

	return
}

type githubUsersService interface {
	Get(ctx context.Context, userName string) (user *github.User, response *github.Response, err error)
}

type githubSearchService interface {
	Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult,
		*github.Response, error)
}

// GithubClient is this program's facade of github API client
type GithubClient struct {
	Users  githubUsersService
	Search githubSearchService
}

type IGithubClientProvider interface {
	GetGithubClient(accessToken string) (client GithubClient, err error)
}

// GithubClientProvider allows consumer to get a GithubClient
type GithubClientProvider struct{}

// GetGithubClient implements providing GithubClient
func (provider *GithubClientProvider) GetGithubClient(accessToken string) (client GithubClient, err error) {
	// TODO: Check access token against some rules? (e.g. "not empty")
	tokenSource := &tokenSource{
		AccessToken: accessToken,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	ghClient := github.NewClient(oauthClient)

	client = GithubClient{
		Users:  ghClient.Users,
		Search: ghClient.Search,
	}

	return
}

// GetMyDotfiles gets dotfile repos from my github
func GetMyDotfiles() (repositories []github.Repository, err error) {
	ctx := context.TODO()

	githubAuthToken, err := tokenProvider.GetAuthToken()
	if err != nil {
		return
	}

	ghClient, err := ghProvider.GetGithubClient(githubAuthToken)
	if err != nil {
		return
	}

	user, _, err := ghClient.Users.Get(ctx, "")
	if err != nil {
		return
	}

	searchQuery := fmt.Sprintf("dotfiles user:%s", *user.Login)
	reposResponse, _, err := ghClient.Search.Repositories(ctx, searchQuery, nil)
	if err != nil {
		return
	}

	repositories = reposResponse.Repositories
	return
}
