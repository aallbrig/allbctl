package ansible

import (
	"errors"
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/aallbrig/allbctl/pkg/ansible"
	"github.com/manifoldco/promptui"
	"log"

	"github.com/spf13/cobra"
)

var groupVarName string

func groupVarNamePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Host var file name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("empty input -- please provide file name for Ansible group var file")
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

var groupVarCmd = &cobra.Command{
	Use:   "groupVar",
	Short: "code generation for ansible group var",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if groupVarName == "" {
			if pkg.Interactive {
				groupVarName, err = groupVarNamePrompt()
				if err != nil {
					log.Fatalf("Error acquiring host var name: %v\n", err)
				}
			} else {
				groupVarName = ansible.DefaultGroupVarFilename
			}
		}

		groupVar := ansible.GroupVar{
			Name: groupVarName,
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
		err = groupVar.RenderFiles()
		if err != nil {
			log.Fatalf("Error rendering group var file(s): %v", err)
		}
	},
}

func init() {
	groupVarCmd.Flags().StringVarP(&groupVarName, "name", "n", "", "Name of Ansible group var")

	Cmd.AddCommand(groupVarCmd)
}
