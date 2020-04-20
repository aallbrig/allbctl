package ansible

import (
	"os"
	"text/template"

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

var groupVarCmd = &cobra.Command{
	Use:   "groupVar",
	Short: "code generation for ansible group var",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := template.Must(template.ParseFiles("./templates/ansible/group_var.yaml.tmpl"))
		_ = tmpl.Execute(os.Stdout, DefaultGroupValues)
	},
}

func init() {
	Cmd.AddCommand(groupVarCmd)
}
