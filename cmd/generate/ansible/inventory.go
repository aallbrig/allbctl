package ansible

import (
	"fmt"
	"github.com/spf13/cobra"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ansible inventory generator called")
	},
}

func init() {
	Cmd.AddCommand(inventoryCmd)
}
