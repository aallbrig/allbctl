package ansible

import (
	"github.com/aallbrig/allbctl/pkg/ansible"
	"github.com/spf13/cobra"
	"log"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "code generation for ansible config",
	Run: func(cmd *cobra.Command, args []string) {
		var config = ansible.Config{}
		err := config.RenderFiles()
		if err != nil {
			log.Fatalf("Error rendering ansible config file: %v\n", err)
		}
	},
}

func init() {
	Cmd.AddCommand(configCmd)
}
