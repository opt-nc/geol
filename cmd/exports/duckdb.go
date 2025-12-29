package exports

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/cmd/product"
	"github.com/opt-nc/geol/utilities"
)

const (
	padding  = 2
	maxWidth = 60
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type productProcessedMsg string

type model struct {
	progress          progress.Model
	totalProducts     int
	processed         int
	done              bool
	progressMessage   string
	completionMessage string
}

type identifiers struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type productData struct {
	Name        string                `json:"name"`
	Aliases     []string              `json:"aliases"`
	Label       string                `json:"label"`
	Category    string                `json:"category"`
	Tags        []utilities.Tag       `json:"tags"`
	Identifiers []identifiers         `json:"identifiers"`
	URI         string                `json:"links.html"`
	Releases    []product.ReleaseInfo `json:"releases"`
}

type allProductsData struct {
	Products map[string]*productData
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		width := msg.Width - padding*2 - 4
		if width > maxWidth {
			width = maxWidth
		}
		m.progress = progress.New(progress.WithWidth(width))
		return m, nil

	case productProcessedMsg:
		m.processed++
		percent := float64(m.processed) / float64(m.totalProducts)
		if m.processed >= m.totalProducts {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Quit,
			)
		}
		return m, m.progress.SetPercent(percent)

	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() tea.View {
	pad := strings.Repeat(" ", padding)
	count := fmt.Sprintf(" %d/%d", m.processed, m.totalProducts)
	content := "\n" +
		pad + m.progress.View() + count + "\n\n"
	if m.done {
		content += pad + lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓") + " " + m.completionMessage + "\n"
	} else {
		content += pad + helpStyle(m.progressMessage)
	}

	return tea.NewView(content)
}

// fetchAllProductData retrieves all product information and details from the API in a single pass
func fetchAllProductData(cmd *cobra.Command) (*allProductsData, error) {
	// Get products from cache
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products path")
		return nil, err
	}

	products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products from cache")
		return nil, err
	}

	allData := &allProductsData{
		Products: make(map[string]*productData),
	}

	// Initialize the bubbletea model with progress bar
	m := model{
		progress:          progress.New(progress.WithWidth(maxWidth)),
		totalProducts:     len(products.Products),
		processed:         0,
		done:              false,
		progressMessage:   "Fetching all product data from API...",
		completionMessage: "All product data fetched!",
	}

	// Start the TUI in a goroutine
	p := tea.NewProgram(m)

	// Process products in a goroutine
	go func() {
		for productName := range products.Products {
			// Fetch basic product info
			url := utilities.ApiUrl + "products/" + productName
			resp, err := http.Get(url)
			if err != nil {
				log.Warn().Err(err).Msgf("Error requesting %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if cerr := resp.Body.Close(); cerr != nil {
				log.Warn().Err(cerr).Msgf("Error closing HTTP body for %s", productName)
			}
			if err != nil {
				log.Warn().Err(err).Msgf("Error reading response for %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
				continue
			}

			if resp.StatusCode != 200 {
				log.Warn().Msgf("Product %s not found on the API (status %d), skipping", productName, resp.StatusCode)
				p.Send(productProcessedMsg(productName))
				continue
			}

			// Parse JSON response
			var apiResp struct {
				Result struct {
					Name        string        `json:"name"`
					Aliases     []string      `json:"aliases"`
					Label       string        `json:"label"`
					Category    string        `json:"category"`
					Tags        []string      `json:"tags"`
					Identifiers []identifiers `json:"identifiers"`
					Links       struct {
						Html string `json:"html"`
					} `json:"links"`
				} `json:"result"`
			}

			if err := json.Unmarshal(body, &apiResp); err != nil {
				log.Warn().Err(err).Msgf("Error decoding JSON for %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
				continue
			}

			// Fetch product details (releases)
			prodData, err := product.FetchProductData(productName)
			if err != nil {
				log.Warn().Err(err).Msgf("Error fetching product data for %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
				continue
			}

			// Store all data
			allData.Products[productName] = &productData{
				Name:        apiResp.Result.Name,
				Aliases:     apiResp.Result.Aliases,
				Label:       apiResp.Result.Label,
				Category:    apiResp.Result.Category,
				Tags:        utilities.ConvertTagStringsToTags(apiResp.Result.Tags),
				Identifiers: apiResp.Result.Identifiers,
				URI:         apiResp.Result.Links.Html,
				Releases:    prodData.Releases,
			}

			p.Send(productProcessedMsg(productName))
		}
	}()

	// Run the program and wait for completion
	if _, err := p.Run(); err != nil {
		log.Error().Err(err).Msg("Error running progress display")
		return nil, err
	}

	log.Info().Msgf("Fetched data for %d products", len(allData.Products))
	return allData, nil
}

// fetchAllCategories retrieves all categories from the API
func fetchAllCategories() (map[string]utilities.Category, error) {
	// Fetch categories from API
	resp, err := http.Get(utilities.ApiUrl + "categories")
	if err != nil {
		log.Error().Err(err).Msg("Error fetching categories")
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing response body")
		}
	}()

	var apiResp struct {
		Result []utilities.Category `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Error().Err(err).Msg("Error decoding categories JSON")
		return nil, err
	}

	categories := make(map[string]utilities.Category)
	for _, category := range apiResp.Result {
		categories[category.Name] = category
	}

	log.Info().Msgf("Fetched %d categories", len(categories))
	return categories, nil
}

// fetchAllTags retrieves all tags from the API
func fetchAllTags() (map[string]utilities.Tag, error) {
	// Fetch tags from API
	resp, err := http.Get(utilities.ApiUrl + "tags")
	if err != nil {
		log.Error().Err(err).Msg("Error fetching tags")
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing response body")
		}
	}()

	var apiResp struct {
		Result []utilities.Tag `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Error().Err(err).Msg("Error decoding tags JSON")
		return nil, err
	}

	tags := make(map[string]utilities.Tag)
	for _, tag := range apiResp.Result {
		tags[tag.Name] = tag
	}

	log.Info().Msgf("Fetched %d tags", len(tags))
	return tags, nil
}

// createAboutTable creates the 'about' table and inserts metadata
func createAboutTable(db *sql.DB) error {

	// Create 'about' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS about (
			git_version TEXT,
			git_commit TEXT,
			go_version TEXT,
			platform TEXT,
			github_URL TEXT,
			generated_at TIMESTAMP DEFAULT date_trunc('second', CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			generated_at_TZ TIMESTAMPTZ DEFAULT date_trunc('second', CURRENT_TIMESTAMP)
		)`)
	if err != nil {
		return fmt.Errorf("error creating 'about' table: %w", err)
	}

	// Add comment to 'about' table
	_, err = db.Exec(`COMMENT ON TABLE about IS 'Metadata about the geol version and platform that generated this database'`)
	if err != nil {
		return fmt.Errorf("error adding comment to 'about' table: %w", err)
	}

	// Add comments to 'about' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN about.git_version IS 'Git tag version of geol used to generate this database';
		COMMENT ON COLUMN about.git_commit IS 'Git commit hash of geol used to generate this database';
		COMMENT ON COLUMN about.go_version IS 'Go compiler version used to build geol';
		COMMENT ON COLUMN about.platform IS 'Operating system and architecture where geol was executed';
		COMMENT ON COLUMN about.github_URL IS 'GitHub repository URL for the geol project';
		COMMENT ON COLUMN about.generated_at IS 'UTC timestamp when this database was generated';
		COMMENT ON COLUMN about.generated_at_TZ IS 'Local timestamp with timezone when this database was generated';
	`)
	if err != nil {
		return fmt.Errorf("error adding comments to 'about' columns: %w", err)
	}

	// Insert values into 'about' table
	_, err = db.Exec(`INSERT INTO about (git_version, git_commit, go_version, platform, github_URL) 
		VALUES (?, ?, ?, ?, ?)`,
		utilities.Version, utilities.Commit, utilities.GoVersion,
		fmt.Sprintf("%s/%s", utilities.PlatformOs, utilities.PlatformArch),
		"https://github.com/opt-nc/geol")
	log.Info().Msg("Inserted metadata into \"about\" table")
	if err != nil {
		return fmt.Errorf("error inserting into 'about' table: %w", err)
	}

	return nil
}

// createTempDetailsTable creates the 'details_temp' table and inserts product details
func createTempDetailsTable(db *sql.DB, allData *allProductsData) error {
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

	// Add comment to 'details_temp' table
	_, err = db.Exec(`COMMENT ON TABLE details_temp IS 'Temporary table for storing raw product release details before type conversion'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'details_temp' table")
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

	//log.Info().Msg("Populated \"details_temp\" table")
	return nil
}

// createDetailsTable creates the final 'details' table from 'details_temp' with proper date types
func createDetailsTable(db *sql.DB) error {
	// Create 'details' table with DATE columns for release_date, latest_release_date, and eol
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS details (
			product_id TEXT,
			cycle TEXT,
			release_date DATE,
			latest TEXT,
			latest_release_date DATE,
			eol_date DATE,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'details' table")
		return err
	}

	// Insert data from details_temp, converting empty strings to NULL and casting to DATE
	_, err = db.Exec(`INSERT INTO details (product_id, cycle, release_date, latest, latest_release_date, eol_date)
		SELECT 
			product_id,
			cycle,
			CASE WHEN release = '' THEN NULL ELSE TRY_CAST(release AS DATE) END,
			latest,
			CASE WHEN latest_release = '' THEN NULL ELSE TRY_CAST(latest_release AS DATE) END,
			CASE WHEN eol = '' THEN NULL ELSE TRY_CAST(eol AS DATE) END
		FROM details_temp`)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting data into 'details' table")
		return err
	}

	// Add comment to 'details' table
	_, err = db.Exec(`COMMENT ON TABLE details IS 'Product release lifecycle details including release dates, latest versions, and end-of-life (EOL) dates'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'details' table")
		return err
	}

	// Add comments to 'details' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN details.product_id IS 'Product id referencing the products table';
		COMMENT ON COLUMN details.cycle IS 'Product release cycle or version number';
		COMMENT ON COLUMN details.release_date IS 'Initial release date for this cycle';
		COMMENT ON COLUMN details.latest IS 'Latest patch version within this cycle';
		COMMENT ON COLUMN details.latest_release_date IS 'Release date of the latest patch version';
		COMMENT ON COLUMN details.eol_date IS 'End-of-life date when this cycle stops receiving support';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'details' columns")
		return err
	}

	log.Info().Msg("Created and populated \"details\" table")

	return nil
}

// createProductsTable creates the 'products' table and inserts product information
func createProductsTable(db *sql.DB, allData *allProductsData) error {
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

	// Add comment to 'products' table
	_, err = db.Exec(`COMMENT ON TABLE products IS 'Catalog of all products tracked by geol with their labels, categories, and documentation URIs'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'products' table")
		return err
	}

	// Add comments to 'products' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN products.id IS 'Unique product id (primary key)';
		COMMENT ON COLUMN products.label IS 'Human-readable display name for the product';
		COMMENT ON COLUMN products.category_id IS 'Category id grouping related products';
		COMMENT ON COLUMN products.uri IS 'URI to the product documentation on endoflife.date';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'products' columns")
		return err
	}

	// Insert product data from pre-fetched data
	for _, prodData := range allData.Products {
		_, err = db.Exec(`INSERT INTO products (id, label, category_id, uri) 
				VALUES (?, ?, ?, ?)`,
			prodData.Name,
			prodData.Label,
			prodData.Category,
			prodData.URI,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting product data for %s", prodData.Name)
		}
	}

	log.Info().Msg("Created and populated \"products\" table")

	return nil
}

func createAliasesTable(db *sql.DB, allData *allProductsData) error {
	// Create 'aliases' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS aliases (
			product_id TEXT,
			alias TEXT,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'aliases' table")
		return err
	}

	// Add comment to 'aliases' table
	_, err = db.Exec(`COMMENT ON TABLE aliases IS 'Alternative names or aliases for products to facilitate searching and identification'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'aliases' table")
		return err
	}
	// Add comments to 'aliases' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN aliases.product_id IS 'Product id referencing the products table';
		COMMENT ON COLUMN aliases.alias IS 'Alternative name or alias for the product';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'aliases' columns")
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

	// Insert aliases for each product in the database
	for _, productID := range productIDs {
		if prodData, exists := allData.Products[productID]; exists {
			for _, alias := range prodData.Aliases {
				_, err = db.Exec(`INSERT INTO aliases (product_id, alias) 
						VALUES (?, ?)`,
					productID,
					alias,
				)
				if err != nil {
					log.Error().Err(err).Msgf("Error inserting alias %s for product %s", alias, productID)
				}
			}
		}
	}

	log.Info().Msg("Created and populated \"aliases\" table")

	return nil
}

func createProductIdentifiersTable(db *sql.DB, allData *allProductsData) error {
	// Create 'product_identifiers' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS product_identifiers (
			product_id TEXT,
			identifier_type TEXT,
			identifier_value TEXT,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'product_identifiers' table")
		return err
	}

	// Add comment to 'product_identifiers' table
	_, err = db.Exec(`COMMENT ON TABLE product_identifiers IS 'Various identifiers for products such as SKUs, model numbers, or codes used by manufacturers'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'product_identifiers' table")
		return err
	}

	// Add comments to 'product_identifiers' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN product_identifiers.product_id IS 'Product id referencing the products table';
		COMMENT ON COLUMN product_identifiers.identifier_type IS 'Type of identifier (e.g., repology, purl, cpe)';
		COMMENT ON COLUMN product_identifiers.identifier_value IS 'Value of the identifier';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'product_identifiers' columns")
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
	// Insert product identifiers for each product in the database
	for _, productID := range productIDs {
		if prodData, exists := allData.Products[productID]; exists {
			for _, identifier := range prodData.Identifiers {
				// Special handling for repology identifiers - store full URL
				identifierValue := identifier.ID
				if identifier.Type == "repology" {
					identifierValue = "https://repology.org/project/" + identifier.ID
				}

				_, err = db.Exec(`INSERT INTO product_identifiers (product_id, identifier_type, identifier_value)
						VALUES (?, ?, ?)`,
					productID,
					identifier.Type,
					identifierValue,
				)
				if err != nil {
					log.Error().Err(err).Msgf("Error inserting identifier %s for product %s", identifier.ID, productID)
				}
			}
		}
	}

	log.Info().Msg("Created and populated \"product_identifiers\" table")

	return nil
}

func createTagsTable(db *sql.DB, allTags map[string]utilities.Tag) error {
	// Create 'tags' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			uri TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'tags' table")
		return err
	}

	// Add comment to 'tags' table
	_, err = db.Exec(`COMMENT ON TABLE tags IS 'Tags used to categorize and group products by common characteristics or use cases'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'tags' table")
		return err
	}

	// Add comments to 'tags' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN tags.id IS 'Unique tag identifier (primary key)';
		COMMENT ON COLUMN tags.uri IS 'URI to the tag page on endoflife.date';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'tags' columns")
		return err
	}

	// Insert tags data
	for _, tag := range allTags {
		_, err = db.Exec(`INSERT INTO tags (id, uri) 
				VALUES (?, ?)`,
			tag.Name,
			tag.Uri,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting tag %s", tag.Name)
		}
	}

	log.Info().Msg("Created and populated \"tags\" table")

	return nil
}

func createCategoriesTable(db *sql.DB, allCategories map[string]utilities.Category) error {
	// Create 'categories' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			uri TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'categories' table")
		return err
	}

	// Add comment to 'categories' table
	_, err = db.Exec(`COMMENT ON TABLE categories IS 'Categories used to group products by their primary function or domain'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'categories' table")
		return err
	}

	// Add comments to 'categories' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN categories.id IS 'Unique category identifier (primary key)';
		COMMENT ON COLUMN categories.uri IS 'URI to the category page on endoflife.date';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'categories' columns")
		return err
	}

	// Insert categories data
	for _, category := range allCategories {
		_, err = db.Exec(`INSERT INTO categories (id, uri) 
				VALUES (?, ?)`,
			category.Name,
			category.Uri,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting category %s", category.Name)
		}
	}

	log.Info().Msg("Created and populated \"categories\" table")

	return nil
}

func createProductTagsTable(db *sql.DB, allData *allProductsData) error {
	// Create 'product_tags' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS product_tags (
			product_id TEXT,
			tag_id TEXT,
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'product_tags' table")
		return err
	}

	// Add comment to 'product_tags' table
	_, err = db.Exec(`COMMENT ON TABLE product_tags IS 'Junction table linking products to their associated tags for many-to-many relationships'`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comment to 'product_tags' table")
		return err
	}

	// Add comments to 'product_tags' table columns
	_, err = db.Exec(`
		COMMENT ON COLUMN product_tags.product_id IS 'Product id referencing the products table';
		COMMENT ON COLUMN product_tags.tag_id IS 'Tag id referencing the tags table';
	`)
	if err != nil {
		log.Error().Err(err).Msg("Error adding comments to 'product_tags' columns")
		return err
	}

	// Insert product-tag relationships
	for _, prodData := range allData.Products {
		for _, tag := range prodData.Tags {
			_, err = db.Exec(`INSERT INTO product_tags (product_id, tag_id) 
					VALUES (?, ?)`,
				prodData.Name,
				tag.Name,
			)
			if err != nil {
				log.Error().Err(err).Msgf("Error inserting product-tag relationship for product %s and tag %s", prodData.Name, tag.Name)
			}
		}
	}

	log.Info().Msg("Created and populated \"product_tags\" table")

	return nil
}

// duckdbCmd represents the duckdb command
var duckdbCmd = &cobra.Command{
	Use:   "duckdb",
	Short: "Export data to a DuckDB database",
	Long: `Export all known products and their end-of-life (EOL) metadata into a DuckDB database file.
This command creates a new DuckDB file (default: geol.duckdb) and populates it with
information such as version details, platform info, and comprehensive product lifecycle data.

You can specify the output filename using the --output flag.
If the file already exists, use the --force flag to overwrite it.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

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
		if err := createTagsTable(db, allTags); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'tags' table")
		}

		// Create 'categories' table and insert categories
		if err := createCategoriesTable(db, allCategories); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'categories' table")
		}
		// Create 'products' table and insert product information
		if err := createProductsTable(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'products' table")
		}

		// Create 'details_temp' table and insert product details
		if err := createTempDetailsTable(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details_temp' table")
		}

		// Create 'details' table from 'details_temp' with proper date types
		if err := createDetailsTable(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details' table")
		}

		// Create 'aliases' table and insert product aliases
		if err := createAliasesTable(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'aliases' table")
		}

		// Create 'product_identifiers' table and insert product identifiers
		if err := createProductIdentifiersTable(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'product_identifiers' table")
		}

		// Create 'product_tags' junction table
		if err := createProductTagsTable(db, allProductsData); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'product_tags' table")
		}

		// Create 'about' table and insert metadata
		if err := createAboutTable(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'about' table")
		}

		duration := time.Since(startTime)
		log.Info().Msgf("DuckDB database created successfully at %s (took %v)", dbPath, duration.Round(time.Millisecond))
		log.Info().Msg("You can query the database using DuckDB CLI or any compatible client.")
		log.Info().Msgf("Example CLI command: duckdb %s", dbPath)
		log.Info().Msg("Check https://github.com/davidgasquez/awesome-duckdb for more tools and clients.")
	},
}

func init() {
	ExportCmd.AddCommand(duckdbCmd)
	duckdbCmd.Flags().StringP("output", "o", "geol.duckdb", "Output DuckDB database file path")
	duckdbCmd.Flags().BoolP("force", "f", false, "Overwrites the DuckDB database file if it already exists")
}
