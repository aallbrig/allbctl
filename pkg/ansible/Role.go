package ansible

import (
	"errors"
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/manifoldco/promptui"
	"path/filepath"
)

var DefaultRoleName = "DefaultRoleName"

type Role struct {
	Name string
}

func (role *Role) RenderFiles(defaults interface{}) error {
	// TODO: opportunity to further decompose this function?
	err := pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
			Defaults: defaults,
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
			Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
			Defaults: defaults,
		},
		&pkg.ResultingFile{
			Filename:    "main.yaml",
			RelativeDir: filepath.Join("ansible/roles", role.Name, "/defaults"),
		},
	)
	if err != nil {
		fmt.Println("Error creating defaults main file")
		return err
	}

	err = pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
			Defaults: defaults,
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
			Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
			Defaults: defaults,
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

func RoleNamePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Role name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("empty input -- please provide role name for Ansible role")
			}
			return nil
		},
		Default: DefaultRoleName,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}
