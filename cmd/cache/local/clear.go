package local

import (
	"os"

	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	// Initialisation du logger global

}

// ClearCmd represents the clear command
var ClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"c"},
	Short:   "Delete the locally cached products file.",
	Long: `Removes the local products cache file from the user's config directory.

This command is useful for clearing the cached list of products and their aliases previously downloaded from the endoflife.date API. The cache file is stored in the config directory under geol/products.json. If the file does not exist, a message is displayed.`,
	Run: func(cmd *cobra.Command, args []string) {
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products path")
			os.Exit(1)
		}
		if err := utilities.RemoveFileIfExists(productsPath); err != nil {
			log.Error().Err(err).Msg("Error deleting cache file")
			os.Exit(1)
		}
		log.Info().Str("path", productsPath).Msg("Cache file removed.")

		tagsPath, err := utilities.GetTagsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving tags path")
			os.Exit(1)
		}
		if err := utilities.RemoveFileIfExists(tagsPath); err != nil {
			log.Error().Err(err).Msg("Error deleting tags file")
			os.Exit(1)
		}
		log.Info().Str("path", tagsPath).Msg("Tags file removed.")

		categoriesPath, err := utilities.GetCategoriesPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving categories path")
			os.Exit(1)
		}
		if err := utilities.RemoveFileIfExists(categoriesPath); err != nil {
			log.Error().Err(err).Msg("Error deleting categories file")
			os.Exit(1)
		}
		log.Info().Str("path", categoriesPath).Msg("Categories file removed.")
	},
}
