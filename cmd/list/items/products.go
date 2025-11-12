package items

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/fatih/color"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	ProductsCmd.Flags().BoolP("tree", "t", false, "List all products including aliases in a tree structure.")
}

// ProductsCmd represents the products command
var ProductsCmd = &cobra.Command{
	Use:     "products",
	Aliases: []string{"p"},
	Short:   "List all cached product names.",
	Long:    `Displays the list of all product names currently available on https://endoflife.date.`,
	Example: `geol list products
geol list products --tree
geol l p -t`,
	Run: func(cmd *cobra.Command, args []string) {
		// List the cached products
		utilities.AnalyzeCacheProductsValidity(cmd)
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products path")
			os.Exit(1)
		}

		products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products from cache")
			os.Exit(1)
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
				os.Exit(1)
			}
		} else {
			// Print the list of products with a green '+ product' prefix using lipgloss
			plusStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render("+")
			for _, name := range names {
				if _, err := fmt.Fprintf(os.Stdout, "%s %s\n", plusStyle, productColor.Sprint(name)); err != nil {
					log.Error().Err(err).Msg("Error writing product name to stdout")
					os.Exit(1)
				}
			}
		}

		nbProducts := strconv.Itoa(len(names))
		nbProductsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := fmt.Fprintf(os.Stdout, "\n%s products listed\n", nbProductsStyle.Render(nbProducts)); err != nil {
			log.Error().Err(err).Msg("Error writing product count to stdout")
			os.Exit(1)
		}
	},
}
