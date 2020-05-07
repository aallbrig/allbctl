package dockerfile

import (
	"github.com/aallbrig/allbctl/pkg"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "code generation for Dockerfiles",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.HelpTextIfEmpty(cmd, args)
	},
}
