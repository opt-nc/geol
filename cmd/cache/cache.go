package cache

import (
	"github.com/opt-nc/geol/cmd/cache/local"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

// CacheCmd represents the cache command
var CacheCmd = &cobra.Command{
	Use:     "cache",
	Aliases: []string{"c"},
	Short:   "Update the local cache",
	Long:    `The cache command is used to update the local cache in the user's config directory, in 'geol/products.json'. It provides subcommands to refresh, clear, and check the status of the cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error().Err(err).Msg("Error displaying help")
		}
	},
}

func init() {
	CacheCmd.AddCommand(local.StatusCmd)
	CacheCmd.AddCommand(local.RefreshCmd)
	CacheCmd.AddCommand(local.ClearCmd)
}
