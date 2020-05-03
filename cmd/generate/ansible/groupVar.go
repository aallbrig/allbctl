package ansible

import (
	"github.com/aallbrig/allbctl/pkg"

	"github.com/spf13/cobra"
)

type GroupKeyValue struct {
	Key   string
	Value string
}

type GroupValues struct {
	Values []GroupKeyValue
}

var DefaultGroupValues = GroupValues{
	Values: []GroupKeyValue{
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

func GenerateGroupVar(filename string) {
	if filename == "" {
		filename = "group_var.yaml"
	}

	pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/ansible/group_var.yaml.tmpl",
			Defaults: DefaultGroupValues,
		},
		&pkg.ResultingFile{
			Filename:    filename,
			RelativeDir: "ansible/inventory/group_vars",
		},
	)
}

var groupVarCmd = &cobra.Command{
	Use:   "groupVar",
	Short: "code generation for ansible group var",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateGroupVar("")
	},
}

func init() {
	Cmd.AddCommand(groupVarCmd)
}
