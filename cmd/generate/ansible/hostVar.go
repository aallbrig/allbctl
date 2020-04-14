package ansible

import (
	"os"
	"text/template"

	"github.com/aallbrig/allbctl/templates/ansible"
	"github.com/spf13/cobra"
)

var hostVarCmd = &cobra.Command{
	Use:   "hostVar",
	Short: "code generation for ansible host var",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := template.Must(template.New("hostVar").Parse(ansible.HostVarFile))
		_ = tmpl.Execute(os.Stdout, ansible.DefaultHostValues)
	},
}

func init() {
	Cmd.AddCommand(hostVarCmd)
}
