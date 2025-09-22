package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/phuslu/log"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Variables injected at build time with ldflags
var (
	Commit    = "none"
	Date      = "unknown"
	BuiltBy   = "GoReleaser"
	GoVersion = runtime.Version()
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:     "about",
	Aliases: []string{"a"},
	Example: `geol about`,
	Short:   "Information about geol",
	Long:    `This disruptive innovation CLI (functional scope, stack, Open Source) is the result of a whole process of innovations and context: the end user should be able to learn about it...from the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Define colors using the company's hexadecimal codes
		titleColor := color.New(color.FgHiYellow)               // Yellow text
		sloganColor := color.New(color.FgHiWhite, color.Italic) // #FFFFFF (white) in italic
		sectionColor := color.New(color.FgHiYellow)             // #FBBC09 (yellow)
		labelColor := color.New(color.FgHiWhite)                // #FFFFFF (white) for labels
		valueColor := color.New(color.FgHiWhite)                // #FFFFFF (white) for values

		// Generate ASCII art in yellow on blue background
		myFigure := figure.NewFigure("geol", "starwars", true)
		asciiArtLines := myFigure.Slicify()
		maxLength := 0
		for _, line := range asciiArtLines {
			if len(line) > maxLength {
				maxLength = len(line)
			}
		}
		for _, line := range asciiArtLines {
			paddedLine := line + strings.Repeat(" ", maxLength-len(line))
			if _, err := titleColor.Println(paddedLine); err != nil {
				log.Error().Err(err).Msg("Error printing ASCII art line")
			}
		}

		// Display the slogan in white and italic
		if _, err := sloganColor.Println("⏳ Tech doesn’t last forever. Awareness does."); err != nil {
			log.Error().Err(err).Msg("Error printing slogan")
		}

		// Display injected metadata
		fmt.Println()
		if _, err := sectionColor.Println("--- Build Info ---"); err != nil {
			log.Error().Err(err).Msg("Error printing build info section")
		}

		if _, err := labelColor.Printf("%-20s ", "GitVersion:"); err != nil {
			log.Error().Err(err).Msg("Error printing GitVersion label")
		}
		if _, err := valueColor.Println(Version); err != nil {
			log.Error().Err(err).Msg("Error printing GitVersion value")
		}

		if _, err := labelColor.Printf("%-20s ", "Git Commit:"); err != nil {
			log.Error().Err(err).Msg("Error printing Git Commit label")
		}
		if _, err := valueColor.Println(Commit); err != nil {
			log.Error().Err(err).Msg("Error printing Git Commit value")
		}

		if _, err := labelColor.Printf("%-20s ", "BuildDate:"); err != nil {
			log.Error().Err(err).Msg("Error printing BuildDate label")
		}
		if _, err := valueColor.Println(Date); err != nil {
			log.Error().Err(err).Msg("Error printing BuildDate value")
		}

		if _, err := labelColor.Printf("%-20s ", "BuiltBy:"); err != nil {
			log.Error().Err(err).Msg("Error printing BuiltBy label")
		}
		if _, err := valueColor.Println(BuiltBy); err != nil {
			log.Error().Err(err).Msg("Error printing BuiltBy value")
		}

		if _, err := labelColor.Printf("%-20s ", "GoVersion:"); err != nil {
			log.Error().Err(err).Msg("Error printing GoVersion label")
		}
		if _, err := valueColor.Println(GoVersion); err != nil {
			log.Error().Err(err).Msg("Error printing GoVersion value")
		}

		if _, err := labelColor.Printf("%-20s ", "Compiler:"); err != nil {
			log.Error().Err(err).Msg("Error printing Compiler label")
		}
		if _, err := valueColor.Println(runtime.Compiler); err != nil {
			log.Error().Err(err).Msg("Error printing Compiler value")
		}

		if _, err := labelColor.Printf("%-20s ", "Platform:"); err != nil {
			log.Error().Err(err).Msg("Error printing Platform label")
		}
		if _, err := valueColor.Printf("%s/%s\n", runtime.GOOS, runtime.GOARCH); err != nil {
			log.Error().Err(err).Msg("Error printing Platform value")
		}

		// Affichage des ressources
		fmt.Println()
		if _, err := sectionColor.Println("--- Ressources ---"); err != nil {
			log.Error().Err(err).Msg("Error printing resources section")
		}

		if _, err := labelColor.Printf("%-20s ", "Licence:"); err != nil {
			log.Error().Err(err).Msg("Error printing Licence label")
		}
		if _, err := valueColor.Println("Apache-2.0"); err != nil {
			log.Error().Err(err).Msg("Error printing Licence value")
		}

		if _, err := labelColor.Printf("%-20s ", "Code:"); err != nil {
			log.Error().Err(err).Msg("Error printing Code label")
		}
		if _, err := valueColor.Println("https://github.com/opt-nc/geol"); err != nil {
			log.Error().Err(err).Msg("Error printing Code value")
		}

		if _, err := labelColor.Printf("%-20s ", "Roadmap:"); err != nil {
			log.Error().Err(err).Msg("Error printing Roadmap label")
		}
		if _, err := valueColor.Println("https://github.com/orgs/opt-nc/projects/28"); err != nil {
			log.Error().Err(err).Msg("Error printing Roadmap value")
		}

		if _, err := labelColor.Printf("%-20s ", "API:"); err != nil {
			log.Error().Err(err).Msg("Error printing API label")
		}
		if _, err := valueColor.Println("https://endoflife.date"); err != nil {
			log.Error().Err(err).Msg("Error printing API value")
		}
	},
}

func init() {
	// Set up pretty console writer for phuslu/log
	log.DefaultLogger.Writer = &log.ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    true,
		EndWithMessage: true,
	}
	rootCmd.AddCommand(aboutCmd)
}
