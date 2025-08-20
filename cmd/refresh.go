package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:     "refresh",
	Aliases: []string{"r"},
	Short:   "Download the latest list of products and their aliases from the endoflife.date API and save it locally.",
	Long: `Fetches the current list of products and their aliases from the endoflife.date API, processes the data into a local JSON file under the user's config directory, and ensures the file is updated with the latest information.

This command is useful for keeping the local product list in sync with the upstream source for further use by the application. The resulting file is stored in the config directory under geol/products.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		productsPath, err := getProductsPath()
		if err != nil {
			cmd.PrintErrln("Error retrieving products path:", err)
			return
		}

		// Ensure the directory exists
		if err := createDirectoryIfNotExists(productsPath); err != nil {
			cmd.PrintErrln("Error ensuring directory exists:", err)
			return
		}

		// HTTP GET request (extracted)
		resp, err := getAPIResponse(ApiUrl + "products")
		if err != nil {
			cmd.PrintErrln("Error during HTTP request:", err)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				cmd.PrintErrln("Error closing response body:", err)
			}
		}()

		// Define structure to parse the response (Product is now imported from tools.go)
		type ApiResponse struct {
			Result []Product `json:"result"`
		}
		var apiResp ApiResponse

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&apiResp); err != nil {
			cmd.PrintErrln("Error decoding JSON:", err)
			return
		}

		// Build the hashmap
		type ProductsFile struct {
			Products map[string][]string `json:"products"`
		}
		products := ProductsFile{Products: make(map[string][]string)}
		for _, p := range apiResp.Result {
			aliases := []string{p.Name}
			aliases = append(aliases, p.Aliases...)
			products.Products[p.Name] = aliases
		}

		// Marshal the data to JSON
		data, err := json.MarshalIndent(products, "", "  ")
		if err != nil {
			cmd.PrintErrln("Error serializing JSON:", err)
			return
		}

		// Remove the file if it exists (extracted)
		if err := removeFileIfExists(productsPath); err != nil {
			cmd.PrintErrln("Error removing old file:", err)
			return
		}
		// Save to file
		if err := os.WriteFile(productsPath, data, 0644); err != nil {
			cmd.PrintErrln("Error writing file:", err)
			return
		}
		// Print the number of products written and elapsed time
		elapsed := time.Since(start).Milliseconds()
		cmd.Printf("Products file updated from API. \nNumber of products: %d \n(elapsed time: %d ms)\n", len(products.Products), elapsed)
	},
}

func init() {
	cacheCmd.AddCommand(refreshCmd)
}
