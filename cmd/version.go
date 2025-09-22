package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"}, // Alias for the command
	Short:   "Display the application version",
	Long:    `Display the application version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
