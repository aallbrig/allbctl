package cmd

import (
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
)

// RootCmd used in allbctl root
var RootCmd = &cobra.Command{
	Use:   "computersetup",
	Short: "Setup computer in expected way",
	Run: func(cmd *cobra.Command, args []string) {
		usrHomeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal("Error getting user home directory")
		}
		err = computerSetup.DirectoryForSourceCode(usrHomeDir)
		if err != nil {
			log.Fatal("Error creating src directory")
		}

		log.Println("src directory exists in home directory")
	},
}
