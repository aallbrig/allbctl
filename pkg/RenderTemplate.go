package pkg

import (
	"bytes"
	"fmt"
	"github.com/markbates/pkger"
	"html/template"
	"io"
)

type GenerateFile struct {
	FileName     string
	FileContents *bytes.Buffer
	RelativeDir  string
}

type TemplateFile struct {
	Path string
	Data interface{}
}

type ResultingFile struct {
	Filename    string
	RelativeDir string
}

var FilesToGenerate []GenerateFile

func AddFileToRenderQueue(f GenerateFile) {
	FilesToGenerate = append(FilesToGenerate, f)
}

func RenderTemplateByFile(tf *TemplateFile, rf *ResultingFile) error {
	templateFile, err := pkger.Open(tf.Path)
	if err != nil {
		fmt.Printf("Unable to open file: %v", err)
		return err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, templateFile)
	if err != nil {
		fmt.Printf("Unable to copy template file into buffer: %v", err)
		return err
	}

	tmpl, err := template.New(rf.Filename).Parse(buf.String())
	if err != nil {
		fmt.Printf("Unable to load template: %v", err)
		return err
	}

	fileContents := new(bytes.Buffer)
	err = tmpl.Execute(fileContents, tf.Data)
	if err != nil {
		fmt.Printf("Unable to render template: %v", err)
		return err
	}

	AddFileToRenderQueue(
		GenerateFile{
			RelativeDir:  rf.RelativeDir,
			FileName:     rf.Filename,
			FileContents: fileContents,
		},
	)

	return nil
}
