/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := os.UserConfigDir()
		if err != nil {
			cmd.PrintErrln("Error retrieving config directory:", err)
			return
		}
		productsPath := configDir + "/geol/products.json"

		// Create the directory if it doesn't exist
		if _, err := os.Stat(configDir + "/geol"); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir+"/geol", 0755); err != nil {
				cmd.PrintErrln("Error creating directory:", err)
				return
			}
		}

		// HTTP GET request
		resp, err := http.Get("https://endoflife.date/api/v1/products")
		if err != nil {
			cmd.PrintErrln("Error during HTTP request:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			cmd.PrintErrln("Unexpected HTTP status:", resp.Status)
			return
		}

		// Define structures to parse the response
		type Product struct {
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
		}
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

		// Save to file
		data, err := json.MarshalIndent(products, "", "  ")
		if err != nil {
			cmd.PrintErrln("Error serializing JSON:", err)
			return
		}

		// Remove the file if it exists
		if _, err := os.Stat(productsPath); err == nil {
			if err := os.Remove(productsPath); err != nil {
				cmd.PrintErrln("Error removing old file:", err)
				return
			}
		}
		if err := os.WriteFile(productsPath, data, 0644); err != nil {
			cmd.PrintErrln("Error writing file:", err)
			return
		}
		cmd.Println("Products file updated from API.")
	},
}

func init() {
	cacheCmd.AddCommand(refreshCmd)
}
