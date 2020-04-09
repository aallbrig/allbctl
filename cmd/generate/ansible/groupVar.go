package ansible

import (
	"fmt"
	"github.com/spf13/cobra"
)

var groupVarCmd = &cobra.Command{
	Use:   "groupVar",
	Short: "code generation for ansible group var",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ansible group var generator called")
	},
}

func init() {
	Cmd.AddCommand(groupVarCmd)
}
