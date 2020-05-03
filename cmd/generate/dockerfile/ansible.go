package dockerfile

import (
	"github.com/aallbrig/allbctl/pkg"

	"github.com/spf13/cobra"
)

var defaults = Dockerfile{
	Image:   "ansible/ansible",
	Version: "ubuntu1604",
}

func GenerateAnsibleDockerfile(filename string) {
	if filename == "" {
		filename = "ansible.Dockerfile"
	}

	pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/docker/Dockerfile.tmpl",
			Defaults: defaults,
		},
		&pkg.ResultingFile{
			Filename:    filename,
			RelativeDir: "dockerfiles",
		},
	)
}

func init() {
	Cmd.AddCommand(
		&cobra.Command{
			Use:   "ansible",
			Short: "Generates a dockerfile appropriate for running Ansible playbooks",
			Run: func(cmd *cobra.Command, args []string) {
				GenerateAnsibleDockerfile("")
			},
		},
	)
}
