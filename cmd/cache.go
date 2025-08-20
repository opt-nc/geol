package cmd

import (
	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:     "cache",
	Aliases: []string{"c"},
	Short:   "Update the local cache",
	Long:    `The cache command is used to update the local cache in the user's config directory, in 'geol/products.json'. It provides subcommands to refresh, clear, and check the status of the cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)

}
