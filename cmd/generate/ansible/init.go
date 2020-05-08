package ansible

import (
	"github.com/aallbrig/allbctl/cmd/generate/dockerfile"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		inventoryCmd.Run(cmd, args)
		configCmd.Run(cmd, args)
		groupVarCmd.Run(cmd, args)
		hostVarCmd.Run(cmd, args)
		dockerfile.GenerateAnsibleDockerfile("")
	},
}

func init() {
	Cmd.AddCommand(initCmd)
}
