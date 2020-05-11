package ansible

import "github.com/aallbrig/allbctl/pkg"

type Inventory struct{}

func (c *Inventory) RenderFiles() error {
	err := pkg.RenderTemplateByFile(
		&pkg.TemplateFile{
			Path: "/templates/ansible/ansible.cfg.tmpl",
			Data: KeyValuePairs{
				Values: []KeyValue{
					{
						Key:   "key1",
						Value: "Value1",
					},
					{
						Key:   "key2",
						Value: "Value2",
					},
					{
						Key:   "key3",
						Value: "Value3",
					},
				},
			},
		},
		&pkg.ResultingFile{
			Filename:    "hosts.yaml",
			RelativeDir: "ansible/inventory",
		},
	)
	if err != nil {
		return err
	}
	return nil
}
