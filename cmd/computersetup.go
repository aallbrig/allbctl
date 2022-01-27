package cmd

import (
	"context"
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/externalcmd"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
	"path"
)

// ComputerSetupCmd defines the root of computer setup
var ComputerSetupCmd = &cobra.Command{
	Use: "computer-setup",
	Aliases: []string{
		"computersetup",
		"cs",
		"setup",
	},
	Short: "Configure host to developer preferences (cross platform)",
	Run: func(cmd *cobra.Command, args []string) {
		usrHomeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal("Error getting user home directory")
		}

		err = computerSetup.DirectoryForSourceCode(usrHomeDir)
		if err != nil {
			log.Fatal("Error creating src directory")
		}
		log.Println("$HOME/src is available")

		err = computerSetup.DirectoryForUserBin(usrHomeDir)
		if err != nil {
			log.Fatal("Error creating src directory")
		}
		log.Println("$HOME/bin is available")

		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationForMachine()
		if configProvider != nil {
			tweaker := computerSetup.MachineTweaker{
				MachineConfiguration: configProvider.GetConfiguration(),
			}
			tweaker.ApplyConfiguration()
		}

		tokenProvider := externalapi.GithubAuthTokenProvider{}
		githubClientProvider := externalapi.GithubClientProvider{}
		dotfiles, err := externalapi.GetMyDotfiles()
		if err != nil {
			log.Fatalf("Error getting dotfiles %v", err)
		}

		token, err := tokenProvider.GetAuthToken()
		if err != nil {
			log.Fatalf("Error getting github auth token %v", err)
		}
		ghClient, err := githubClientProvider.GetGithubClient(token)
		if err != nil {
			log.Fatalf("Error getting github client %v", err)
		}
		user, _, err := ghClient.Users.Get(context.TODO(), "")
		if err != nil {
			log.Fatalf("Error getting github user %v", err)
		}

		externalcmd.Auth = &http.BasicAuth{
			Username: *user.Login,
			Password: token,
		}

		for _, repo := range dotfiles {
			_, innerErr := externalcmd.CloneGithubRepo(path.Join(usrHomeDir, "src"), &repo)
			if innerErr != nil {
				log.Fatalf("Error cloning repo %v", innerErr)
			}
		}
	},
}
