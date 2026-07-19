package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ciGithubCmd represents the ci-github command
var ciGithubCmd = &cobra.Command{
	Use:     "ci-github",
	Aliases: []string{"ci-gh"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ci-github called")
	},
}

func init() {
	rootCmd.AddCommand(ciGithubCmd)

}
