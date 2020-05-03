package ansible

import (
	"github.com/aallbrig/allbctl/cmd/generate/dockerfile"
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
		dockerfile.GenerateAnsibleDockerfile("")
	},
}

func init() {
	Cmd.AddCommand(initCmd)
}
