package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

// ProductCmd represents the product command
var ProductCmd = &cobra.Command{
	Use:     "product",
	Aliases: []string{"p"},
	Short:   "Display the latest version of one or more products and the end of life date.",
	Long:    "Show the latest version, release date, and end-of-life information for one or more products. Use the `extended` subcommand for more detailed output.",
	Example: `geol product linux ubuntu
geol product extended golang k8s`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify at least one product.")
			return
		}

		// Load the local cache
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			fmt.Println("Error retrieving cache path:", err)
			return
		}

		// Ensure cache exists, create if missing
		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			fmt.Println("Error ensuring cache exists:", err)
			return
		}

		utilities.CheckCacheTimeAndUpdate(cmd, info.ModTime())

		cacheFile, err := os.Open(productsPath)
		if err != nil {
			fmt.Println("Error opening local cache:", err)
			return
		}

		defer func() {
			if err := cacheFile.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing local cache: %v\n", err)
			}
		}()

		var products utilities.ProductsFile
		if err := json.NewDecoder(cacheFile).Decode(&products); err != nil {
			fmt.Println("Error decoding cache:", err)
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
				fmt.Printf("Error requesting %s: %v\n", prod, err)
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
				fmt.Printf("Error decoding JSON for %s: %v\n", prod, err)
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
			fmt.Println("No product found in cache or API.")
			return
		}
		var buf bytes.Buffer
		// Header
		buf.WriteString("| **Name** | **Version** | **Release Date** | **EOL From** |\n")
		buf.WriteString("|------|--------------|--------------|----------|\n")
		for _, r := range results {
			name := strings.ReplaceAll(r.Name, "|", "\\|")
			//eolLabel := strings.ReplaceAll(r.EolLabel, "|", "\\|")
			relName := strings.ReplaceAll(r.ReleaseName, "|", "\\|")
			relDate := strings.ReplaceAll(r.ReleaseDate, "|", "\\|")
			relEol := strings.ReplaceAll(r.EolFrom, "|", "\\|")
			buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", name, relName, relDate, relEol))
		}
		// r, _ := glamour.NewTermRenderer(
		//      glamour.WithAutoStyle(),
		//      //glamour.WithWordWrap(120),
		// )

		//out, err := glamour.Render(buf.String(), "dark")
		//out, err := r.Render(buf.String())
		out, err := glamour.RenderWithEnvironmentConfig(buf.String())
		if err != nil {
			fmt.Print(buf.String()) // raw fallback
		} else {
			fmt.Print(out)
		}
	},
}

func init() {
	ProductCmd.AddCommand(extendedCmd)
}
