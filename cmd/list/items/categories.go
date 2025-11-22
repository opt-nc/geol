package items

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/fatih/color"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

// CategoriesCmd represents the categories command
var CategoriesCmd = &cobra.Command{
	Use:     "categories",
	Aliases: []string{"c"},
	Short:   "List all cached category names.",
	Long:    `Displays the list of all category names currently available in the cache.`,
	Example: `geol list categories
geol l c`,
	Run: func(cmd *cobra.Command, args []string) {
		// List the cached categories
		utilities.AnalyzeCacheProductsValidity(cmd)
		categoriesPath, err := utilities.GetCategoriesPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving categories path")
			os.Exit(1)
		}

		categories, err := utilities.GetCategoriesWithCacheRefresh(cmd, categoriesPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving categories from cache")
			os.Exit(1)
		}

		var names []string
		for name := range categories {
			names = append(names, name)
		}
		sort.Strings(names)
		categoryColor := color.New(color.Bold)

		plusStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render("+")
		for _, name := range names {
			if _, err := fmt.Fprintf(os.Stdout, "%s %s\n", plusStyle, categoryColor.Sprint(name)); err != nil {
				log.Error().Err(err).Msg("Error writing category name to stdout")
				os.Exit(1)
			}
		}

		nbCategories := strconv.Itoa(len(names))
		nbCategoriesStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := fmt.Fprintf(os.Stdout, "\n%s categories listed\n", nbCategoriesStyle.Render(nbCategories)); err != nil {
			log.Error().Err(err).Msg("Error writing category count to stdout")
			os.Exit(1)
		}
	},
}

func init() {
}
