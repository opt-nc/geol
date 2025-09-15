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

// getProductsWithCacheRefresh tries to unmarshal products from file, refreshes cache if needed, and returns the products.
func getProductsWithCacheRefresh(cmd *cobra.Command, productsPath string) (utilities.ProductsFile, error) {
	var products utilities.ProductsFile
	if err := readAndUnmarshalProducts(productsPath, &products); err != nil {
		log.Error().Err(err).Msg("Error parsing JSON")
		log.Warn().Msg("Trying to refresh the cache now...")
		if err := utilities.FetchAndSaveProducts(cmd); err != nil {
			log.Error().Err(err).Msg("Error refreshing cache")
			return products, err
		}
		log.Info().Msg("Cache refreshed successfully. Now getting the products...")
		if err := readAndUnmarshalProducts(productsPath, &products); err != nil {
			log.Error().Err(err).Msg("Error parsing JSON after refresh")
			return products, err
		}
	}
	return products, nil
}

// readAndUnmarshalProducts lit le fichier et fait l'unmarshal JSON dans products.
func readAndUnmarshalProducts(productsPath string, products *utilities.ProductsFile) error {
	data, err := os.ReadFile(productsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error reading products file")
		return err
	}
	if err := json.Unmarshal(data, products); err != nil {
		return err
	}
	return nil
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

		products, err := getProductsWithCacheRefresh(cmd, productsPath)
		if err != nil {
			return
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
