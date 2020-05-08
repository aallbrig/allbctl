package ansible

import (
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"path/filepath"
)

var DefaultRoleName = "DefaultRoleName"

type Role struct {
	Name string
}

func (role *Role) RenderFiles(data interface{}) error {
	// TODO: opportunity to further decompose this function?
	err := pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/key_value_dict.yaml.tmpl",
			Data: data,
		},
		&pkg.ResultingFile{
			Filename:    "main.yaml",
			RelativeDir: filepath.Join("ansible/roles", role.Name, "/vars"),
		},
	)
	if err != nil {
		fmt.Println("Error creating vars main file")
		return err
	}

	err = pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/key_value_dict.yaml.tmpl",
			Data: data,
		},
		&pkg.ResultingFile{
			Filename:    "main.yaml",
			RelativeDir: filepath.Join("ansible/roles", role.Name, "/Data"),
		},
	)
	if err != nil {
		fmt.Println("Error creating Data main file")
		return err
	}

	err = pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/key_value_dict.yaml.tmpl",
			Data: data,
		},
		&pkg.ResultingFile{
			Filename:    "main.yaml",
			RelativeDir: filepath.Join("ansible/roles", role.Name, "/tasks"),
		},
	)
	if err != nil {
		fmt.Println("Error creating tasks main file")
		return err
	}

	err = pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/key_value_dict.yaml.tmpl",
			Data: data,
		},
		&pkg.ResultingFile{
			Filename:    "main.yaml",
			RelativeDir: filepath.Join("ansible/roles", role.Name, "/handlers"),
		},
	)
	if err != nil {
		fmt.Println("Error creating role main file")
		return err
	}

	return nil

}
