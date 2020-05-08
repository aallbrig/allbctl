package ansible

import (
	"errors"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/manifoldco/promptui"
)

var DefaultGroupVarFilename = "localhost.yaml"

var GroupVarNamePrompt = promptui.Prompt{
	Label:    "Group var file name",
	Validate: func(input string) error {
		if input == "" {
			return errors.New("empty input -- please provide file name for Ansible group var file")
		}
		return nil
	},
	Default: DefaultGroupVarFilename,
}

type GroupVar struct {
	Name string
	Data interface{}
}

func (h *GroupVar) RenderFiles() error {
	err := pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/key_value_dict.yaml.tmpl",
			Data: h.Data,
		},
		&pkg.ResultingFile{
			Filename:    h.Name,
			RelativeDir: "ansible/inventory/group_vars",
		},
	)
	if err != nil {
		return err
	}
	return nil
}
