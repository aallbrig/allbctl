package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/aallbrig/allbctl/cmd/generate/ansible"
	"github.com/aallbrig/allbctl/cmd/generate/dockerfile"
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
		pkg.HelpTextIfEmpty(cmd, args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		for _, file := range pkg.FilesToGenerate {
			if pkg.WriteStdOut {
				fmt.Println(file.FileContents)
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatalf("Unable to get current directory: %v", err)
				}

				err = os.MkdirAll(path.Join(cwd, file.RelativeDir), os.ModePerm)
				if err != nil {
					log.Fatalf("Unable to create directories: %v", err)
				}

				err = ioutil.WriteFile(path.Join(cwd, file.RelativeDir, file.FileName), file.FileContents.Bytes(), os.ModePerm)
				if err != nil {
					log.Fatalf("Unable to write file: %v", err)
				}
			}
		}
		for _, action := range pkg.ActionsToExecute {
			fmt.Printf("Executing action: %s\n", action.Name)

			out, err := action.Cmd.Output()
			if err != nil {
				log.Fatalf("Unable to execute action: %v", err)
			}
			fmt.Printf("Action output:\n%v", string(out))
		}
	},
}

func init() {
	GenerateCmd.PersistentFlags().BoolVarP(&pkg.WriteStdOut, "stdout", "o", false, "")

	GenerateCmd.AddCommand(ansible.Cmd)
	GenerateCmd.AddCommand(dockerfile.Cmd)
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
