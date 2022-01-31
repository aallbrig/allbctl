package macbook

import (
	"bytes"
	"os/exec"
)

var dockOnLeftDefaultsCmd = DefaultsCommand{
	Domain:        "com.apple.dock",
	Key:           "orientation",
	ExpectedValue: "left",
	ValueType:     DefaultsString,
}

type DockOnLeft struct{}

func NewDockOnLeft() *DockOnLeft {
	return &DockOnLeft{}
}

func (d DockOnLeft) Name() string {
	return "Dock on left"
}

func (d DockOnLeft) Validate() (err error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")
	err, validateOut := dockOnLeftDefaultsCmd.Validate()
	out.WriteString(validateOut.String())
	return
}

func (d DockOnLeft) Install() (err error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")
	err, installOut := dockOnLeftDefaultsCmd.WriteExpectedValue()
	out.WriteString(installOut.String())

	cmd := exec.Command("killall", "Dock")
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return
}

func (d DockOnLeft) Uninstall() (err error, out *bytes.Buffer) {
	out = bytes.NewBufferString("")
	err, deleteOut := dockOnLeftDefaultsCmd.Delete()
	out.WriteString(deleteOut.String())

	cmd := exec.Command("killall", "Dock")
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return
}
