package pkg

import "bytes"

type GenerateFile struct {
	RelativeDir  string
	FileName     string
	FileContents *bytes.Buffer
}

var FilesToGenerate []GenerateFile
