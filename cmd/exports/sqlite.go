package exports

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/v2/utilities"
)

// sqliteCmd represents the sqlite command
var sqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Export data to a SQLite database",
	Long: `Export all known products and their end-of-life (EOL) metadata into a SQLite database file.
This command creates a new SQLite file (default: geol.sqlite) by first generating a DuckDB
database and then converting it to SQLite format, ensuring compatibility with SQLite tools.

You can specify the output filename using the --output flag.
If the file already exists, use the --force flag to overwrite it.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		utilities.AnalyzeCacheProductsValidity(cmd)
		sqlitePath, _ := cmd.Flags().GetString("output")
		forceSQLite, _ := cmd.Flags().GetBool("force")

		if _, err := os.Stat(sqlitePath); err == nil {
			if !forceSQLite {
				log.Fatal().Msgf("File %s already exists. Use --force to overwrite.", sqlitePath)
			}
			if err := os.Remove(sqlitePath); err != nil {
				log.Fatal().Err(err).Msgf("Error removing existing %s file", sqlitePath)
			}
		}

		// Create temporary DuckDB file in user config directory
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting user config directory")
		}
		geolConfigDir := filepath.Join(configDir, "geol")
		if err := os.MkdirAll(geolConfigDir, 0o755); err != nil {
			log.Fatal().Err(err).Msg("Error creating geol config directory")
		}
		tempDuckDB := filepath.Join(geolConfigDir, "geol_export.duckdb.tmp")

		// Populate DuckDB then export to SQLite, cleaning up DuckDB resources before FK step
		if err := buildDuckDBAndExportToSQLite(cmd, tempDuckDB, sqlitePath); err != nil {
			log.Fatal().Err(err).Msg("Error building DuckDB and exporting to SQLite")
		}

		// Clean up temporary DuckDB files
		for _, f := range []string{tempDuckDB, tempDuckDB + ".wal"} {
			if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
				log.Warn().Err(err).Msgf("Error removing temporary file %s", f)
			}
		}

		log.Info().Msg("Adding foreign key constraints to SQLite...")

		// Add foreign key constraints to SQLite (mirrors add_fk.sql)
		if err := addSQLiteForeignKeys(sqlitePath); err != nil {
			log.Fatal().Err(err).Msg("Error adding foreign key constraints to SQLite")
		}

		duration := time.Since(startTime)
		log.Info().Msgf("SQLite database created successfully at %s (took %v)", sqlitePath, duration.Round(time.Millisecond))
		log.Info().Msg("You can query the database using SQLite CLI or any compatible client.")
		log.Info().Msgf("Example CLI command: sqlite3 %s", sqlitePath)
	},
}

// buildDuckDBAndExportToSQLite creates a temporary DuckDB, populates it, and exports to SQLite.
// The DuckDB connection is closed before returning so the caller can safely open the SQLite file.
func buildDuckDBAndExportToSQLite(cmd *cobra.Command, tempDuckDB, sqlitePath string) error {
	log.Info().Msg("Creating temporary DuckDB database...")

	db, err := sql.Open("duckdb", tempDuckDB)
	if err != nil {
		return fmt.Errorf("error creating temporary DuckDB database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Warn().Err(err).Msg("Error closing DuckDB database")
		}
	}()

	if err := populateDuckDB(cmd, db); err != nil {
		return err
	}

	log.Info().Msg("Converting DuckDB to SQLite...")

	// Export to SQLite using the DuckDB sqlite extension (mirrors export_complete.sql)
	return exportDuckDBToSQLite(db, sqlitePath)
}

// exportDuckDBToSQLite exports the DuckDB database to SQLite using the DuckDB sqlite extension.
// Tables are created with proper PRIMARY KEY and UNIQUE constraints (mirrors export_complete.sql).
func exportDuckDBToSQLite(db *sql.DB, sqlitePath string) error {
	if _, err := db.Exec("INSTALL sqlite; LOAD sqlite;"); err != nil {
		return fmt.Errorf("error loading SQLite extension: %w", err)
	}

	if _, err := db.Exec(fmt.Sprintf("ATTACH '%s' AS sqlite_db (TYPE SQLITE)", sqlitePath)); err != nil {
		return fmt.Errorf("error attaching SQLite database: %w", err)
	}

	// Create tables with PRIMARY KEYS and UNIQUE constraints
	tableCreations := []string{
		`CREATE TABLE sqlite_db.categories(
			id TEXT PRIMARY KEY,
			uri TEXT
		)`,
		`CREATE TABLE sqlite_db.tags(
			id TEXT PRIMARY KEY,
			uri TEXT UNIQUE,
			www TEXT
		)`,
		`CREATE TABLE sqlite_db.about(
			git_version TEXT,
			git_commit TEXT,
			go_version TEXT,
			platform TEXT,
			github_url TEXT,
			generated_at TEXT,
			generated_at_tz TEXT
		)`,
		`CREATE TABLE sqlite_db.products(
			id TEXT PRIMARY KEY,
			label TEXT,
			category_id TEXT,
			uri TEXT
		)`,
		`CREATE TABLE sqlite_db.aliases(
			id TEXT,
			product_id TEXT,
			PRIMARY KEY(id, product_id)
		)`,
		`CREATE TABLE sqlite_db.details(
			product_id TEXT,
			cycle TEXT,
			is_lts INTEGER,
			release_date TEXT,
			latest TEXT,
			latest_release_date TEXT,
			eol_date TEXT,
			PRIMARY KEY(product_id, cycle)
		)`,
		`CREATE TABLE sqlite_db.product_identifiers(
			product_id TEXT,
			identifier_type TEXT,
			identifier_value TEXT,
			PRIMARY KEY(product_id, identifier_type, identifier_value)
		)`,
		`CREATE TABLE sqlite_db.product_tags(
			product_id TEXT,
			tag_id TEXT,
			PRIMARY KEY(product_id, tag_id)
		)`,
	}

	for _, ddl := range tableCreations {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("error creating SQLite table: %w", err)
		}
	}

	// Copy data with proper CASTs for timestamp/date columns
	dataInserts := []string{
		`INSERT INTO sqlite_db.about
		SELECT
			git_version, git_commit, go_version, platform, github_url,
			CAST(generated_at AS VARCHAR),
			CAST(generated_at_tz AS VARCHAR)
		FROM main.about`,

		`INSERT INTO sqlite_db.categories
		SELECT id, uri FROM main.categories`,

		`INSERT INTO sqlite_db.tags
		SELECT id, uri, www FROM main.tags`,

		`INSERT INTO sqlite_db.products
		SELECT id, label, category_id, uri FROM main.products`,

		`INSERT INTO sqlite_db.aliases
		SELECT id, product_id FROM main.aliases`,

		`INSERT INTO sqlite_db.details
		SELECT
			product_id, cycle,
			CAST(is_lts AS INTEGER),
			CAST(release_date AS VARCHAR),
			latest,
			CAST(latest_release_date AS VARCHAR),
			CAST(eol_date AS VARCHAR)
		FROM main.details`,

		`INSERT INTO sqlite_db.product_identifiers
		SELECT product_id, identifier_type, identifier_value
		FROM main.product_identifiers`,

		`INSERT INTO sqlite_db.product_tags
		SELECT product_id, tag_id FROM main.product_tags`,
	}

	for _, dml := range dataInserts {
		if _, err := db.Exec(dml); err != nil {
			return fmt.Errorf("error copying data to SQLite: %w", err)
		}
	}

	if _, err := db.Exec("DETACH sqlite_db"); err != nil {
		return fmt.Errorf("error detaching SQLite database: %w", err)
	}

	log.Info().Msg("Data exported to SQLite")
	return nil
}

// addSQLiteForeignKeys reopens the SQLite file and adds foreign key constraints
// using the rename-recreate pattern (mirrors add_fk.sql).
func addSQLiteForeignKeys(sqlitePath string) error {
	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return fmt.Errorf("error opening SQLite database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Warn().Err(err).Msg("Error closing SQLite database")
		}
	}()

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("error enabling foreign keys: %w", err)
	}

	fkMigrations := []struct {
		name       string
		statements []string
	}{
		{
			name: "products",
			statements: []string{
				"ALTER TABLE products RENAME TO products_old",
				`CREATE TABLE products(
					id TEXT PRIMARY KEY,
					label TEXT,
					category_id TEXT,
					uri TEXT,
					FOREIGN KEY (category_id) REFERENCES categories(id)
				)`,
				"INSERT INTO products SELECT * FROM products_old",
				"DROP TABLE products_old",
			},
		},
		{
			name: "aliases",
			statements: []string{
				"ALTER TABLE aliases RENAME TO aliases_old",
				`CREATE TABLE aliases(
					id TEXT,
					product_id TEXT,
					PRIMARY KEY(id, product_id),
					FOREIGN KEY (product_id) REFERENCES products(id)
				)`,
				"INSERT INTO aliases SELECT * FROM aliases_old",
				"DROP TABLE aliases_old",
			},
		},
		{
			name: "details",
			statements: []string{
				"ALTER TABLE details RENAME TO details_old",
				`CREATE TABLE details(
					product_id TEXT,
					cycle TEXT,
					is_lts INTEGER,
					release_date TEXT,
					latest TEXT,
					latest_release_date TEXT,
					eol_date TEXT,
					PRIMARY KEY(product_id, cycle),
					FOREIGN KEY (product_id) REFERENCES products(id)
				)`,
				"INSERT INTO details SELECT * FROM details_old",
				"DROP TABLE details_old",
			},
		},
		{
			name: "product_identifiers",
			statements: []string{
				"ALTER TABLE product_identifiers RENAME TO product_identifiers_old",
				`CREATE TABLE product_identifiers(
					product_id TEXT,
					identifier_type TEXT,
					identifier_value TEXT,
					PRIMARY KEY(product_id, identifier_type, identifier_value),
					FOREIGN KEY (product_id) REFERENCES products(id)
				)`,
				"INSERT INTO product_identifiers SELECT * FROM product_identifiers_old",
				"DROP TABLE product_identifiers_old",
			},
		},
		{
			name: "product_tags",
			statements: []string{
				"ALTER TABLE product_tags RENAME TO product_tags_old",
				`CREATE TABLE product_tags(
					product_id TEXT,
					tag_id TEXT,
					PRIMARY KEY(product_id, tag_id),
					FOREIGN KEY (product_id) REFERENCES products(id),
					FOREIGN KEY (tag_id) REFERENCES tags(id)
				)`,
				"INSERT INTO product_tags SELECT * FROM product_tags_old",
				"DROP TABLE product_tags_old",
			},
		},
	}

	for _, migration := range fkMigrations {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction for %s: %w", migration.name, err)
		}
		for _, stmt := range migration.statements {
			if _, err := tx.Exec(stmt); err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Warn().Err(rbErr).Msgf("Error rolling back transaction for %s", migration.name)
				}
				return fmt.Errorf("error in FK migration for %s: %w", migration.name, err)
			}
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing FK migration for %s: %w", migration.name, err)
		}
		log.Info().Msgf("Added foreign key constraints to \"%s\" table", migration.name)
	}

	return nil
}

func init() {
	ExportCmd.AddCommand(sqliteCmd)
	sqliteCmd.Flags().StringP("output", "o", "geol.sqlite", "Output SQLite database file path")
	sqliteCmd.Flags().BoolP("force", "f", false, "Overwrites the SQLite database file if it already exists")
}
