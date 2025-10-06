package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/fatih/color"
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
	listCmd.Flags().BoolP("tree", "t", false, "List all products including aliases in a tree structure.")
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

		treeFlag, _ := cmd.Flags().GetBool("tree")

		var names []string
		for name := range products.Products {
			names = append(names, name)
		}
		sort.Strings(names)
		productColor := color.New(color.Bold)

		if treeFlag {
			// Print the list of products with aliases in a tree structure using lipgloss
			productTree := tree.Root(".")

			for _, name := range names {
				aliases := products.Products[name]
				if len(aliases) > 0 {
					aliases = aliases[1:]
				}
				sort.Strings(aliases)
				t := tree.New().Root(productColor.Sprint(name))
				for _, item := range aliases {
					t.Child(item)
				}
				productTree.Child(t)
			}
			if _, err := fmt.Fprintln(os.Stdout, productTree.String()); err != nil {
				log.Error().Err(err).Msg("Error writing product tree to stdout")
				return
			}
		} else {
			// Print the list of products with a green '+ product' prefix using lipgloss
			plusStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render("+")
			for _, name := range names {
				if _, err := fmt.Fprintf(os.Stdout, "%s %s\n", plusStyle, productColor.Sprint(name)); err != nil {
					log.Error().Err(err).Msg("Error writing product name to stdout")
					return
				}
			}
		}

		nbProducts := strconv.Itoa(len(names))
		nbProductsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := fmt.Fprintf(os.Stdout, "\n%s products listed\n", nbProductsStyle.Render(nbProducts)); err != nil {
			log.Error().Err(err).Msg("Error writing product count to stdout")
			return
		}

	},
}
