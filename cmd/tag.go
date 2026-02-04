/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

var boldStyle = lipgloss.NewStyle().Bold(true)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:     "tag",
	Aliases: []string{"t"},
	Short:   "Display all products associated with a tag.",
	Long:    `Show all products associated with a given tag. The tag must exist in the cache. Results are displayed in a tree structure.`,
	Example: `geol tag os
geol tag canonical`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal().Msg("Please specify a tag")
		}
		if len(args) > 1 {
			log.Fatal().Msg("Please specify only one tag")
		}
		tag := args[0]

		utilities.AnalyzeCacheProductsValidity(cmd)
		tagsPath, err := utilities.GetTagsPath()
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving tags path")
		}
		tags, err := utilities.GetTagsWithCacheRefresh(cmd, tagsPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving tags from cache")
		}
		if _, ok := tags[tag]; !ok {
			log.Fatal().Msgf("Tag '%s' not found in cache", tag)
		}

		url := utilities.ApiUrl + "tags/" + tag
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error requesting tag '%s'", tag)
		}
		body, err := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if err != nil {
			log.Fatal().Err(err).Msgf("Error reading response for tag '%s'", tag)
		}
		if closeErr != nil {
			log.Fatal().Err(closeErr).Msgf("Error closing response body for tag '%s'", tag)
		}
		if resp.StatusCode != 200 {
			log.Fatal().Msgf("Tag '%s' not found on the API", tag)
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
			log.Fatal().Err(err).Msgf("Error decoding JSON for tag '%s'", tag)
		}

		treeRoot := tree.Root(".")
		tagNode := tree.New().Root(boldStyle.Render(tag))
		for _, prod := range apiResp.Result {
			tagNode.Child(boldStyle.Render(prod.Name))
		}
		treeRoot.Child(tagNode)
		if _, err := lipgloss.Println(treeRoot.String()); err != nil {
			log.Fatal().Err(err).Msgf("Error printing tree for tag '%s'", tag)
		}
		nbProducts := len(apiResp.Result)
		nbProductsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := lipgloss.Printf("\n%s products listed for tag '%s'\n", nbProductsStyle.Render(fmt.Sprintf("%d", nbProducts)), tag); err != nil {
			log.Fatal().Err(err).Msgf("Error printing product count for tag '%s'", tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
