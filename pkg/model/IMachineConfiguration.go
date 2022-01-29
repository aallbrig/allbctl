package model

import "bytes"

type IMachineConfiguration interface {
	Name() string
	Validate() (err error, out *bytes.Buffer)
	Install() (err error, out *bytes.Buffer)
	Uninstall() (err error, out *bytes.Buffer)
}
