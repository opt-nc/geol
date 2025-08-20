/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"c"},
	Short:   "Delete the locally cached products file.",
	Long: `Removes the local products cache file from the user's config directory.

This command is useful for clearing the cached list of products and their aliases previously downloaded from the endoflife.date API. The cache file is stored in the config directory under geol/products.json. If the file does not exist, a message is displayed.`,
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := os.UserConfigDir()
		if err != nil {
			cmd.PrintErrln("Error retrieving config directory:", err)
			return
		}
		productsPath := configDir + "/geol/products.json"

		if _, err := os.Stat(productsPath); os.IsNotExist(err) {
			cmd.Println("No cache file to delete.")
			return
		}
		if err := os.Remove(productsPath); err != nil {
			cmd.PrintErrln("Error deleting cache file:", err)
			return
		}
		cmd.Println("Cache file deleted.")
	},
}

func init() {
	cacheCmd.AddCommand(clearCmd)
}
