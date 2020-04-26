package ansible

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/aallbrig/allbctl/pkg"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
)

func GenerateConfig() {
	templateFile, _ := pkger.Open("/templates/ansible/ansible.cfg.tmpl")

	buf := new(bytes.Buffer)
	io.Copy(buf, templateFile)
	tmpl, _ := template.New("config").Parse(buf.String())
	fileContents := new(bytes.Buffer)
	_ = tmpl.Execute(fileContents, nil)

	pkg.FilesToGenerate = append(pkg.FilesToGenerate, pkg.GenerateFile{
		FileName:     "ansible.cfg",
		FileContents: fileContents,
	})
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "code generation for ansible config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ansible config generator called")
	},
}

func init() {
	Cmd.AddCommand(configCmd)
}
