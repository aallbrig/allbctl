package model

import "bytes"

type IMachineConfiguration interface {
	Name() string
	Validate() (out *bytes.Buffer, err error)
	Install() (out *bytes.Buffer, err error)
	Uninstall() (out *bytes.Buffer, err error)
}
