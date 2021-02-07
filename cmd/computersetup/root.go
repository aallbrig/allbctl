package computersetup

import (
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
)

// RootCmd defines the root of computer setup
var RootCmd = &cobra.Command{
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
	},
}
