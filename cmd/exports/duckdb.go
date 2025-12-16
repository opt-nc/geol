package exports

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/cmd/product"
	"github.com/opt-nc/geol/utilities"
)

// createAboutTable creates the 'about' table and inserts metadata
func createAboutTable(db *sql.DB) error {
	// Create 'about' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS about (
			gitVersion TEXT,
			gitCommit TEXT,
			goVersion TEXT,
			platform TEXT,
			githubURL TEXT,
			generatedAt TIMESTAMP DEFAULT date_trunc('second', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			generatedAtTZ TIMESTAMPTZ DEFAULT date_trunc('second', CURRENT_TIMESTAMP)
		)`)
	if err != nil {
		return fmt.Errorf("error creating 'about' table: %w", err)
	}

	// Insert values into 'about' table
	_, err = db.Exec(`INSERT INTO about (gitVersion, gitCommit, goVersion, platform, githubURL) 
		VALUES (?, ?, ?, ?, ?)`,
		utilities.Version, utilities.Commit, utilities.GoVersion,
		fmt.Sprintf("%s/%s", utilities.PlatformOs, utilities.PlatformArch),
		"https://github.com/opt-nc/geol")
	if err != nil {
		return fmt.Errorf("error inserting into 'about' table: %w", err)
	}

	return nil
}

// createDetailsTable creates the 'details' table and inserts product details
func createDetailsTable(cmd *cobra.Command, db *sql.DB) error {
	// Create 'details' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS details (
			product_id TEXT,
			cycle TEXT,
			release TEXT,
			latest TEXT,
			latest_release TEXT,
			eol TEXT
		)`)
	if err != nil {
		return fmt.Errorf("error creating 'details' table: %w", err)
	}

	// Get products from cache
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		return fmt.Errorf("error retrieving products path: %w", err)
	}

	products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
	if err != nil {
		return fmt.Errorf("error retrieving products from cache: %w", err)
	}

	// Insert data for each product
	for productName := range products.Products {
		prodData, err := product.FetchProductData(productName)
		if err != nil {
			log.Warn().Err(err).Msgf("Error fetching product data for %s, skipping", productName)
			continue
		}

		// Insert each release into the details table
		for _, release := range prodData.Releases {
			_, err = db.Exec(`INSERT INTO details (product_id, cycle, release, latest, latest_release, eol) 
				VALUES (?, ?, ?, ?, ?, ?)`,
				prodData.Name,
				release.Name,
				release.ReleaseDate,
				release.LatestName,
				release.LatestDate,
				release.EolFrom,
			)
			if err != nil {
				return fmt.Errorf("error inserting release data for %s: %w", productName, err)
			}
		}
	}

	return nil
}

// duckdbCmd represents the duckdb command
var duckdbCmd = &cobra.Command{
	Use:   "duckdb",
	Short: "Export data to a DuckDB database",
	Long: `Export the application data, metadata, and product information from the endoflife.date API into a DuckDB database file.
This command creates a new DuckDB file (default: geol.duckdb) and populates it with
information such as version details, platform info, and comprehensive product lifecycle data.

You can specify the output filename using the --output flag.
If the file already exists, use the --force flag to overwrite it.`,
	Run: func(cmd *cobra.Command, args []string) {
		utilities.AnalyzeCacheProductsValidity(cmd)
		dbPath, _ := cmd.Flags().GetString("output")
		forceDuckDB, _ := cmd.Flags().GetBool("force")
		if _, err := os.Stat(dbPath); err == nil {
			if !forceDuckDB {
				log.Fatal().Msgf("File %s already exists. Use --force to overwrite.", dbPath)
			}
			if err := os.Remove(dbPath); err != nil {
				log.Fatal().Err(err).Msgf("Error removing existing %s file", dbPath)
			}
		}

		// Open or create geol.duckdb at project root
		db, err := sql.Open("duckdb", dbPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while creating DuckDB database")
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Fatal().Err(err).Msg("Error closing DuckDB database")
			}
		}()

		// Create 'about' table and insert metadata
		if err := createAboutTable(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'about' table")
		}

		// Create 'details' table and insert product details
		if err := createDetailsTable(cmd, db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details' table")
		}

		log.Info().Msg("DuckDB database created successfully at " + dbPath)
	},
}

func init() {
	ExportCmd.AddCommand(duckdbCmd)
	duckdbCmd.Flags().StringP("output", "o", "geol.duckdb", "Output DuckDB database file path")
	duckdbCmd.Flags().BoolP("force", "f", false, "Force overwrite of existing geol.duckdb file")
}
