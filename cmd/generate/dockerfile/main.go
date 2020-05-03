package dockerfile

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "code generation for Dockerfiles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dockerfile root called")
	},
}
