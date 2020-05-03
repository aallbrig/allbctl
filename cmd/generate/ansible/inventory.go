package ansible

import (
	"github.com/aallbrig/allbctl/pkg"

	"github.com/spf13/cobra"
)

type InventoryKeyValue struct{}

type InventoryValues struct {
	Values []InventoryKeyValue
}

var DefaultInventoryValues = InventoryValues{
	Values: []InventoryKeyValue{},
}

func GenerateInventory() {
	pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/ansible/inventory.yaml.tmpl",
			Defaults: DefaultInventoryValues,
		},
		&pkg.ResultingFile{
			Filename:    "hosts.yaml",
			RelativeDir: "ansible/inventory",
		},
	)
}

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateInventory()
	},
}

func init() {
	Cmd.AddCommand(inventoryCmd)
}
