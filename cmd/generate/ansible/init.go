package ansible

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateInventory()
		GenerateConfig()
		GenerateGroupVar("")
		GenerateHostVar("")
	},
}

func init() {
	Cmd.AddCommand(initCmd)
}
