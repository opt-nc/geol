package exports

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/utilities"
)

// duckdbCmd represents the duckdb command
var duckdbCmd = &cobra.Command{
	Use:   "duckdb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Open or create geol.duckdb at project root
		db, err := sql.Open("duckdb", "geol.duckdb")
		if err != nil {
			log.Fatalf("Error while creating DuckDB database: %v", err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Fatalf("Error closing DuckDB database: %v", err)
			}
		}()

		// Create 'about' table if not exists
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS about (
				gitVersion TEXT,
				gitCommit TEXT,
				goVersion TEXT,
				platform TEXT,
				githubURL TEXT,
				generatedAt TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::TIMESTAMP,
				generatedAtTZ TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
			)`)
		if err != nil {
			log.Fatalf("Error creating 'about' table: %v", err)
		}

		// Insert values into 'about' table
		_, err = db.Exec(`INSERT INTO about (gitVersion, gitCommit, goVersion, platform, githubURL) 
			VALUES (?, ?, ?, ?, ?)`,
			utilities.Version, utilities.Commit, utilities.GoVersion,
			fmt.Sprintf("%s/%s", utilities.PlatformOs, utilities.PlatformArch),
			"https://github.com/opt-nc/geol")
		if err != nil {
			log.Fatalf("Error inserting into 'about' table: %v", err)
		}

		fmt.Println("Table 'about' created and metadata inserted into geol.duckdb.")
	},
}

func init() {
	ExportCmd.AddCommand(duckdbCmd)

}
