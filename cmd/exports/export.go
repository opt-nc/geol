package exports

import (
	"github.com/spf13/cobra"
)

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data to various formats",
	Long: `Export all known products and their end-of-life (EOL) metadata in different file formats.
By default, this command exports to a DuckDB database (equivalent to 'export duckdb').

Available formats:
- duckdb: Export to a DuckDB database file (default)`,
	Run: func(cmd *cobra.Command, args []string) {
		duckdbCmd.Run(cmd, args)
	},
}

func init() {
	ExportCmd.Flags().StringP("output", "o", "geol.duckdb", "Output DuckDB database file path")
	ExportCmd.Flags().BoolP("force", "f", false, "Overwrites the DuckDB database file if it already exists")
}
