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

var roleName string

func roleNamePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Role name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("empty input -- please provide role name for Ansible role")
			}
			return nil
		},
		Default: ansible.DefaultRoleName,
	}

	result, err := prompt.Run()
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
				roleName, err = roleNamePrompt()
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
