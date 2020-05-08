package ansible

import (
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/aallbrig/allbctl/pkg/ansible"
	"github.com/spf13/cobra"
	"log"
)

var hostVarName string

func promptForHostVarName() (string, error) {
	result, err := ansible.HostVarNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}

var hostVarCmd = &cobra.Command{
	Use:   "hostVar",
	Short: "code generation for ansible host var",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if hostVarName == "" {
			if pkg.Interactive {
				hostVarName, err = promptForHostVarName()
				if err != nil {
					log.Fatalf("Error acquiring host var name: %v\n", err)
				}
			} else {
				hostVarName = ansible.DefaultHostVarFilename
			}
		}

		hostVar := ansible.HostVar{
			Name: hostVarName,
			Data: ansible.KeyValuePairs{
				Values: []ansible.KeyValue{
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
		}
		err = hostVar.RenderFiles()
		if err != nil {
			log.Fatalf("Error rendering host var file(s): %v", err)
		}
	},
}

func init() {
	hostVarCmd.Flags().StringVarP(&hostVarName, "name", "n", "", "Name of Ansible host var")

	Cmd.AddCommand(hostVarCmd)
}
