package externalcmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/google/go-github/github"
	"os"
	"path"
)

var gitProvider gitClientProvider = &GitClientProvider{}
var Auth transport.AuthMethod

// GitClient is a facade for git operations
type GitClient struct{}

// PlainClone is a facade for git plain clone
func (gitClient *GitClient) PlainClone(dir string, url string) (repo *git.Repository, err error) {
	repo, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:  url,
		Auth: Auth,
	})

	return
}

type gitClientProvider interface {
	GetGitClient() (GitClient, error)
}

// GitClientProvider provides git client facade
type GitClientProvider struct{}

// GetGitClient is the implementation that actually returns the git client
func (provider *GitClientProvider) GetGitClient() (client GitClient, err error) {
	client = GitClient{}
	return
}

// CloneGithubRepo is the function that actually clones a github repo
func CloneGithubRepo(targetDir string, repository *github.Repository) (localRepo *git.Repository, err error) {
	client, err := gitProvider.GetGitClient()
	if err != nil {
		return
	}

	repositoryDir := path.Join(targetDir, *repository.Name)
	if _, statErr := os.Stat(repositoryDir); os.IsNotExist(statErr) {
		localRepo, err = client.PlainClone(repositoryDir, *repository.CloneURL)
	}

	return
}
