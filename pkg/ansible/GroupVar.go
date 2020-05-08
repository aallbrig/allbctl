package ansible

import "github.com/aallbrig/allbctl/pkg"

type GroupVar struct {
	Name string
	Data interface{}
}

var DefaultGroupVarFilename = "localhost.yaml"

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
