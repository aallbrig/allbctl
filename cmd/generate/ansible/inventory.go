package ansible

import (
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

type InventoryKeyValue struct{}

type InventoryValues struct {
	Values []InventoryKeyValue
}

var DefaultInventoryValues = InventoryValues{
	Values: []InventoryKeyValue{},
}

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "initialize ansible project",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := template.Must(template.ParseFiles("./templates/ansible/host_var.yaml.tmpl"))
		_ = tmpl.Execute(os.Stdout, DefaultInventoryValues)
	},
}

func init() {
	Cmd.AddCommand(inventoryCmd)
}
