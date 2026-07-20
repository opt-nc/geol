/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package ci_github

import (
	"github.com/opt-nc/geol/v2/cmd/templates"
	"github.com/spf13/cobra"
)

var output string
var force bool

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Generate a ready to use GitHub Actions workflow file for ci-github command",
	Long: `The init command generates a GitHub Actions workflow file for the ci-github command.
Use this command to create a starter workflow that you can customize for your environment.
You can specify the output path with the --output flag.
Use --force to overwrite an existing file.`,
	Example: `geol ci-github init
geol ci-github init --output .github/workflows/geol-check.yml
geol ci-github init --output .github/workflows/geol-check.yml --force`,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")
		templates.GenerateCGithubTemplate(output, force)
	},
}

func init() {
	InitCmd.Flags().StringVarP(&output, "output", "o", ".github/workflows/geol-action.yml", "Path to the output file")
	InitCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite the file if it already exists")
}
