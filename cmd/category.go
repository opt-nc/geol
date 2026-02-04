package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/opt-nc/geol/v2/utilities"
	"github.com/phuslu/log"
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
			log.Fatal().Msg("Please specify a category")
		}
		if len(args) > 1 {
			log.Fatal().Msg("Please specify only one category")
		}
		category := args[0]

		utilities.AnalyzeCacheProductsValidity(cmd)
		categoriesPath, err := utilities.GetCategoriesPath()
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving categories path")
		}
		categories, err := utilities.GetCategoriesWithCacheRefresh(cmd, categoriesPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving categories from cache")
		}
		if _, ok := categories[category]; !ok {
			log.Fatal().Msgf("Category '%s' not found in cache", category)
		}

		url := utilities.ApiUrl + "categories/" + category
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error requesting category '%s'", category)
		}
		body, err := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if err != nil {
			log.Fatal().Err(err).Msgf("Error reading response for category '%s'", category)
		}
		if closeErr != nil {
			log.Fatal().Err(closeErr).Msgf("Error closing response body for category '%s'", category)
		}
		if resp.StatusCode != 200 {
			log.Fatal().Msgf("Category '%s' not found on the API", category)
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
			log.Fatal().Err(err).Msgf("Error decoding JSON for category '%s'", category)
		}

		treeRoot := tree.Root(".")
		categoryNode := tree.New().Root(boldStyle.Render(category))
		for _, prod := range apiResp.Result {
			categoryNode.Child(boldStyle.Render(prod.Name))
		}
		treeRoot.Child(categoryNode)
		if _, err := lipgloss.Println(treeRoot.String()); err != nil {
			log.Fatal().Err(err).Msgf("Error printing tree for category '%s'", category)
		}
		nbProducts := len(apiResp.Result)
		nbProductsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := lipgloss.Printf("\n%s products listed for category '%s'\n", nbProductsStyle.Render(fmt.Sprintf("%d", nbProducts)), category); err != nil {
			log.Fatal().Err(err).Msgf("Error printing product count for category '%s'", category)
		}
	},
}

func init() {
	rootCmd.AddCommand(categoryCmd)
}
