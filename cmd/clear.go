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
	Use:   "clear",
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
