package ansible

import (
	"errors"
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/aallbrig/allbctl/pkg/ansible"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
)

var hostVarName string

func hostVarNamePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Host var file name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("empty input -- please provide file name for Ansible host var file")
			}
			return nil
		},
		Default: ansible.DefaultHostVarFilename,
	}

	result, err := prompt.Run()
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
				hostVarName, err = hostVarNamePrompt()
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
