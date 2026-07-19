/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package check

import (
	"github.com/opt-nc/geol/v2/cmd/templates"
	"github.com/spf13/cobra"
)

var output string
var force bool
var appName string
var appID string

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Generate a valid template configuration file for check command",
	Long: `The init command generates a default configuration file for the check command.
Use this command to create a starter YAML file that you can customize for your environment.
You can specify the output path with the --output flag.
Use --force to overwrite an existing file.`,
	Example: `geol check init
geol check init --output stack.yaml
geol check init --output stack.yaml --force`,
	Run: func(cmd *cobra.Command, args []string) {
		templates.GenerateTemplate(output, force, appName, appID)
	},
}

func init() {
	InitCmd.Flags().StringVarP(&output, "output", "o", ".geol.yaml", "Path to the output file")
	InitCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite the file if it already exists")
	InitCmd.Flags().StringVarP(&appName, "app-name", "a", "", "Application name to use in the generated template")
	InitCmd.Flags().StringVar(&appID, "app-id", "", "Application ID to use in the generated template")
}
