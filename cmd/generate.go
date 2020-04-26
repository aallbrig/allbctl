package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/aallbrig/allbctl/cmd/generate/ansible"
	"github.com/aallbrig/allbctl/cmd/generate/git"
	"github.com/aallbrig/allbctl/cmd/generate/golang"
	"github.com/aallbrig/allbctl/cmd/generate/java"
	"github.com/aallbrig/allbctl/cmd/generate/kubernetes"
	"github.com/aallbrig/allbctl/cmd/generate/node"
	"github.com/aallbrig/allbctl/cmd/generate/python"
	"github.com/aallbrig/allbctl/cmd/generate/ruby"
	"github.com/aallbrig/allbctl/cmd/generate/scala"
	"github.com/aallbrig/allbctl/pkg"
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "root command for code generation commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// TODO: Add option to output to stdout or write to filesystem
		fmt.Println("Files to be generated:")
		for _, file := range pkg.FilesToGenerate {
			fmt.Println(filepath.Join(file.RelativeDir, file.FileName))
		}
		fmt.Println()

		for _, file := range pkg.FilesToGenerate {
			fmt.Println(filepath.Join(file.RelativeDir, file.FileName))
			fmt.Println(file.FileContents)
		}
	},
}

func init() {
	GenerateCmd.AddCommand(ansible.Cmd)
	GenerateCmd.AddCommand(git.Cmd)
	GenerateCmd.AddCommand(golang.Cmd)
	GenerateCmd.AddCommand(java.Cmd)
	GenerateCmd.AddCommand(kubernetes.Cmd)
	GenerateCmd.AddCommand(node.Cmd)
	GenerateCmd.AddCommand(python.Cmd)
	GenerateCmd.AddCommand(ruby.Cmd)
	GenerateCmd.AddCommand(scala.Cmd)
	rootCmd.AddCommand(GenerateCmd)
}
