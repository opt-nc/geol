package cmd

import (
	"sort"

	"github.com/charmbracelet/lipgloss"
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

		products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products from cache")
			return
		}

		// Print the list of products with a green '+ product' prefix using lipgloss
		var names []string
		for name := range products.Products {
			names = append(names, name)
		}
		sort.Strings(names)
		plusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("+")
		for _, name := range names {
			cmd.Printf("%s %s\n", plusStyle, name)
		}
		cmd.Printf("\n%d products listed\n", len(names))

	},
}
