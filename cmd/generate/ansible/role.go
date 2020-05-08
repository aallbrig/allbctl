package ansible

import (
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/aallbrig/allbctl/pkg/ansible"
	"github.com/spf13/cobra"
	"log"
)

var roleName string

func promptForRoleName() (string, error) {
	result, err := ansible.RoleNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "code generation for ansible role",
	Run: func(cmd *cobra.Command, args []string) {
		if roleName == "" {
			if pkg.Interactive {
				var err error
				roleName, err = promptForRoleName()
				if err != nil {
					log.Fatalf("Error acquiring role name: %v\n", err)
				}
			} else {
				roleName = ansible.DefaultRoleName
			}
		}

		var role = &ansible.Role{
			Name: roleName,
		}
		role.RenderFiles(ansible.DefaultKeyValue)
	},
}

func init() {
	roleCmd.Flags().StringVarP(&roleName, "roleName", "n", "", "Name of Ansible role")

	Cmd.AddCommand(roleCmd)
}
