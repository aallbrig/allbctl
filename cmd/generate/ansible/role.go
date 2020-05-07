package ansible

import (
	"github.com/aallbrig/allbctl/pkg"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "code generation for ansible role",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.HelpTextIfEmpty(cmd, args)
	},
}

func init() {
	Cmd.AddCommand(roleCmd)
}
