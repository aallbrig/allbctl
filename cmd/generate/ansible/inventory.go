package ansible

import (
	"os"
	"text/template"

	"github.com/aallbrig/allbctl/templates/ansible"
	"github.com/spf13/cobra"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := template.Must(template.New("inventory").Parse(ansible.InventoryFile))
		_ = tmpl.Execute(os.Stdout, ansible.DefaultInventoryValues)
	},
}

func init() {
	Cmd.AddCommand(inventoryCmd)
}
