package local

import (
	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

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
			cmd.PrintErrln("Error retrieving products path:", err)
			return
		}
		if err := utilities.RemoveFileIfExists(productsPath); err != nil {
			cmd.PrintErrln("Error deleting cache file:", err)
			return
		}
		cmd.Println("Local products cache cleared successfully.")
	},
}

func init() {
}
