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
	"fmt"
	"github.com/aallbrig/allbctl/cmd/generate/ansible"
	"github.com/aallbrig/allbctl/cmd/generate/git"
	"github.com/aallbrig/allbctl/cmd/generate/golang"
	"github.com/aallbrig/allbctl/cmd/generate/java"
	"github.com/aallbrig/allbctl/cmd/generate/kubernetes"
	"github.com/aallbrig/allbctl/cmd/generate/node"
	"github.com/aallbrig/allbctl/cmd/generate/python"
	"github.com/aallbrig/allbctl/cmd/generate/ruby"
	"github.com/aallbrig/allbctl/cmd/generate/scala"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "root command for code generation commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")
	},
}

func init() {
	generateCmd.AddCommand(ansible.Cmd)
	generateCmd.AddCommand(git.Cmd)
	generateCmd.AddCommand(golang.Cmd)
	generateCmd.AddCommand(java.Cmd)
	generateCmd.AddCommand(kubernetes.Cmd)
	generateCmd.AddCommand(node.Cmd)
	generateCmd.AddCommand(python.Cmd)
	generateCmd.AddCommand(ruby.Cmd)
	generateCmd.AddCommand(scala.Cmd)
	rootCmd.AddCommand(generateCmd)
}
