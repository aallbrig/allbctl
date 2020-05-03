package dockerfile

import (
	"github.com/aallbrig/allbctl/pkg"

	"github.com/spf13/cobra"
)

func GenerateAlpineDockerfile(filename string) {
	if filename == "" {
		filename = "alpine.Dockerfile"
	}

	pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/docker/Dockerfile.tmpl",
			Defaults: Dockerfile{
				Image:   "alpine",
				Version: "latest",
			},
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
			Use:   "alpine",
			Short: "Generates a minimal alpine dockerfile",
			Run: func(cmd *cobra.Command, args []string) {
				GenerateAlpineDockerfile("")
			},
		},
	)
}
