package cmd

import (
	"fmt"
	"runtime"
	"strings"

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
	Short:   "Information about this CLI",
	Long:    `This disruptive innovation CLI (functional scope, stack, Open Source) is the result of a whole process of innovations and context: the end user should be able to learn about it...from the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Define colors using the company's hexadecimal codes
		titleColor := color.New(color.FgHiYellow).Add(color.BgBlue, color.Bold) // Yellow text on dark blue background
		sloganColor := color.New(color.FgHiWhite, color.Italic)                 // #FFFFFF (white) in italic
		sectionColor := color.New(color.FgHiYellow)                             // #FBBC09 (yellow)
		labelColor := color.New(color.FgHiWhite)                                // #FFFFFF (white) for labels
		valueColor := color.New(color.FgHiWhite)                                // #FFFFFF (white) for values

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
			if _, err := titleColor.Add(color.BgBlue).Println(paddedLine); err != nil {
				fmt.Printf("Error printing ASCII art line: %v\n", err)
			}
		}

		// Display the slogan in white and italic
		if _, err := sloganColor.Println("\"TODO\""); err != nil { //TODO
			fmt.Printf("Error printing slogan: %v\n", err)
		}

		// Display injected metadata
		fmt.Println()
		if _, err := sectionColor.Println("--- Build Info ---"); err != nil {
			fmt.Printf("Error printing build info section: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "GitVersion:"); err != nil {
			fmt.Printf("Error printing GitVersion label: %v\n", err)
		}
		if _, err := valueColor.Println(Version); err != nil {
			fmt.Printf("Error printing GitVersion value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "Git Commit:"); err != nil {
			fmt.Printf("Error printing Git Commit label: %v\n", err)
		}
		if _, err := valueColor.Println(Commit); err != nil {
			fmt.Printf("Error printing Git Commit value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "BuildDate:"); err != nil {
			fmt.Printf("Error printing BuildDate label: %v\n", err)
		}
		if _, err := valueColor.Println(Date); err != nil {
			fmt.Printf("Error printing BuildDate value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "BuiltBy:"); err != nil {
			fmt.Printf("Error printing BuiltBy label: %v\n", err)
		}
		if _, err := valueColor.Println(BuiltBy); err != nil {
			fmt.Printf("Error printing BuiltBy value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "GoVersion:"); err != nil {
			fmt.Printf("Error printing GoVersion label: %v\n", err)
		}
		if _, err := valueColor.Println(GoVersion); err != nil {
			fmt.Printf("Error printing GoVersion value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "Compiler:"); err != nil {
			fmt.Printf("Error printing Compiler label: %v\n", err)
		}
		if _, err := valueColor.Println(runtime.Compiler); err != nil {
			fmt.Printf("Error printing Compiler value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "Platform:"); err != nil {
			fmt.Printf("Error printing Platform label: %v\n", err)
		}
		if _, err := valueColor.Printf("%s/%s\n", runtime.GOOS, runtime.GOARCH); err != nil {
			fmt.Printf("Error printing Platform value: %v\n", err)
		}

		// Affichage des ressources
		fmt.Println()
		if _, err := sectionColor.Println("--- Ressources ---"); err != nil {
			fmt.Printf("Error printing resources section: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "Licence:"); err != nil {
			fmt.Printf("Error printing Licence label: %v\n", err)
		}
		if _, err := valueColor.Println("Apache-2.0"); err != nil {
			fmt.Printf("Error printing Licence value: %v\n", err)
		}

		if _, err := labelColor.Printf("%-20s ", "Code:"); err != nil {
			fmt.Printf("Error printing Code label: %v\n", err)
		}
		if _, err := valueColor.Println("https://github.com/opt-nc/geol"); err != nil {
			fmt.Printf("Error printing Code value: %v\n", err)
		}

		// if _, err := labelColor.Printf("%-20s ", "Roadmap:"); err != nil {
		// 	fmt.Printf("Error printing Roadmap label: %v\n", err)
		// }
		// if _, err := valueColor.Println("https://github.com/orgs/opt-nc/projects/24"); err != nil {
		// 	fmt.Printf("Error printing Roadmap value: %v\n", err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
