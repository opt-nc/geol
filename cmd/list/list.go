package list

import (
	"github.com/spf13/cobra"

	"github.com/phuslu/log"

	"github.com/opt-nc/geol/cmd/list/items"
	"github.com/opt-nc/geol/utilities"
)

func init() {
	utilities.InitLogger()
	ListCmd.AddCommand(items.ProductsCmd)
}

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List cached data.",
	Long:    `Displays the list of cached data such as products, tags, and categories.`,
	Example: `geol list products`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error().Err(err).Msg("Error displaying help")
		}
	},
}
