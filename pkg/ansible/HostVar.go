package ansible

import "github.com/aallbrig/allbctl/pkg"

type HostVar struct {
	Name string
	Data interface{}
}

var DefaultHostVarFilename = "localhost.yaml"

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
