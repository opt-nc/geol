/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/opt-nc/geol/cmd/templates"
	"github.com/spf13/cobra"
)

var output string

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"t"},
	Short:   "Generate a template configuration file for check command",
	Long: `The template command generates a default configuration file for the check command.
Use this command to create a starter YAML file that you can customize for your environment.
You can specify the output path with the --output flag.`,
	Example: `geol check template
geol check template --output stack.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		templates.GenerateTemplate(output)
	},
}

func init() {
	checkCmd.AddCommand(templateCmd)
	templateCmd.Flags().StringVarP(&output, "output", "o", ".geol.yaml", "Path to the output file")
}
