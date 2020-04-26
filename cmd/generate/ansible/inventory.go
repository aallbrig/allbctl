package ansible

import (
	"bytes"
	"io"
	"text/template"

	"github.com/aallbrig/allbctl/pkg"

	"github.com/markbates/pkger"

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
	templateFile, _ := pkger.Open("/templates/ansible/inventory.yaml.tmpl")

	buf := new(bytes.Buffer)
	io.Copy(buf, templateFile)
	tmpl, _ := template.New("inventory").Parse(buf.String())
	fileContents := new(bytes.Buffer)
	_ = tmpl.Execute(fileContents, DefaultInventoryValues)

	pkg.FilesToGenerate = append(pkg.FilesToGenerate, pkg.GenerateFile{
		RelativeDir:  "ansible/inventory",
		FileName:     "hosts.yaml",
		FileContents: fileContents,
	})
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
