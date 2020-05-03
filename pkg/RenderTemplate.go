package pkg

import (
	"bytes"
	"html/template"
	"io"
	"log"

	"github.com/markbates/pkger"
)

type TemplateFile struct {
	Path     string
	Defaults interface{}
}

type ResultingFile struct {
	Filename    string
	RelativeDir string
}

func RenderTemplateByFile(tf *TemplateFile, rf *ResultingFile) {
	templateFile, err := pkger.Open(tf.Path)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, templateFile)
	if err != nil {
		log.Fatalf("Unable to copy template file into buffer: %v", err)
	}

	tmpl, err := template.New(rf.Filename).Parse(buf.String())
	if err != nil {
		log.Fatalf("Unable to load template: %v", err)
	}

	fileContents := new(bytes.Buffer)
	err = tmpl.Execute(fileContents, tf.Defaults)
	if err != nil {
		log.Fatalf("Unable to render template: %v", err)
	}

	FilesToGenerate = append(FilesToGenerate, GenerateFile{
		RelativeDir:  rf.RelativeDir,
		FileName:     rf.Filename,
		FileContents: fileContents,
	})
}
