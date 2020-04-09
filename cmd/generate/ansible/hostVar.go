package ansible

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)
type HostKeyValue struct {
	Key string
	Value string
}
type HostValues struct {
	Values []HostKeyValue
}
var defaultHostValues = HostValues{
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
		dir, _ := os.Getwd()
		templatePath := filepath.Join(
			dir,
			"cmd/generate/ansible/templates",
			"host_variable.yaml.tmpl",
		)
		tmpl := template.Must(template.ParseFiles(templatePath))
		tmpl.Execute(os.Stdout, defaultHostValues)
	},
}

func init() {
	Cmd.AddCommand(hostVarCmd)
}
