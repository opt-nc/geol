package local

import (
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
		if err == nil {
			if info, err2 := utilities.EnsureCacheExistsGeneric(productsPath, cmd); err2 == nil {
				modTime := info.ModTime()
				log.Info().Msg("Cache last updated " + modTime.Format("2006-01-02 15:04:05"))
			}
		}
		utilities.AnalyzeCacheProductsValidity(cmd)
		var errorOccurred = false

		products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products from cache")
			errorOccurred = true
		}
		log.Info().Int("Number of products", len(products.Products)).Msg("")

		tagsPath, err := utilities.GetTagsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving tags path")
			errorOccurred = true
		}

		tags, err := utilities.GetTagsWithCacheRefresh(cmd, tagsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving tags from cache")
			errorOccurred = true
		}
		log.Info().Int("Number of tags", len(tags)).Msg("")

		categoriesPath, err := utilities.GetCategoriesPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving categories path")
			errorOccurred = true
		}

		categories, err := utilities.GetCategoriesWithCacheRefresh(cmd, categoriesPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving categories from cache")
			errorOccurred = true
		}
		log.Info().Int("Number of categories", len(categories)).Msg("")

		if errorOccurred {
			os.Exit(1)
		}
	},
}
