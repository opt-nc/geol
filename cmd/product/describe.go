package product

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

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
			cmd.Println("Please specify exactly one product.")
			return
		}
		prodArg := args[0]

		// Check the cache
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			cmd.PrintErrln("Error retrieving cache path:", err)
			return
		}
		info, err := utilities.EnsureCacheExists(cmd, productsPath)
		if err != nil {
			cmd.PrintErrln("Error ensuring cache exists:", err)
			return
		}
		utilities.CheckCacheTimeAndUpdate(cmd, info.ModTime())

		cacheFile, err := os.Open(productsPath)
		if err != nil {
			cmd.PrintErrln("Error opening local cache:", err)
			return
		}
		defer cacheFile.Close()

		var products utilities.ProductsFile
		if err := json.NewDecoder(cacheFile).Decode(&products); err != nil {
			cmd.PrintErrln("Error decoding cache:", err)
			return
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
			cmd.Printf("Product '%s' not found in cache.\n", prodArg)
			return
		}

		// Build the markdown URL
		mdUrl := "https://raw.githubusercontent.com/endoflife-date/endoflife.date/refs/heads/master/products/" + mainName + ".md"

		// Retrieve the Markdown content
		resp, err := http.Get(mdUrl)
		if err != nil {
			cmd.PrintErrln("Error fetching markdown:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			cmd.Printf("Failed to fetch markdown. Status: %s\n", resp.Status)
			return
		}

		mdBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			cmd.PrintErrln("Error reading markdown:", err)
			return
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
			cmd.Println("No description found in markdown.")
			return
		}

		// Print a product title as in extended: # ProductName, with color and background
		styledTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFF88")).
			Background(lipgloss.Color("#5F5FFF")).
			Render("# `" + mainName + "`")
		os.Stdout.Write([]byte(styledTitle))

		// Glamour rendering only on the description
		out, err := glamour.RenderWithEnvironmentConfig(desc)
		if err != nil {
			cmd.PrintErrln("Error rendering markdown:", err)
			return
		}
		os.Stdout.Write([]byte(out))

	},
}

func init() {
}
