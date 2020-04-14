package ansible

import (
	"os"
	"text/template"

	"github.com/aallbrig/allbctl/templates/ansible"
	"github.com/spf13/cobra"
)

var groupVarCmd = &cobra.Command{
	Use:   "groupVar",
	Short: "code generation for ansible group var",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := template.Must(template.New("groupVar").Parse(ansible.GroupVarFile))
		_ = tmpl.Execute(os.Stdout, ansible.DefaultGroupValues)
	},
}

func init() {
	Cmd.AddCommand(groupVarCmd)
}
