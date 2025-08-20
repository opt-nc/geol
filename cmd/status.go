/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s"},
	Short:   "Show information about the local products cache file.",
	Long: `Displays the status of the local products cache file stored in the user's config directory.

This command prints the last update date and the number of products currently cached in geol/products.json. It helps verify if the cache is present and up to date.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Define the cache file path
		configDir, err := os.UserConfigDir()
		if err != nil {
			cmd.PrintErrln("Error retrieving config directory:", err)
			return
		}
		productsPath := configDir + "/geol/products.json"

		// Check if the file exists
		info, err := os.Stat(productsPath)
		if err != nil {
			cmd.PrintErrln("Cache file not found:", productsPath)
			return
		}

		// Print the last update date
		cmd.Printf("Cache last update: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))

		// Read and parse the file to count the keys
		data, err := os.ReadFile(productsPath)
		if err != nil {
			cmd.PrintErrln("Error reading cache file:", err)
			return
		}
		type ProductsFile struct {
			Products map[string][]string `json:"products"`
		}
		var products ProductsFile
		if err := json.Unmarshal(data, &products); err != nil {
			cmd.PrintErrln("Error parsing JSON:", err)
			return
		}
		cmd.Printf("Number of items in cache: %d\n", len(products.Products))
	},
}

func init() {
	cacheCmd.AddCommand(statusCmd)
}
