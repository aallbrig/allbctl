package dotfiles

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aallbrig/allbctl/pkg/computersetup/os_agnostic"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/externalcmd"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"os"
	"path"
)

type DotFilerGremlin struct {
	GithubAuthTokenProvider externalapi.IGithubTokenProvider
	GithubClientProvider    externalapi.IGithubClientProvider
}

func NewDotfilesGremlin() *DotFilerGremlin {
	return &DotFilerGremlin{
		GithubAuthTokenProvider: &externalapi.GithubAuthTokenProvider{},
		GithubClientProvider:    &externalapi.GithubClientProvider{},
	}
}

func (d DotFilerGremlin) Name() string {
	return "Dotfiles Gremlin"
}

func (d DotFilerGremlin) Validate() (err error, out *bytes.Buffer) {
	osAgonstic := os_agnostic.OperatingSystem{}
	dotfiles, err := d.GetMyDotfiles()
	_, homeDir := osAgonstic.HomeDir()
	for _, repo := range dotfiles {
		repoLocation := path.Join(homeDir, "src", *repo.Name)
		if _, statErr := os.Stat(repoLocation); os.IsNotExist(statErr) {
			err = statErr
		}
	}
	return
}

func (d DotFilerGremlin) Install() (err error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")
	osAgnostic := os_agnostic.OperatingSystem{}
	dotfiles, err := d.GetMyDotfiles()
	if err != nil {
		out.WriteString(fmt.Sprintf("❌ Error getting dotfiles %v\n", err))
		return
	}

	token, err := d.GithubAuthTokenProvider.GetAuthToken()
	if err != nil {
		out.WriteString(fmt.Sprintf("❌ Error getting github auth token %v", err))
		return
	}
	ghClient, err := d.GithubClientProvider.GetGithubClient(token)
	if err != nil {
		out.WriteString(fmt.Sprintf("❌ Error getting github client %v", err))
		return
	}
	user, _, err := ghClient.Users.Get(context.TODO(), "")
	if err != nil {
		out.WriteString(fmt.Sprintf("❌ Error getting github user %v", err))
		return
	}

	externalcmd.Auth = &http.BasicAuth{
		Username: *user.Login,
		Password: token,
	}

	_, homeDir := osAgnostic.HomeDir()
	for _, repo := range dotfiles {
		repoLocation := path.Join(homeDir, "src", *repo.Name)
		if _, statErr := os.Stat(repoLocation); os.IsNotExist(statErr) {
			_, innerErr := externalcmd.CloneGithubRepo(path.Join(homeDir, "src"), &repo)
			if innerErr != nil {
				out.WriteString(fmt.Sprintf("Error cloning repo %v", innerErr))
				err = innerErr
				return
			}
		}
	}
	return
}

func (d DotFilerGremlin) Uninstall() (err error, out *bytes.Buffer) {
	return
}

func (d DotFilerGremlin) GetMyDotfiles() (repositories []github.Repository, err error) {
	ctx := context.TODO()

	githubAuthToken, err := d.GithubAuthTokenProvider.GetAuthToken()
	if err != nil {
		return
	}

	ghClient, err := d.GithubClientProvider.GetGithubClient(githubAuthToken)
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
