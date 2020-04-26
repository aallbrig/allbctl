package ansible

import (
	"bytes"
	"io"
	"text/template"

	"github.com/aallbrig/allbctl/pkg"

	"github.com/markbates/pkger"

	"github.com/spf13/cobra"
)

type GroupKeyValue struct {
	Key   string
	Value string
}

type GroupValues struct {
	Values []GroupKeyValue
}

var DefaultGroupValues = GroupValues{
	Values: []GroupKeyValue{
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

func GenerateGroupVar(filename string) {
	if filename == "" {
		filename = "host_var.yaml"
	}

	templateFile, _ := pkger.Open("/templates/ansible/group_var.yaml.tmpl")

	buf := new(bytes.Buffer)
	io.Copy(buf, templateFile)
	tmpl, _ := template.New("groupVar").Parse(buf.String())
	fileContents := new(bytes.Buffer)
	_ = tmpl.Execute(fileContents, DefaultGroupValues)

	pkg.FilesToGenerate = append(pkg.FilesToGenerate, pkg.GenerateFile{
		RelativeDir:  "ansible/inventory/group_vars",
		FileName:     "group_var.yaml",
		FileContents: fileContents,
	})
}

var groupVarCmd = &cobra.Command{
	Use:   "groupVar",
	Short: "code generation for ansible group var",
	Run: func(cmd *cobra.Command, args []string) {
		GenerateGroupVar("")
	},
}

func init() {
	Cmd.AddCommand(groupVarCmd)
}
