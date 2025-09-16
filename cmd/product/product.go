package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	ProductCmd.AddCommand(extendedCmd)
	ProductCmd.AddCommand(describeCmd)
	utilities.InitLogger()
}

// ProductCmd represents the product command
var ProductCmd = &cobra.Command{
	Use:     "product",
	Aliases: []string{"p"},
	Short:   "Display the latest version of one or more products and the end of life date.",
	Long:    "Show the latest version, release date, and end-of-life information for one or more products. Use the `extended` subcommand for more detailed output.",
	Example: `geol product linux ubuntu
geol product extended golang k8s
geol product describe nodejs`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Warn().Msg("Please specify at least one product.")
			return
		}

		// Load the local cache
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving cache path")
			return
		}

		// Ensure cache exists, create if missing
		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error ensuring cache exists")
			return
		}

		utilities.CheckCacheTimeAndUpdate(cmd, info.ModTime())

		products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving products from cache")
			return
		}

		// Structure to store results
		type ProductResult struct {
			Name        string
			EolLabel    string
			ReleaseName string
			ReleaseDate string
			EolFrom     string
		}
		var results []ProductResult

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

			// API request for this product
			url := utilities.ApiUrl + "products/" + prod
			resp, err := http.Get(url)
			if err != nil {
				log.Error().Err(err).Msgf("Error requesting %s", prod)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if cerr := resp.Body.Close(); cerr != nil {
				log.Error().Err(cerr).Msgf("Error closing HTTP body for %s", prod)
			}
			if err != nil {
				log.Error().Err(err).Msgf("Error reading response for %s", prod)
				continue
			}
			if resp.StatusCode != 200 {
				log.Warn().Msgf("Product %s not found on the API.", prod)
				continue
			}

			// JSON decoding
			var apiResp struct {
				Result struct {
					Name   string `json:"name"`
					Labels struct {
						Eol string `json:"eol"`
					} `json:"labels"`
					Releases []struct {
						Name        string `json:"name"`
						ReleaseDate string `json:"releaseDate"`
						EolFrom     string `json:"eolFrom"`
					} `json:"releases"`
				} `json:"result"`
			}
			if err := json.Unmarshal(body, &apiResp); err != nil {
				log.Error().Err(err).Msgf("Error decoding JSON for %s", prod)
				continue
			}
			var relName, relDate, relEol string
			if len(apiResp.Result.Releases) > 0 {
				relName = apiResp.Result.Releases[0].Name
				relDate = apiResp.Result.Releases[0].ReleaseDate
				relEol = apiResp.Result.Releases[0].EolFrom
			}
			results = append(results, ProductResult{
				Name: apiResp.Result.Name,
				//EolLabel:    apiResp.Result.Labels.Eol,
				ReleaseName: relName,
				ReleaseDate: relDate,
				EolFrom:     relEol,
			})
		}

		// Display markdown table with glamour
		if len(results) == 0 {
			log.Warn().Msg("No product found in cache or API.")
			return
		}
		var buf bytes.Buffer
		// Header
		buf.WriteString("| **Name** | **Latest Cycle** | **Release Date** | **EOL From** |\n")
		buf.WriteString("|------|--------------|--------------|----------|\n")
		for _, r := range results {
			name := strings.ReplaceAll(r.Name, "|", "\\|")
			relName := strings.ReplaceAll(r.ReleaseName, "|", "\\|")
			relDate := strings.ReplaceAll(r.ReleaseDate, "|", "\\|")
			relEol := strings.ReplaceAll(r.EolFrom, "|", "\\|")
			buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", name, relName, relDate, relEol))
		}

		out, err := glamour.RenderWithEnvironmentConfig(buf.String())
		if err != nil {
			fmt.Print(buf.String()) // raw fallback
		} else {
			fmt.Print(out)
		}
	},
}
