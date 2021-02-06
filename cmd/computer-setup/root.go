/*
Copyright © 2020 Andrew Allbright

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	computer_setup "github.com/aallbrig/allbctl/pkg/computer-setup"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
)

var RootCmd = &cobra.Command{
	Use:   "computer-setup",
	Short: "Setup computer in expected way",
	Run: func(cmd *cobra.Command, args []string) {
		usrHomeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal("Error getting user home directory")
		}
		err = computer_setup.DirectoryForSourceCode(usrHomeDir)
		if err != nil {
			log.Fatal("Error creating src directory")
		}

		log.Println("src directory exists in home directory")
	},
}
