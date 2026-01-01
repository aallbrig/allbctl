package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RuntimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Display detected runtimes and their versions",
	Long: `Display detected programming language runtimes and their versions.

This is the same output shown in the 'Runtimes:' section of 'allbctl status'.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintRuntimes()
	},
}

// PrintRuntimes outputs the runtimes in inline format (same as status command)
func PrintRuntimes() {
	runtimesInline := detectRuntimesInline()
	if runtimesInline != "" {
		fmt.Println(runtimesInline)
	} else {
		fmt.Println("No runtimes detected")
	}
}
