package ci_github

import (
	"github.com/spf13/cobra"
)

// CiGithubCmd represents the ci-github command
var CiGithubCmd = &cobra.Command{
	Use:     "ci-github",
	Aliases: []string{"ci-gh"},
	Short:   "Manage GitHub Actions CI configuration",
	Long: `Manage GitHub Actions CI configuration for Geol.
By default, this command generates a GitHub Actions workflow (equivalent to 'ci-github init').

Available subcommands:
- init: Generate a ready to use GitHub Actions workflow file (default)`,
	Run: func(cmd *cobra.Command, args []string) {
		InitCmd.Run(cmd, args)
	},
}

func init() {
	CiGithubCmd.AddCommand(InitCmd)
	CiGithubCmd.Flags().StringP("output", "o", ".github/workflows/geol-action.yml", "Path to the output file")
	CiGithubCmd.Flags().BoolP("force", "f", false, "Overwrite the file if it already exists")
}
