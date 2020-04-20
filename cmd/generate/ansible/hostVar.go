package ansible

import (
	"os"
	"text/template"

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

var hostVarCmd = &cobra.Command{
	Use:   "hostVar",
	Short: "code generation for ansible host var",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := template.Must(template.ParseFiles("./templates/ansible/host_var.yaml.tmpl"))
		_ = tmpl.Execute(os.Stdout, DefaultHostValues)
	},
}

func init() {
	Cmd.AddCommand(hostVarCmd)
}
