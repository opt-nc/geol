package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

var Version string = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"}, // Alias for the command
	Short:   "Display the application version",
	Long:    `Display the application version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
		log.Info().Msg("Checking the latest geol version...")
		latestVersion := getLatestVersionFromGitHub()
		if latestVersion == "" {
			log.Warn().Msg("Could not check the latest version from GitHub")
			return
		}
		if latestVersion != Version {
			log.Warn().Msg("There is a new version available ! Latest version: " + latestVersion + ", you have: " + Version)
		} else {
			log.Info().Msg("You have the latest version !")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// getLatestVersionFromGitHub fetches the latest release tag from GitHub
func getLatestVersionFromGitHub() string {
	resp, err := http.Get("https://api.github.com/repos/opt-nc/geol/releases/latest")
	if err != nil {
		return ""
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warn().Err(err).Msg("Error closing response body in getLatestVersionFromGitHub")
		}
	}()

	var result struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}
	tag := result.TagName
	if strings.HasPrefix(tag, "v") {
		return tag[1:]
	}
	return tag
}
