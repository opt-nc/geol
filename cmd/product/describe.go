package product

import (
	"io"
	"net/http"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour/v2"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
}

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:     "describe [product]",
	Aliases: []string{"d"},
	Example: `geol product describe nodejs`,
	Short:   "Display the product summary",
	Long:    `Display the description for a single given product. Useful for quickly viewing product summary.`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Error().Msg("Please specify exactly one product.")
			os.Exit(1)
		}
		prodArg := args[0]

		// Check the cache
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

		// Find the main product name (key)
		var mainName string
		found := false
		for name, aliases := range products.Products {
			if strings.EqualFold(prodArg, name) {
				mainName = name
				found = true
				break
			}
			for _, alias := range aliases {
				if strings.EqualFold(prodArg, alias) {
					mainName = name
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			log.Error().Msgf("Product '%s' not found in the API.", prodArg)
			os.Exit(1)
		}

		// Build the markdown URL
		mdUrl := "https://raw.githubusercontent.com/endoflife-date/endoflife.date/refs/heads/master/products/" + mainName + ".md"

		// Retrieve the Markdown content
		resp, err := http.Get(mdUrl)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching markdown of the product " + mainName)
			os.Exit(1)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Error().Err(err).Msg("Error closing response body")
			}
		}()

		if resp.StatusCode != http.StatusOK {
			log.Error().Msgf("Failed to fetch markdown. Status: %s", resp.Status)
			os.Exit(1)
		}

		mdBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error reading markdown")
			os.Exit(1)
		}

		// Extract the description between the second '---' and the first empty line after
		mdLines := strings.Split(string(mdBytes), "\n")
		sepCount := 0
		descLines := []string{}
		collecting := false
		for _, line := range mdLines {
			if strings.TrimSpace(line) == "---" {
				sepCount++
				if sepCount == 2 {
					collecting = true
					continue
				}
			}
			if collecting {
				// Stop if a new markdown section (title) is encountered
				if strings.HasPrefix(strings.TrimSpace(line), "#") && len(descLines) > 0 {
					break
				}
				descLines = append(descLines, line)
			}
		}
		desc := strings.TrimRight(strings.Join(descLines, "\n"), "\n")
		if desc == "" {
			log.Error().Msg("No description found in markdown.")
			os.Exit(1)
		}

		// Add iCalendar feed information
		desc += "\n\nYou can subscribe to the iCalendar feed at `webcal://endoflife.date/calendar/" + mainName + ".ics`"

		// Add A JSON version of this page is available at /api/v1/products/neo4j/
		desc += "\n\nA JSON version of this page is available at `https://endoflife.date/api/v1/products/" + mainName + "`"

		// Print a product title as in extended: # ProductName, with color and background
		styledTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFF88")).
			Background(lipgloss.Color("#5F5FFF")).
			Render("# " + mainName)
		if _, err := os.Stdout.Write([]byte(styledTitle)); err != nil {
			log.Error().Err(err).Msg("Error writing styled title")
			os.Exit(1)
		}

		// Glamour rendering only on the description
		out, err := glamour.RenderWithEnvironmentConfig(desc)
		if err != nil {
			log.Error().Err(err).Msg("Error rendering markdown")
			os.Exit(1)
		}
		if _, err := os.Stdout.Write([]byte(out)); err != nil {
			log.Error().Err(err).Msg("Error writing rendered markdown")
			os.Exit(1)
		}
	},
}
