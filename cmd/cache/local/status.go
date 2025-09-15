package local

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s"},
	Short:   "Show information about the local products cache file.",
	Long: `Displays the status of the local products cache file stored in the user's config directory.

This command prints the last update date and the number of products currently cached in geol/products.json. It helps verify if the cache is present and up to date.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the cache file path using shared function
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			cmd.PrintErrln("Error retrieving products path:", err)
			return
		}

		// Ensure cache exists, create if missing
		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			return
		}

		// Print the last update date with system timezone
		modTime := info.ModTime()
		zone, offset := modTime.Zone()
		tz := zone
		if offset != 0 {
			// Format offset as "+02:00" or "-07:00"
			sign := "+"
			if offset < 0 {
				sign = "-"
				offset = -offset
			}
			hours := offset / 3600
			minutes := (offset % 3600) / 60
			tz = fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
		}
		cmd.Printf("Last cache update: %s %s\n", modTime.Format("2006-01-02 15:04:05"), tz)

		utilities.CheckCacheTimeAndUpdate(cmd, modTime)

		// Read and parse the file to count the keys
		data, err := os.ReadFile(productsPath)
		if err != nil {
			cmd.PrintErrln("Error reading cache file:", err)
			return
		}
		var products utilities.ProductsFile
		if err := json.Unmarshal(data, &products); err != nil {
			cmd.PrintErrln("Error parsing JSON:", err)
			return
		}
		cmd.Printf("Number of items in cache: %d\n", len(products.Products))
	},
}

func init() {
}
