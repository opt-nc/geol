package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/fatih/color"
	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

// categoryCmd represents the category command
var categoryCmd = &cobra.Command{
	Use:     "category",
	Aliases: []string{"cat"},
	Short:   "Display all products associated with a category.",
	Long:    `Show all products associated with a given category. The category must exist in the cache. Results are displayed in a tree structure.`,
	Example: `geol category os
geol category cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a category.")
			os.Exit(1)
		}
		if len(args) > 1 {
			fmt.Println("Please specify only one category.")
			os.Exit(1)
		}
		category := args[0]

		utilities.AnalyzeCacheProductsValidity(cmd)
		categoriesPath, err := utilities.GetCategoriesPath()
		if err != nil {
			fmt.Println("Error retrieving categories path:", err)
			os.Exit(1)
		}
		categories, err := utilities.GetCategoriesWithCacheRefresh(cmd, categoriesPath)
		if err != nil {
			fmt.Println("Error retrieving categories from cache:", err)
			os.Exit(1)
		}
		if _, ok := categories[category]; !ok {
			fmt.Printf("Category '%s' not found in cache.\n", category)
			os.Exit(1)
		}

		url := utilities.ApiUrl + "categories/" + category
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error requesting category '%s': %v\n", category, err)
			os.Exit(1)
		}
		body, err := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response for category '%s': %v\n", category, err)
			os.Exit(1)
		}
		if closeErr != nil {
			fmt.Printf("Error closing response body for category '%s': %v\n", category, closeErr)
			os.Exit(1)
		}
		if resp.StatusCode != 200 {
			fmt.Printf("Category '%s' not found on the API.\n", category)
			os.Exit(1)
		}

		var apiResp struct {
			Result []struct {
				Name     string   `json:"name"`
				Label    string   `json:"label"`
				Aliases  []string `json:"aliases"`
				Category string   `json:"category"`
				Tags     []string `json:"tags"`
				Uri      string   `json:"uri"`
			} `json:"result"`
		}
		if err := json.Unmarshal(body, &apiResp); err != nil {
			fmt.Printf("Error decoding JSON for category '%s': %v\n", category, err)
			os.Exit(1)
		}

		categoryColor := color.New(color.Bold)
		productColor := color.New(color.Bold)
		treeRoot := tree.Root(".")
		categoryNode := tree.New().Root(categoryColor.Sprint(category))
		for _, prod := range apiResp.Result {
			categoryNode.Child(productColor.Sprint(prod.Name))
		}
		treeRoot.Child(categoryNode)
		if _, err := fmt.Fprintln(os.Stdout, treeRoot.String()); err != nil {
			fmt.Printf("Error printing tree for category '%s': %v\n", category, err)
			os.Exit(1)
		}
		nbProducts := len(apiResp.Result)
		nbProductsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := fmt.Fprintf(os.Stdout, "\n%s products listed for category '%s'\n", nbProductsStyle.Render(fmt.Sprintf("%d", nbProducts)), category); err != nil {
			fmt.Printf("Error printing product count for category '%s': %v\n", category, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(categoryCmd)
}
