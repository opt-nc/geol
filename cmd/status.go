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
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
