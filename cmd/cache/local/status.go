package local

import (
	"encoding/json"
	"os"

	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	utilities.InitLogger()
}

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s"},
	Short:   "Show information about the local products cache file.",
	Long: `Displays the status of the local products cache file stored in the user's config directory.

This command prints the last update date and the number of products currently cached in geol/products.json. It helps verify if the cache is present and up to date.`,
	Run: func(cmd *cobra.Command, args []string) {
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products path")
			return
		}

		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			return
		}

		modTime := info.ModTime()

		utilities.CheckCacheTimeAndUpdate(cmd, modTime)

		data, err := os.ReadFile(productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error reading cache file")
			return
		}
		var products utilities.ProductsFile
		if err := json.Unmarshal(data, &products); err != nil {
			log.Error().Err(err).Msg("Error parsing JSON")
			return
		}
		log.Info().Int("Number of products", len(products.Products)).Msg("")
	},
}

func init() {
}
