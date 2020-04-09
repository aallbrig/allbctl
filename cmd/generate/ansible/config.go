package ansible

import (
	"fmt"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "code generation for ansible config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ansible config generator called")
	},
}

func init() {
	Cmd.AddCommand(configCmd)
}
