package dockerfile

import (
	"github.com/aallbrig/allbctl/pkg"
	"github.com/aallbrig/allbctl/pkg/docker"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
)

var dockerfileName string
var dockerfilePrompt = promptui.Select{
	Label: "Which dockerfile would you like to generate?",
	Items: docker.ListDockerFileNames(docker.Dockerfiles),
}

var Cmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "code generation for Dockerfiles",
	Run: func(cmd *cobra.Command, args []string) {
		var dockerfile docker.Dockerfile
		if docker.NameIsInDockerfiles(docker.Dockerfiles, dockerfileName) {
			dockerfile = docker.GetDockerfileByName(docker.Dockerfiles, dockerfileName)
		} else if dockerfileName == "" {
			if pkg.Interactive {
				i, _, err := dockerfilePrompt.Run()
				if err != nil {
					log.Fatalf("Error acquiring dockerfile name: %v\n", err)
				}
				dockerfile = docker.Dockerfiles[i]
			} else {
				pkg.HelpText(cmd, args)
			}
		}

		pkg.RenderTemplateByFile(
			&pkg.TemplateFile{
				Path: "/templates/docker/Dockerfile.tmpl",
				Data: dockerfile,
			},
			&pkg.ResultingFile{
				Filename:    "Dockerfile",
				RelativeDir: ".",
			},
		)
	},
}

func init() {
	Cmd.Flags().StringVarP(&dockerfileName, "name", "n", "", "Which dockerfile to generate?")
}
