package cmd

import (
	"context"
	"fmt"
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/computersetup/os_agnostic"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/externalcmd"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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
		os := os_agnostic.OperatingSystem{}
		identifier := computerSetup.MachineIdentifier{}
		err, name := os.GetName()
		if err != nil {
			log.Fatalf("Issues getting operating system identifier")
		}

		configProvider := identifier.ConfigurationProviderForOperatingSystem(name)
		if configProvider == nil {
			log.Fatal(fmt.Sprintf("No configuration provider found for operationg system %s", os))
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, out := tweaker.ApplyConfiguration()
		log.Print(out)

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

		_, homeDir := os.HomeDir()
		for _, repo := range dotfiles {
			_, innerErr := externalcmd.CloneGithubRepo(path.Join(homeDir, "src"), &repo)
			if innerErr != nil {
				log.Fatalf("Error cloning repo %v", innerErr)
			}
		}
	},
}
