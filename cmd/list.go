package cmd

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	// Set up pretty console writer for phuslu/log
	log.DefaultLogger.Writer = &log.ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    true,
		EndWithMessage: true,
	}
	rootCmd.AddCommand(listCmd)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all cached product names.",
	Long:    `Displays the list of all product names currently available on https://endoflife.date.`,
	Example: `geol list
geol l`,
	Run: func(cmd *cobra.Command, args []string) {
		// List the cached products
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products path")
			return
		}
		// Ensure cache exists, create if missing
		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error ensuring cache exists")
			return
		}

		modTime := info.ModTime()
		utilities.CheckCacheTimeAndUpdate(cmd, modTime)

		// Read and parse the products file
		data, err := os.ReadFile(productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error reading products file - try running `geol cache refresh`")
			return
		}

		var products utilities.ProductsFile
		if err := json.Unmarshal(data, &products); err != nil {
			log.Error().Err(err).Msg("Error parsing JSON")
			log.Warn().Msg("Trying to refresh the cache now...")
			if err := utilities.FetchAndSaveProducts(cmd); err != nil {
				log.Error().Err(err).Msg("Error refreshing cache")
				return
			}
			log.Info().Msg("Cache refreshed successfully. Now getting the products...")
			data, err = os.ReadFile(productsPath)
			if err != nil {
				log.Error().Err(err).Msg("Error reading products file after refresh")
				return
			}
			if err := json.Unmarshal(data, &products); err != nil {
				log.Error().Err(err).Msg("Error parsing JSON after refresh")
				return
			}
		}

		// Print the list of products
		cmd.Println("Cached products:")
		cmd.Println("")
		// Collect and sort product names
		var names []string
		for name := range products.Products {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			cmd.Printf("%s\n", name)
		}

	},
}
