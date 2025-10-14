/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/fatih/color"
	"github.com/opt-nc/geol/utilities"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Display all products associated with a tag.",
	Long:  `Show all products associated with a given tag. The tag must exist in the cache. Results are displayed in a tree structure.`,
	Example: `geol tag os
geol tag canonical`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a tag.")
			os.Exit(1)
		}
		if len(args) > 1 {
			fmt.Println("Please specify only one tag.")
			os.Exit(1)
		}
		tag := args[0]

		// Vérifier le cache des tags
		utilities.AnalyzeCacheProductsValidity(cmd)
		tagsPath, err := utilities.GetTagsPath()
		if err != nil {
			fmt.Println("Error retrieving tags path:", err)
			return
		}
		tags, err := utilities.GetTagsWithCacheRefresh(cmd, tagsPath)
		if err != nil {
			fmt.Println("Error retrieving tags from cache:", err)
			return
		}
		if _, ok := tags[tag]; !ok {
			fmt.Printf("Tag '%s' not found in cache.\n", tag)
			return
		}

		url := utilities.ApiUrl + "tags/" + tag
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error requesting tag '%s': %v\n", tag, err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response for tag '%s': %v\n", tag, err)
			return
		}
		if resp.StatusCode != 200 {
			fmt.Printf("Tag '%s' not found on the API.\n", tag)
			return
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
			fmt.Printf("Error decoding JSON for tag '%s': %v\n", tag, err)
			return
		}

		tagColor := color.New(color.Bold)
		productColor := color.New(color.Bold)
		treeRoot := tree.Root(".")
		tagNode := tree.New().Root(tagColor.Sprint(tag))
		for _, prod := range apiResp.Result {
			tagNode.Child(productColor.Sprint(prod.Name))
		}
		treeRoot.Child(tagNode)
		fmt.Fprintln(os.Stdout, treeRoot.String())
		nbProducts := len(apiResp.Result)
		nbProductsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		fmt.Fprintf(os.Stdout, "\n%s products listed for tag '%s'\n", nbProductsStyle.Render(fmt.Sprintf("%d", nbProducts)), tag)
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

}
