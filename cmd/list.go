package cmd

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/opt-nc/geol/cmd/cache/utilities"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all cached product names.",
	Long: `Displays a sorted list of all product names currently stored in the local cache file.

This command reads the products cache (geol/products.json) and prints each product name to the console. If the cache file is missing or cannot be read, an error message is shown and a suggestion to refresh the cache is provided. Use this command to quickly view which products are available in your local cache.`,
	Example: `geol cache list
geol cl`,
	Run: func(cmd *cobra.Command, args []string) {
		// List the cached products
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			cmd.PrintErrln("Error retrieving products path:", err)
			return
		}

		// Read and parse the products file
		data, err := os.ReadFile(productsPath)
		if err != nil {
			cmd.PrintErrln("Error reading products file:", err, "- try running `geol cache refresh`")
			return
		}

		var products utilities.ProductsFile
		if err := json.Unmarshal(data, &products); err != nil {
			cmd.PrintErrln("Error parsing JSON:", err)
			return
		}

		// Print the list of products
		cmd.Println("Cached products:")
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

func init() {
	rootCmd.AddCommand(listCmd)
}
