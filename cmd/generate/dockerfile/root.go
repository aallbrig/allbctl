package dockerfile

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "code generation for Dockerfiles",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				log.Fatalf("Error generating help text: %v", err)
			}
			os.Exit(0)
		}
	},
}
