package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	extendedCmd.Flags().IntP("number", "n", 10, "Number of latest versions to display (default: 10, 0 to show all)")
	utilities.InitLogger()
}

// extendedCmd represents the extended command
var extendedCmd = &cobra.Command{
	Use:     "extended",
	Aliases: []string{"e"},
	Example: `geol product extended golang k8s
geol product extended quarkus -n 15`,
	Short:   "Display extended release information for specified products (latest 10 versions by default).",
	Long:    `Retrieve and display detailed release data for one or more products, including cycle, release dates, support periods, and end-of-life information. By default, the latest 10 versions are shown for each product; use the --number flag to display the latest n versions instead. Results are formatted in a styled table for easy reading. Products must exist in the local cache or be available via the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		numberFlag, _ := cmd.Flags().GetInt("number")

		if numberFlag < 0 {
			log.Error().Msg("The number of rows must be zero or positive.")
			os.Exit(1)
		}

		if len(args) == 0 {
			log.Error().Msg("Please specify at least one product.")
			os.Exit(1)
		}

		utilities.AnalyzeCacheValidity(cmd)

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

		type ProductReleases struct {
			Name     string
			Releases []struct {
				Name        string
				ReleaseDate string
				LatestName  string
				LatestDate  string
				EoasFrom    string
				EolFrom     string
				LTS         bool
			}
		}
		var allProducts []ProductReleases

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
				continue // product not found in cache
			}

			url := utilities.ApiUrl + "products/" + prod
			resp, err := http.Get(url)
			if err != nil {
				log.Error().Err(err).Msgf("Error requesting %s", prod)
				os.Exit(1)
			}
			body, err := io.ReadAll(resp.Body)
			if cerr := resp.Body.Close(); cerr != nil {
				log.Error().Err(cerr).Msgf("Error closing HTTP body for %s", prod)
				os.Exit(1)
			}
			if err != nil {
				log.Error().Err(err).Msgf("Error reading response for %s", prod)
				os.Exit(1)
			}
			if resp.StatusCode != 200 {
				log.Error().Msgf("Product %s not found on the API.", prod)
				os.Exit(1)
			}

			var apiResp struct {
				Result struct {
					Name     string `json:"name"`
					Releases []struct {
						Name        string `json:"name"`
						ReleaseDate string `json:"releaseDate"`
						Latest      struct {
							Name string `json:"name"`
							Date string `json:"date"`
						} `json:"latest"`
						EoasFrom string `json:"eoasFrom"`
						EolFrom  string `json:"eolFrom"`
						IsLTS    bool   `json:"isLTS"`
					} `json:"releases"`
				}
			}
			if err := json.Unmarshal(body, &apiResp); err != nil {
				log.Error().Err(err).Msgf("Error decoding JSON for %s", prod)
				os.Exit(1)
			}

			var releases []struct {
				Name        string
				ReleaseDate string
				LatestName  string
				LatestDate  string
				EoasFrom    string
				EolFrom     string
				LTS         bool
			}
			for _, r := range apiResp.Result.Releases {
				releases = append(releases, struct {
					Name        string
					ReleaseDate string
					LatestName  string
					LatestDate  string
					EoasFrom    string
					EolFrom     string
					LTS         bool
				}{
					Name:        r.Name,
					ReleaseDate: r.ReleaseDate,
					LatestName:  r.Latest.Name,
					LatestDate:  r.Latest.Date,
					EoasFrom:    r.EoasFrom,
					EolFrom:     r.EolFrom,
					LTS:         r.IsLTS,
				})
			}
			allProducts = append(allProducts, ProductReleases{
				Name:     apiResp.Result.Name,
				Releases: releases,
			})
		}

		if len(allProducts) == 0 {
			log.Error().Msg("Aucun produit trouvé dans le cache ou l'API.")
			os.Exit(1)
		}

		// Print as a title "# Products"
		styledTitle := lipgloss.NewStyle().
			Bold(true).Foreground(lipgloss.Color("#FFFF88")).
			Background(lipgloss.Color("#5F5FFF")).
			Render("# Products")
		fmt.Println(styledTitle)

		// Lipgloss table rendering with lipgloss/table
		for _, prod := range allProducts {
			styledTitle := lipgloss.NewStyle().
				Bold(true).Foreground(lipgloss.Color("#00AFF8")).
				Render("\n## " + prod.Name + "\n")
			fmt.Println(styledTitle)

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
				fmt.Println(lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244")).Render("No release data available."))
				continue
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
			t.Border(lipgloss.RoundedBorder())
			t.BorderTop(false)
			t.BorderBottom(true)
			t.BorderLeft(false)
			t.BorderRight(false)
			t.BorderStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("63")))
			t.StyleFunc(func(row, col int) lipgloss.Style {
				padding := 1
				return lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Align(lipgloss.Left).Padding(0, padding)
			})
			renderedTable := t.Render()
			fmt.Println(renderedTable)
			// Always show a summary line below the table
			tableLines := strings.Split(renderedTable, "\n")
			maxLen := 0
			for _, l := range tableLines {
				if len(l) > maxLen {
					maxLen = len(l)
				}
			}
			summary := fmt.Sprintf("%d rows (%d shown)", len(prod.Releases), displayCount)
			fmt.Printf("%s\n", summary)
		}
	},
}
