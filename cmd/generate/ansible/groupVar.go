package ansible

import (
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/aallbrig/allbctl/pkg/ansible"
	"log"

	"github.com/spf13/cobra"
)

var groupVarName string

func promptForGroupVarName() (string, error) {
	result, err := ansible.GroupVarNamePrompt.Run()
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
				groupVarName, err = promptForGroupVarName()
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
