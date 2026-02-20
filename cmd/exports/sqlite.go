package exports

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/v2/utilities"
)

// createAboutTableSQLite creates the 'about' table and inserts metadata for SQLite
func createAboutTableSQLite(db *sql.DB) error {
	// Create 'about' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS about (
			git_version TEXT,
			git_commit TEXT,
			go_version TEXT,
			platform TEXT,
			github_url TEXT,
			generated_at TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
			generated_at_tz TEXT DEFAULT (datetime('now', 'localtime'))
		)`)
	if err != nil {
		return fmt.Errorf("error creating 'about' table: %w", err)
	}

	// Insert values into 'about' table
	_, err = db.Exec(`INSERT INTO about (git_version, git_commit, go_version, platform, github_url) 
		VALUES (?, ?, ?, ?, ?)`,
		utilities.Version, utilities.Commit, utilities.GoVersion,
		fmt.Sprintf("%s/%s", utilities.PlatformOs, utilities.PlatformArch),
		"https://github.com/opt-nc/geol/v2")
	log.Info().Msg("Inserted metadata into \"about\" table")
	if err != nil {
		return fmt.Errorf("error inserting into 'about' table: %w", err)
	}

	return nil
}

// createTempDetailsTableSQLite creates the 'details_temp' table and inserts product details for SQLite
func createTempDetailsTableSQLite(db *sql.DB, allData *productDataMap) error {
	// Create 'details_temp' table if not exists
	_, err := db.Exec(`CREATE TEMP TABLE IF NOT EXISTS details_temp (
			product_id TEXT,
			cycle TEXT,
			release TEXT,
			latest TEXT,
			latest_release TEXT,
			eol TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'details_temp' table")
		return err
	}

	// Get product IDs from the products table
	rows, err := db.Query(`SELECT id FROM products`)
	if err != nil {
		log.Error().Err(err).Msg("Error querying products from database")
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing rows")
		}
	}()

	var productIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Error().Err(err).Msg("Error scanning product ID")
			return err
		}
		productIDs = append(productIDs, id)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over product rows")
		return err
	}

	// Insert product details for each product in the database
	for _, productID := range productIDs {
		if prodData, exists := allData.Products[productID]; exists {
			// Insert each release into the details_temp table
			for _, release := range prodData.Releases {
				_, err = db.Exec(`INSERT INTO details_temp (product_id, cycle, release, latest, latest_release, eol) 
						VALUES (?, ?, ?, ?, ?, ?)`,
					productID,
					release.Name,
					release.ReleaseDate,
					release.LatestName,
					release.LatestDate,
					release.EolFrom,
				)
				if err != nil {
					log.Error().Err(err).Msgf("Error inserting release data for %s", productID)
				}
			}
		}
	}

	return nil
}

// createDetailsTableSQLite creates the final 'details' table from 'details_temp' with proper date types for SQLite
func createDetailsTableSQLite(db *sql.DB) error {
	// Create 'details' table with TEXT columns for dates (SQLite standard)
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS details (
			product_id TEXT,
			cycle TEXT,
			release_date TEXT,
			latest TEXT,
			latest_release_date TEXT,
			eol_date TEXT,
			PRIMARY KEY (product_id, cycle),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'details' table")
		return err
	}

	// Insert data from details_temp, converting empty strings to NULL
	_, err = db.Exec(`INSERT INTO details (product_id, cycle, release_date, latest, latest_release_date, eol_date)
		SELECT 
			product_id,
			cycle,
			CASE WHEN release = '' THEN NULL ELSE release END,
			latest,
			CASE WHEN latest_release = '' THEN NULL ELSE latest_release END,
			CASE WHEN eol = '' THEN NULL ELSE eol END
		FROM details_temp
		ORDER BY product_id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting data into 'details' table")
		return err
	}

	log.Info().Msg("Created and populated \"details\" table")

	return nil
}

// createProductsTableSQLite creates the 'products' table and inserts product information for SQLite
func createProductsTableSQLite(db *sql.DB, allData *productDataMap) error {
	// Create 'products' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			label TEXT,
			category_id TEXT,
			uri TEXT,
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'products' table")
		return err
	}

	// Collect all products
	type productEntry struct {
		id         string
		label      string
		categoryID string
		uri        string
	}
	var allProductsSlice []productEntry

	for _, prodData := range allData.Products {
		allProductsSlice = append(allProductsSlice, productEntry{
			id:         prodData.Name,
			label:      prodData.Label,
			categoryID: prodData.Category,
			uri:        prodData.URI,
		})
	}

	// Sort products by id using SQLite
	// First insert all data into a temporary table, then insert sorted
	_, err = db.Exec(`CREATE TEMP TABLE IF NOT EXISTS products_temp (
			id TEXT,
			label TEXT,
			category_id TEXT,
			uri TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating products_temp table")
		return err
	}

	for _, entry := range allProductsSlice {
		_, err = db.Exec(`INSERT INTO products_temp (id, label, category_id, uri) VALUES (?, ?, ?, ?)`,
			entry.id,
			entry.label,
			entry.categoryID,
			entry.uri,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting product %s into temp table", entry.id)
		}
	}

	// Insert from temp table sorted by id
	_, err = db.Exec(`INSERT INTO products (id, label, category_id, uri) 
		SELECT id, label, category_id, uri FROM products_temp ORDER BY id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting sorted products")
		return err
	}

	log.Info().Msg("Created and populated \"products\" table")

	return nil
}

func createAliasesTableSQLite(db *sql.DB, allData *productDataMap) error {
	// Create 'aliases' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS aliases (
			id TEXT,
			product_id TEXT,
			PRIMARY KEY (id, product_id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'aliases' table")
		return err
	}

	// Get product IDs from the products table
	rows, err := db.Query(`SELECT id FROM products`)
	if err != nil {
		log.Error().Err(err).Msg("Error querying products from database")
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing rows")
		}
	}()

	var productIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Error().Err(err).Msg("Error scanning product ID")
			return err
		}
		productIDs = append(productIDs, id)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over product rows")
		return err
	}

	// Collect all aliases with their product IDs
	type aliasEntry struct {
		id        string
		productID string
	}
	var allAliases []aliasEntry

	for _, productID := range productIDs {
		if prodData, exists := allData.Products[productID]; exists {
			for _, alias := range prodData.Aliases {
				allAliases = append(allAliases, aliasEntry{
					id:        alias,
					productID: productID,
				})
			}
		}
	}

	// Sort aliases by id using SQLite
	// First insert all data into a temporary table, then insert sorted
	_, err = db.Exec(`CREATE TEMP TABLE IF NOT EXISTS aliases_temp (
			id TEXT,
			product_id TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating aliases_temp table")
		return err
	}

	for _, entry := range allAliases {
		_, err = db.Exec(`INSERT INTO aliases_temp (id, product_id) VALUES (?, ?)`,
			entry.id,
			entry.productID,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting alias %s into temp table", entry.id)
		}
	}

	// Insert from temp table sorted by id
	_, err = db.Exec(`INSERT INTO aliases (id, product_id) 
		SELECT id, product_id FROM aliases_temp ORDER BY id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting sorted aliases")
		return err
	}

	log.Info().Msg("Created and populated \"aliases\" table")

	return nil
}

func createProductIdentifiersTableSQLite(db *sql.DB, allData *productDataMap) error {
	// Create 'product_identifiers' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS product_identifiers (
			product_id TEXT,
			identifier_type TEXT,
			identifier_value TEXT,
			PRIMARY KEY (product_id, identifier_type, identifier_value),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'product_identifiers' table")
		return err
	}

	// Get product IDs from the products table
	rows, err := db.Query(`SELECT id FROM products`)
	if err != nil {
		log.Error().Err(err).Msg("Error querying products from database")
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing rows")
		}
	}()

	var productIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Error().Err(err).Msg("Error scanning product ID")
			return err
		}
		productIDs = append(productIDs, id)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over product rows")
		return err
	}

	// Collect all product identifiers
	type identifierEntry struct {
		productID       string
		identifierType  string
		identifierValue string
	}
	var allIdentifiers []identifierEntry

	for _, productID := range productIDs {
		if prodData, exists := allData.Products[productID]; exists {
			for _, identifier := range prodData.Identifiers {
				// Special handling for repology identifiers - store full URL
				identifierValue := identifier.ID
				if identifier.Type == "repology" {
					identifierValue = "https://repology.org/project/" + identifier.ID
				}

				allIdentifiers = append(allIdentifiers, identifierEntry{
					productID:       productID,
					identifierType:  identifier.Type,
					identifierValue: identifierValue,
				})
			}
		}
	}

	// Sort product identifiers by product_id using SQLite
	// First insert all data into a temporary table, then insert sorted
	_, err = db.Exec(`CREATE TEMP TABLE IF NOT EXISTS product_identifiers_temp (
			product_id TEXT,
			identifier_type TEXT,
			identifier_value TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating product_identifiers_temp table")
		return err
	}

	for _, entry := range allIdentifiers {
		_, err = db.Exec(`INSERT INTO product_identifiers_temp (product_id, identifier_type, identifier_value) VALUES (?, ?, ?)`,
			entry.productID,
			entry.identifierType,
			entry.identifierValue,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting identifier into temp table")
		}
	}

	// Insert from temp table sorted by product_id
	_, err = db.Exec(`INSERT INTO product_identifiers (product_id, identifier_type, identifier_value) 
		SELECT product_id, identifier_type, identifier_value FROM product_identifiers_temp ORDER BY product_id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting sorted product identifiers")
		return err
	}

	log.Info().Msg("Created and populated \"product_identifiers\" table")

	return nil
}

func createTagsTableSQLite(db *sql.DB, allTags map[string]utilities.Tag) error {
	// Create 'tags' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			uri TEXT UNIQUE,
			www TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'tags' table")
		return err
	}

	// Collect all tags
	type tagEntry struct {
		id  string
		uri string
		www string
	}
	var allTagsSlice []tagEntry

	for _, tag := range allTags {
		allTagsSlice = append(allTagsSlice, tagEntry{
			id:  tag.Name,
			uri: tag.Uri,
			www: "https://endoflife.date/tags/" + tag.Name,
		})
	}

	// Sort tags by id using SQLite
	// First insert all data into a temporary table, then insert sorted
	_, err = db.Exec(`CREATE TEMP TABLE IF NOT EXISTS tags_temp (
			id TEXT,
			uri TEXT,
			www TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating tags_temp table")
		return err
	}

	for _, entry := range allTagsSlice {
		_, err = db.Exec(`INSERT INTO tags_temp (id, uri, www) VALUES (?, ?, ?)`,
			entry.id,
			entry.uri,
			entry.www,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting tag %s into temp table", entry.id)
		}
	}

	// Insert from temp table sorted by id
	_, err = db.Exec(`INSERT INTO tags (id, uri, www) 
		SELECT id, uri, www FROM tags_temp ORDER BY id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting sorted tags")
		return err
	}

	log.Info().Msg("Created and populated \"tags\" table")

	return nil
}

func createCategoriesTableSQLite(db *sql.DB, allCategories map[string]utilities.Category) error {
	// Create 'categories' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			uri TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'categories' table")
		return err
	}

	// Collect all categories
	type categoryEntry struct {
		id  string
		uri string
	}
	var allCategoriesSlice []categoryEntry

	for _, category := range allCategories {
		allCategoriesSlice = append(allCategoriesSlice, categoryEntry{
			id:  category.Name,
			uri: category.Uri,
		})
	}

	// Sort categories by id using SQLite
	// First insert all data into a temporary table, then insert sorted
	_, err = db.Exec(`CREATE TEMP TABLE IF NOT EXISTS categories_temp (
			id TEXT,
			uri TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating categories_temp table")
		return err
	}

	for _, entry := range allCategoriesSlice {
		_, err = db.Exec(`INSERT INTO categories_temp (id, uri) VALUES (?, ?)`,
			entry.id,
			entry.uri,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting category %s into temp table", entry.id)
		}
	}

	// Insert from temp table sorted by id
	_, err = db.Exec(`INSERT INTO categories (id, uri) 
		SELECT id, uri FROM categories_temp ORDER BY id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting sorted categories")
		return err
	}

	log.Info().Msg("Created and populated \"categories\" table")

	return nil
}

func createProductTagsTableSQLite(db *sql.DB, allData *productDataMap) error {
	// Create 'product_tags' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS product_tags (
			product_id TEXT,
			tag_id TEXT,
			PRIMARY KEY (product_id, tag_id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'product_tags' table")
		return err
	}

	// Collect all product-tag relationships
	type productTagEntry struct {
		productID string
		tagID     string
	}
	var allProductTags []productTagEntry

	for _, prodData := range allData.Products {
		for _, tag := range prodData.Tags {
			allProductTags = append(allProductTags, productTagEntry{
				productID: prodData.Name,
				tagID:     tag.Name,
			})
		}
	}

	// Sort product-tag relationships by product_id using SQLite
	// First insert all data into a temporary table, then insert sorted
	_, err = db.Exec(`CREATE TEMP TABLE IF NOT EXISTS product_tags_temp (
			product_id TEXT,
			tag_id TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating product_tags_temp table")
		return err
	}

	for _, entry := range allProductTags {
		_, err = db.Exec(`INSERT INTO product_tags_temp (product_id, tag_id) VALUES (?, ?)`,
			entry.productID,
			entry.tagID,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting product-tag into temp table")
		}
	}

	// Insert from temp table sorted by product_id
	_, err = db.Exec(`INSERT INTO product_tags (product_id, tag_id) 
		SELECT product_id, tag_id FROM product_tags_temp ORDER BY product_id`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting sorted product-tags")
		return err
	}

	log.Info().Msg("Created and populated \"product_tags\" table")

	return nil
}

// sqliteCmd represents the sqlite command
var sqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Export data to a SQLite database",
	Long: `Export all known products and their end-of-life (EOL) metadata into a SQLite database file.
This command creates a new SQLite file (default: geol.sqlite) and populates it with
information such as version details, platform info, and comprehensive product lifecycle data.

You can specify the output filename using the --output flag.
If the file already exists, use the --force flag to overwrite it.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		utilities.AnalyzeCacheProductsValidity(cmd)
		dbPath, _ := cmd.Flags().GetString("output")
		forceSQLite, _ := cmd.Flags().GetBool("force")
		if _, err := os.Stat(dbPath); err == nil {
			if !forceSQLite {
				log.Fatal().Msgf("File %s already exists. Use --force to overwrite.", dbPath)
			}
			if err := os.Remove(dbPath); err != nil {
				log.Fatal().Err(err).Msgf("Error removing existing %s file", dbPath)
			}
		}

		// Open or create SQLite database
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while creating SQLite database")
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Fatal().Err(err).Msg("Error closing SQLite database")
			}
		}()

		// Enable foreign keys in SQLite
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			log.Fatal().Err(err).Msg("Error enabling foreign keys")
		}

		// Fetch all product data from API in a single pass
		allProductsData, err := fetchAllProductData(cmd)
		if err != nil {
			log.Fatal().Err(err).Msg("Error fetching product data from API")
		}

		allTags, err := fetchAllTags()
		if err != nil {
			log.Fatal().Err(err).Msg("Error fetching tags from API")
		}

		allCategories, err := fetchAllCategories()
		if err != nil {
			log.Fatal().Err(err).Msg("Error fetching categories from API")
		}

		// Create 'tags' table and insert tags
		if err := createTagsTableSQLite(db, allTags); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'tags' table")
		}

		// Create 'categories' table and insert categories
		if err := createCategoriesTableSQLite(db, allCategories); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'categories' table")
		}

		// Create 'products' table and insert product information
		if err := createProductsTableSQLite(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'products' table")
		}

		// Create 'details_temp' table and insert product details
		if err := createTempDetailsTableSQLite(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details_temp' table")
		}

		// Create 'details' table from 'details_temp' with proper date types
		if err := createDetailsTableSQLite(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details' table")
		}

		// Create 'aliases' table and insert product aliases
		if err := createAliasesTableSQLite(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'aliases' table")
		}

		// Create 'product_identifiers' table and insert product identifiers
		if err := createProductIdentifiersTableSQLite(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'product_identifiers' table")
		}

		// Create 'product_tags' junction table
		if err := createProductTagsTableSQLite(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'product_tags' table")
		}

		// Create 'about' table and insert metadata
		if err := createAboutTableSQLite(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'about' table")
		}

		duration := time.Since(startTime)
		log.Info().Msgf("SQLite database created successfully at %s (took %v)", dbPath, duration.Round(time.Millisecond))
		log.Info().Msg("You can query the database using SQLite CLI or any compatible client.")
		log.Info().Msgf("Example CLI command: sqlite3 %s", dbPath)
	},
}

func init() {
	ExportCmd.AddCommand(sqliteCmd)
	sqliteCmd.Flags().StringP("output", "o", "geol.sqlite", "Output SQLite database file path")
	sqliteCmd.Flags().BoolP("force", "f", false, "Overwrites the SQLite database file if it already exists")
}
