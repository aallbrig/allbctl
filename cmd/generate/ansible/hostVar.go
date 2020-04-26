package ansible

import (
	"bytes"
	"io"
	"text/template"

	"github.com/aallbrig/allbctl/pkg"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
)

type HostKeyValue struct {
	Key   string
	Value string
}

type HostValues struct {
	Values []HostKeyValue
}

var DefaultHostValues = HostValues{
	Values: []HostKeyValue{
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
}

func GenerateHostVar(filename string) {
	if filename == "" {
		filename = "host_var.yaml"
	}

	templateFile, _ := pkger.Open("/templates/ansible/host_var.yaml.tmpl")

	buf := new(bytes.Buffer)
	io.Copy(buf, templateFile)
	tmpl, _ := template.New("hostVar").Parse(buf.String())
	fileContents := new(bytes.Buffer)
	_ = tmpl.Execute(fileContents, DefaultHostValues)

	pkg.FilesToGenerate = append(pkg.FilesToGenerate, pkg.GenerateFile{
		RelativeDir:  "ansible/inventory/host_vars",
		FileName:     filename,
		FileContents: fileContents,
	})
}

var hostVarCmd = &cobra.Command{
	Use:   "hostVar",
	Short: "code generation for ansible host var",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateHostVar("")
	},
}

func init() {
	Cmd.AddCommand(hostVarCmd)
}
