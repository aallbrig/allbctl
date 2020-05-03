package ansible

import (
	"github.com/aallbrig/allbctl/pkg"
	"github.com/spf13/cobra"
)

type HostKeyValue struct {
	Key   string
	Value string
}

type HostValues struct {
	Values []HostKeyValue
}

var DefaultHostValues = HostValues{
	Values: []HostKeyValue{
		{
			Key:   "key1",
			Value: "Value1",
		},
		{
			Key:   "key2",
			Value: "Value2",
		},
		{
			Key:   "key3",
			Value: "Value3",
		},
	},
}

func GenerateHostVar(filename string) {
	if filename == "" {
		filename = "host_var.yaml"
	}

	pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/ansible/host_var.yaml.tmpl",
			Defaults: DefaultHostValues,
		},
		&pkg.ResultingFile{
			Filename:    filename,
			RelativeDir: "ansible/inventory/host_vars",
		},
	)
}

var hostVarCmd = &cobra.Command{
	Use:   "hostVar",
	Short: "code generation for ansible host var",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateHostVar("")
	},
}

func init() {
	Cmd.AddCommand(hostVarCmd)
}
