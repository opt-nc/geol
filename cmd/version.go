package cmd

import (
	"fmt"

	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

// var Version string = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"}, // Alias for the command
	Short:   "Display the application version",
	Long:    `Display the application version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utilities.Version)
		log.Info().Msg("Checking the latest geol version...")
		latestVersion := utilities.GetLatestVersionFromGitHub()
		if latestVersion == "" {
			log.Warn().Msg("Could not check the latest version from GitHub")
			return
		}
		if latestVersion != utilities.Version {
			log.Warn().Msg("There is a new geol version available ! Latest version: " + latestVersion + ", you have: " + utilities.Version)
		} else {
			log.Info().Msg("You have the latest geol version !")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
