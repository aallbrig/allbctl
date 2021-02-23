package cmd

import (
	"bytes"
	"fmt"
	"github.com/aallbrig/allbctl/pkg/status"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
)

// StatusCmd represents status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the status of the machine for expected setup",
	Run: func(cmd *cobra.Command, args []string) {
		usrHomeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal("Error getting user home directory")
		}

		output := bytes.NewBufferString("")
		directoriesToCheck := []string{"src", "bin"}

		output.WriteString("Directory Expectations\n")
		output.WriteString("-----\n")
		for _, dir := range directoriesToCheck {
			_ = status.CheckForDirectory(output, usrHomeDir, dir)
		}

		fmt.Print(output)
	},
}
