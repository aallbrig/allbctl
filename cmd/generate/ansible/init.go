package ansible

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ansible init project called")
	},
}

func init() {
	Cmd.AddCommand(initCmd)
}
