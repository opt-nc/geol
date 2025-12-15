package exports

import (
	"github.com/spf13/cobra"
)

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data to various formats",
	Long: `Export the application data and product information to different file formats.
By default, this command exports to a DuckDB database (equivalent to 'export duckdb').

Available formats:
- duckdb: Export to a DuckDB database file (default)`,
	Run: func(cmd *cobra.Command, args []string) {
		duckdbCmd.Run(duckdbCmd, args)
	},
}

func init() {

}
