package items

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

// TagsCmd represents the tags command
var TagsCmd = &cobra.Command{
	Use:     "tags",
	Aliases: []string{"t"},
	Short:   "List all cached tag names.",
	Long:    `Displays the list of all tag names currently available in the cache.`,
	Example: `geol list tags\ngeol l t`,
	Run: func(cmd *cobra.Command, args []string) {
		// List the cached tags
		utilities.AnalyzeCacheTagsValidity(cmd)
		tagsPath, err := utilities.GetTagsPath()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving tags path")
			os.Exit(1)
		}

		tags, err := utilities.GetTagsWithCacheRefresh(cmd, tagsPath)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving tags from cache")
			os.Exit(1)
		}

		var names []string
		for name := range tags {
			names = append(names, name)
		}
		sort.Strings(names)
		tagColor := color.New(color.Bold)

		plusStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render("+")
		for _, name := range names {
			if _, err := fmt.Fprintf(os.Stdout, "%s %s\n", plusStyle, tagColor.Sprint(name)); err != nil {
				log.Error().Err(err).Msg("Error writing tag name to stdout")
				os.Exit(1)
			}
		}

		nbTags := strconv.Itoa(len(names))
		nbTagsStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
		if _, err := fmt.Fprintf(os.Stdout, "\n%s tags listed\n", nbTagsStyle.Render(nbTags)); err != nil {
			log.Error().Err(err).Msg("Error writing tag count to stdout")
			os.Exit(1)
		}
	},
}

func init() {

}
