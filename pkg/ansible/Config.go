package ansible

import "github.com/aallbrig/allbctl/pkg"

type Config struct {}

func (c *Config) RenderFiles () error {
	err := pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/ansible.cfg.tmpl",
			Data: nil,
		},
		&pkg.ResultingFile{
			Filename:    "ansible.cfg",
			RelativeDir: "",
		},
	)
	if err != nil {
		return err
	}
	return nil
}
