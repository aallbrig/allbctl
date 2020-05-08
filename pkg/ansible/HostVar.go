package ansible

import (
	"errors"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/manifoldco/promptui"
)

var DefaultHostVarFilename = "localhost.yaml"

var HostVarNamePrompt = promptui.Prompt{
	Label:    "Host var file name",
	Validate: func(input string) error {
		if input == "" {
			return errors.New("empty input -- please provide file name for Ansible host var file")
		}
		return nil
	},
	Default: DefaultHostVarFilename,
}

type HostVar struct {
	Name string
	Data interface{}
}

func (h *HostVar) RenderFiles() error {
	err := pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/key_value_dict.yaml.tmpl",
			Data: h.Data,
		},
		&pkg.ResultingFile{
			Filename:    h.Name,
			RelativeDir: "ansible/inventory/host_vars",
		},
	)
	if err != nil {
		return err
	}
	return nil
}
