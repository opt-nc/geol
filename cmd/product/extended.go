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
	"github.com/spf13/cobra"
)

// extendedCmd represents the extended command
var extendedCmd = &cobra.Command{
	Use:   "extended",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		numberFlag, _ := cmd.Flags().GetInt("number")

		if numberFlag < 0 {
			fmt.Println("Le nombre de lignes doit être positif ou nul.")
			return
		}

		if len(args) == 0 {
			fmt.Println("Please specify at least one product.")
			return
		}

		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			fmt.Println("Error retrieving cache path:", err)
			return
		}

		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			fmt.Println("Error ensuring cache exists:", err)
			return
		}

		utilities.CheckCacheTimeAndUpdate(cmd, info.ModTime())

		cacheFile, err := os.Open(productsPath)
		if err != nil {
			fmt.Println("Error opening cache:", err)
			return
		}
		defer func() {
			if err := cacheFile.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing cache: %v\n", err)
			}
		}()
		var products utilities.ProductsFile
		if err := json.NewDecoder(cacheFile).Decode(&products); err != nil {
			fmt.Println("Error decoding cache:", err)
			return
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
				cmd.Printf("Error requesting %s: %v\n", prod, err)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if cerr := resp.Body.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "Error closing HTTP body for %s: %v\n", prod, cerr)
			}
			if err != nil {
				fmt.Printf("Error reading response for %s: %v\n", prod, err)
				continue
			}
			if resp.StatusCode != 200 {
				fmt.Printf("Product %s not found on the API.\n", prod)
				continue
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
					} `json:"releases"`
				}
			}
			if err := json.Unmarshal(body, &apiResp); err != nil {
				fmt.Printf("Error decoding JSON for %s: %v\n", prod, err)
				continue
			}

			var releases []struct {
				Name        string
				ReleaseDate string
				LatestName  string
				LatestDate  string
				EoasFrom    string
				EolFrom     string
			}
			for _, r := range apiResp.Result.Releases {
				releases = append(releases, struct {
					Name        string
					ReleaseDate string
					LatestName  string
					LatestDate  string
					EoasFrom    string
					EolFrom     string
				}{
					Name:        r.Name,
					ReleaseDate: r.ReleaseDate,
					LatestName:  r.Latest.Name,
					LatestDate:  r.Latest.Date,
					EoasFrom:    r.EoasFrom,
					EolFrom:     r.EolFrom,
				})
			}
			allProducts = append(allProducts, ProductReleases{
				Name:     apiResp.Result.Name,
				Releases: releases,
			})
		}

		if len(allProducts) == 0 {
			fmt.Println("Aucun produit trouvé dans le cache ou l'API.")
			return
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
				if showName {
					row = append(row, r.Name)
				}
				if showReleaseDate {
					row = append(row, r.ReleaseDate)
				}
				if showLatestName {
					row = append(row, r.LatestName)
				}
				if showLatestDate {
					row = append(row, r.LatestDate)
				}
				if showEoasFrom {
					row = append(row, r.EoasFrom)
				}
				if showEolFrom {
					row = append(row, r.EolFrom)
				}
				t.Row(row...)
			}
			t.Border(lipgloss.NormalBorder())
			t.BorderTop(false)
			t.BorderBottom(false)
			t.BorderLeft(false)
			t.BorderRight(false)
			t.BorderStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("63")))
			t.StyleFunc(func(row, col int) lipgloss.Style {
				padding := 1
				return lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Align(lipgloss.Center).Padding(0, padding)
			})
			fmt.Println(t.Render())
			fmt.Println()
		}
	},
}

func init() {
	extendedCmd.Flags().IntP("number", "n", 10, "Number of rows to display (0 to display all)")
}
