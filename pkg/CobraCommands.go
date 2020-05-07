package pkg

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func HelpTextIfEmpty(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cmd.Help()
		if err != nil {
			log.Fatalf("Error generating help text: %v", err)
		}
		os.Exit(0)
	}
}
