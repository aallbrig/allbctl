package docker

import (
	"github.com/aallbrig/allbctl/pkg"
	"log"
	"os"
	"path/filepath"
)

type Dockerfile struct {
	Name          string
	Image         string
	Version       string
	ResultingFile *pkg.ResultingFile
}

var genericDockerfile = &pkg.ResultingFile{
	Filename:    "Dockerfile",
	RelativeDir: ".",
}

var Dockerfiles = []Dockerfile{
	{
		Name:          "Ansible",
		Image:         "ansible/ansible",
		Version:       "ubuntu1604",
		ResultingFile: genericDockerfile,
	},
	{
		Name:          "Alpine",
		Image:         "alpine",
		Version:       "latest",
		ResultingFile: genericDockerfile,
	},
}

func (df *Dockerfile) Filepath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory:\n%v", err)
	}
	return filepath.Join(cwd, df.ResultingFile.RelativeDir, df.ResultingFile.Filename)
}

func ListDockerFileNames(dfs []Dockerfile) []string {
	var names []string
	for i := range dfs {
		names = append(names, dfs[i].Name)
	}
	return names
}

func NameIsInDockerfiles(dfs []Dockerfile, name string) bool {
	found := false
	for i := range dfs {
		if dfs[i].Name == name {
			found = true
		}
	}
	return found
}

func GetDockerfileByName(dfs []Dockerfile, name string) Dockerfile {
	var dockerfile Dockerfile
	for i := range dfs {
		if dfs[i].Name == name {
			dockerfile = dfs[i]
		}
	}
	return dockerfile
}
