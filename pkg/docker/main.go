package docker

type Dockerfile struct {
	Name    string
	Image   string
	Version string
}

var Dockerfiles = []Dockerfile{
	{
		Name:    "Ansible",
		Image:   "ansible/ansible",
		Version: "ubuntu1604",
	},
	{
		Name:    "Alpine",
		Image:   "alpine",
		Version: "latest",
	},
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

