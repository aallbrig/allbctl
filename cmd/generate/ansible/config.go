package ansible

import (
	"github.com/aallbrig/allbctl/pkg"
	"github.com/spf13/cobra"
)

func GenerateConfig() {
	pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/ansible/ansible.cfg.tmpl",
			Defaults: nil,
		},
		&pkg.ResultingFile{
			Filename:    "ansible.cfg",
			RelativeDir: "",
		},
	)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "code generation for ansible config",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateConfig()
	},
}

func init() {
	Cmd.AddCommand(configCmd)
}
