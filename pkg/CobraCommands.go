package pkg

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var Interactive bool
var WriteStdOut bool

func HelpText(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		log.Fatalf("Error generating help text: %v", err)
	}
	os.Exit(0)
}
func HelpTextIfEmpty(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		HelpText(cmd, args)
	}
}
