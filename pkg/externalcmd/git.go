package externalcmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/google/go-github/github"
	"os"
	"path"
)

// GitClient is a facade around git operations
type GitClient struct{}

// PlainClone is a facade around git plain clone
func (gitClient *GitClient) PlainClone(dir string, url string, auth transport.AuthMethod) (repo *git.Repository, err error) {

	repo, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:  url,
		Auth: auth,
	})

	return
}

type gitClientProvider interface {
	GetGitClient() (GitClient, error)
}

// GitClientProvider provides my git client facade
type GitClientProvider struct{}

// GetGitClient is the implementation that actually returns the git client
func (provider *GitClientProvider) GetGitClient() (client GitClient, err error) {
	client = GitClient{}
	return
}

// CloneGithubRepo is the function that actually clones a github repo
func CloneGithubRepo(targetDir string, gitProvider gitClientProvider,
	repository *github.Repository, auth transport.AuthMethod) (localrepo *git.Repository, err error) {
	client, err := gitProvider.GetGitClient()
	if err != nil {
		return
	}

	repositoryDir := path.Join(targetDir, *repository.Name)
	if _, statErr := os.Stat(repositoryDir); os.IsNotExist(statErr) {
		localrepo, err = client.PlainClone(repositoryDir, *repository.CloneURL, auth)
	}

	return
}
