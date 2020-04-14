package ansible

import (
	"fmt"

	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "code generation for ansible role",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ansible role generator called")
	},
}

func init() {
	Cmd.AddCommand(roleCmd)
}
