package ansible

import (
	"errors"
	"fmt"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var roleName string
// TODO: Get interactivity from global flag var
var defaultRoleName = "defaultRoleName"

func roleNamePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Role name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("empty input -- please provide role name for Ansible role")
			}
			return nil
		},
		Default:  defaultRoleName,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}
type Template struct {
	Values []struct {
		key string
		value string
	}
}

type KeyValue struct {
	Key   string
	Value string
}

type KeyValuePairs struct {
	Values []KeyValue
}
var DefaultKeyValue = KeyValuePairs{
	Values: []KeyValue{
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
				roleName = defaultRoleName
			}
		}
		pkg.RenderTemplateByFile(
			&pkg.TemplateFile{
				Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
				Defaults: DefaultKeyValue,
			},
			&pkg.ResultingFile{
				Filename:    "main.yaml",
				RelativeDir: filepath.Join("ansible/roles", roleName, "/vars"),
			},
		)
		pkg.RenderTemplateByFile(
			&pkg.TemplateFile{
				Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
				Defaults: DefaultKeyValue,
			},
			&pkg.ResultingFile{
				Filename:    "main.yaml",
				RelativeDir: filepath.Join("ansible/roles", roleName, "/defaults"),
			},
		)
		pkg.RenderTemplateByFile(
			&pkg.TemplateFile{
				Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
				Defaults: DefaultKeyValue,
			},
			&pkg.ResultingFile{
				Filename:    "main.yaml",
				RelativeDir: filepath.Join("ansible/roles", roleName, "/tasks"),
			},
		)
		pkg.RenderTemplateByFile(
			&pkg.TemplateFile{
				Path:     "/templates/ansible/key_value_dict.yaml.tmpl",
				Defaults: DefaultKeyValue,
			},
			&pkg.ResultingFile{
				Filename:    "main.yaml",
				RelativeDir: filepath.Join("ansible/roles", roleName, "/handlers"),
			},
		)
	},
}

func init() {
	roleCmd.Flags().StringVarP(&roleName, "roleName", "n", "", "Name of Ansible role")
	Cmd.AddCommand(roleCmd)
}
