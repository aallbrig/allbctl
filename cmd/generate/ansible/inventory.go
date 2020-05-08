package ansible

import (
	"github.com/aallbrig/allbctl/pkg/ansible"
	"github.com/spf13/cobra"
	"log"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		inventory := ansible.Inventory{}
		err := inventory.RenderFiles()
		if err != nil {
			log.Fatalf("Error rendering ansible inventory file: %v\n", err)
		}
	},
}

func init() {
	Cmd.AddCommand(inventoryCmd)
}
