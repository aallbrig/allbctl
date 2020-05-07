package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate completions for bash or zsh",
	Long: `To load completion run

. <(allbctl completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(allbctl completion)
# ~/.zshrc
. <(allbctl completion zsh)
`,
	Run: func(cmd *cobra.Command, args []string) {
		completionType := ""
		if len(args) > 0 {
			completionType = args[0]
		}

		switch completionType {
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		default:
			rootCmd.GenBashCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
