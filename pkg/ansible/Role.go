package ansible

import (
	"errors"
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/manifoldco/promptui"
	"path/filepath"
)

var DefaultRoleName = "DefaultRoleName"

var RoleNamePrompt = promptui.Prompt{
	Label:    "Role name",
	Validate: func(input string) error {
		if input == "" {
								 return errors.New("empty input -- please provide role name for Ansible role")
								 }
		return nil
	},
	Default: DefaultRoleName,
}

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
