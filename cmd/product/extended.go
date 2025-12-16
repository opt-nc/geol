package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/charmbracelet/x/term"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	extendedCmd.Flags().IntP("number", "n", 10, "Number of latest versions to display (default: 10, 0 to show all)")
}

// extendedCmd represents the extended command
var extendedCmd = &cobra.Command{
	Use:     "extended",
	Aliases: []string{"e"},
	Example: `# Show the latest 10 versions of Golang and Kubernetes
geol product extended golang k8s
# Show the latest 15 versions of Quarkus
geol product extended quarkus -n 15
# Redirect output to a markdown file
geol product extended quarkus > quarkus-eol.md`,
	Short: "Display extended release information for specified products (latest 10 versions by default).",
	Long:  `Retrieve and display detailed release data for one or more products, including cycle, release dates, support periods, and end-of-life information. By default, the latest 10 versions are shown for each product; use the --number flag to display the latest n versions instead. Results are formatted in a styled table for easy reading. Products must exist in the local cache or be available via the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		numberFlag, _ := cmd.Flags().GetInt("number")
		mdFlag := !term.IsTerminal(os.Stdout.Fd()) // detect if output is not a terminal

		if numberFlag < 0 {
			log.Fatal().Msg("The number of rows must be zero or positive.")
		}

		if len(args) == 0 {
			log.Fatal().Msg("Please specify at least one product.")
		}

		utilities.AnalyzeCacheProductsValidity(cmd)

		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving products path")
		}

		products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving products from cache")
		}

		var allProducts []productReleases

		for _, prod := range args {
			found := false
			for name, aliases := range products.Products {
				if strings.EqualFold(prod, name) {
					found = true
					prod = name
					break
				}
				for _, alias := range aliases {
					if strings.EqualFold(prod, alias) {
						found = true
						prod = name
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				log.Error().Msgf("Product %s not found in the API.", prod)
				continue // product not found in cache
			}

			prodData, err := fetchProductData(prod)
			if err != nil {
				log.Fatal().Err(err).Msg("Error fetching product data")
			}
			allProducts = append(allProducts, prodData)
		}

		if len(allProducts) == 0 {
			log.Fatal().Msg("None of the products were found in the API.")
		}

		// Render tables for all products
		for i, prod := range allProducts {
			renderProductTable(prod, numberFlag, mdFlag, i == 0)
		}
	},
}

// renderProductTable displays a formatted table for a single product's release information
func renderProductTable(prod productReleases, numberFlag int, mdFlag bool, isFirst bool) {
	// Print as a title "# Products" for the first product
	if isFirst {
		mainTitle := lipgloss.NewStyle().
			Bold(true).Foreground(lipgloss.Color("#FFFF88")).
			Background(lipgloss.Color("#5F5FFF")).
			Render("# Products")
		_, _ = lipgloss.Println(mainTitle)
	}

	styledTitle := lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("#00AFF8")).
		Render("\n## " + prod.Name + "\n")
	_, _ = lipgloss.Println(styledTitle)

	// Determine which columns have at least one value
	showName, showReleaseDate, showLatestName, showLatestDate, showEoasFrom, showEolFrom := false, false, false, false, false, false
	for _, r := range prod.Releases {
		if r.Name != "" {
			showName = true
		}
		if r.ReleaseDate != "" {
			showReleaseDate = true
		}
		if r.LatestName != "" {
			showLatestName = true
		}
		if r.LatestDate != "" {
			showLatestDate = true
		}
		if r.EoasFrom != "" {
			showEoasFrom = true
		}
		if r.EolFrom != "" {
			showEolFrom = true
		}
	}

	var columns []string
	if showName {
		columns = append(columns, "Cycle")
	}
	if showReleaseDate {
		columns = append(columns, "Release")
	}
	if showLatestName {
		columns = append(columns, "Latest")
	}
	if showLatestDate {
		columns = append(columns, "Latest Release")
	}
	if showEoasFrom {
		columns = append(columns, "Support")
	}
	if showEolFrom {
		columns = append(columns, "EOL")
	}

	if len(columns) == 0 {
		_, _ = lipgloss.Println(lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244")).Render("No release data available."))
		return
	}

	// Create and display the table with lipgloss/table
	t := table.New()
	// Bold style for column headers
	headerStyle := lipgloss.NewStyle().Bold(true)
	styledHeaders := make([]string, len(columns))
	for i, col := range columns {
		styledHeaders[i] = headerStyle.Render(col)
	}
	t.Headers(styledHeaders...)
	// Limit the number of displayed rows
	displayCount := numberFlag
	if displayCount == 0 || displayCount > len(prod.Releases) {
		displayCount = len(prod.Releases)
	}
	for i := 0; i < displayCount; i++ {
		r := prod.Releases[i]
		var row []string
		today := utilities.TodayDateString() // format: YYYY-MM-DD
		// Helper to color a date string
		colorDate := func(date string) string {
			if date == "" {
				return ""
			}
			styleRed := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			styleGreen := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
			if date <= today {
				return styleRed.Render(date)
			}
			return styleGreen.Render(date)
		}

		// Helper to strikethrough a string if EOL is before today
		strikethroughIfEOL := func(val string) string {
			if r.EolFrom != "" && r.EolFrom < today {
				if mdFlag {
					return "~~" + val + "~~"
				}
				return "\x1b[9m" + val + "\x1b[0m"
			}
			return val
		}

		if showName {
			nameWithBadge := r.Name
			nameWithBadge = strikethroughIfEOL(nameWithBadge)
			var ltsBadge string
			if r.LTS {
				badgeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Padding(0, 1)
				today := utilities.TodayDateString()
				if r.EolFrom == "" || r.EolFrom >= today {
					badgeStyle = badgeStyle.Background(lipgloss.Color("34")).PaddingLeft(0).PaddingRight(0) // dark green
				} else {
					badgeStyle = badgeStyle.Background(lipgloss.Color("196")).PaddingLeft(0).PaddingRight(0) // rouge
				}
				ltsBadge = badgeStyle.Render("LTS")
				nameWithBadge += " " + ltsBadge
			}
			row = append(row, nameWithBadge)
		}
		if showReleaseDate {
			row = append(row, r.ReleaseDate)
		}
		if showLatestName {
			row = append(row, strikethroughIfEOL(r.LatestName))
		}
		if showLatestDate {
			row = append(row, strikethroughIfEOL(r.LatestDate))
		}
		if showEoasFrom {
			row = append(row, colorDate(r.EoasFrom))
		}
		if showEolFrom {
			row = append(row, colorDate(r.EolFrom))
		}
		t.Row(row...)
	}
	// If not all rows are shown, add a final row with '...'
	if displayCount < len(prod.Releases) {
		dotsRow := make([]string, len(columns))
		for i := range dotsRow {
			dotsRow[i] = "..."
		}
		t.Row(dotsRow...)
	}
	if !mdFlag {
		t.Border(lipgloss.RoundedBorder())
		t.BorderBottom(true)
	} else {
		t.Border(lipgloss.MarkdownBorder())
		t.BorderBottom(false)
	}
	t.BorderTop(false)
	t.BorderLeft(false)
	t.BorderRight(false)
	t.BorderStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("63")))
	t.StyleFunc(func(row, col int) lipgloss.Style {
		padding := 1
		return lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Align(lipgloss.Left).Padding(0, padding)
	})
	renderedTable := t.Render()
	_, _ = lipgloss.Println(renderedTable)
	// Always show a summary line below the table
	tableLines := strings.Split(renderedTable, "\n")
	maxLen := 0
	for _, l := range tableLines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	summary := fmt.Sprintf("%d rows (%d shown)", len(prod.Releases), displayCount)
	_, _ = lipgloss.Println(summary)
}

// fetchProductData retrieves product release data from the API
func fetchProductData(productName string) (productReleases, error) {
	url := utilities.ApiUrl + "products/" + productName
	resp, err := http.Get(url)
	if err != nil {
		return productReleases{}, fmt.Errorf("error requesting %s: %w", productName, err)
	}
	body, err := io.ReadAll(resp.Body)
	if cerr := resp.Body.Close(); cerr != nil {
		return productReleases{}, fmt.Errorf("error closing HTTP body for %s: %w", productName, cerr)
	}
	if err != nil {
		return productReleases{}, fmt.Errorf("error reading response for %s: %w", productName, err)
	}
	if resp.StatusCode != 200 {
		return productReleases{}, fmt.Errorf("product %s not found on the API", productName)
	}

	var apiResp ApiRespExtended
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return productReleases{}, fmt.Errorf("error decoding JSON for %s: %w", productName, err)
	}

	var releases []ReleaseInfo
	for _, r := range apiResp.Result.Releases {
		releases = append(releases, ReleaseInfo{
			Name:        r.Name,
			ReleaseDate: r.ReleaseDate,
			LatestName:  r.Latest.Name,
			LatestDate:  r.Latest.Date,
			EoasFrom:    r.EoasFrom,
			EolFrom:     r.EolFrom,
			LTS:         r.IsLTS,
		})
	}

	return productReleases{
		Name:     apiResp.Result.Name,
		Releases: releases,
	}, nil
}
