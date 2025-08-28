package local

import (
	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

// RefreshCmd represents the refresh command
var RefreshCmd = &cobra.Command{
	Use:     "refresh",
	Aliases: []string{"r"},
	Short:   "Download the latest list of products and their aliases from the endoflife.date API and save it locally.",
	Long: `Fetches the current list of products and their aliases from the endoflife.date API, processes the data into a local JSON file under the user's config directory, and ensures the file is updated with the latest information.

This command is useful for keeping the local product list in sync with the upstream source for further use by the application. The resulting file is stored in the config directory under geol/products.json.`,
	Example: `geol cache refresh
geol c r`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utilities.FetchAndSaveProducts(cmd); err != nil {
			return
		}
	},
}

func init() {
}
