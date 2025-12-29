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
func createTempDetailsTable(db *sql.DB) error {
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
	defer rows.Close()

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

	// Initialize the bubbletea model with progress bar
	m := model{
		progress:          progress.New(progress.WithWidth(maxWidth)),
		totalProducts:     len(productIDs),
		processed:         0,
		done:              false,
		progressMessage:   "Fetching product details from API...",
		completionMessage: "All product details fetched!",
	}

	// Start the TUI in a goroutine
	p := tea.NewProgram(m)

	// Process products in a goroutine
	go func() {
		for _, productName := range productIDs {
			prodData, err := product.FetchProductData(productName)
			if err != nil {
				log.Warn().Err(err).Msgf("Error fetching product data for %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
				continue
			}

			// Insert each release into the details_temp table
			for _, release := range prodData.Releases {
				_, err = db.Exec(`INSERT INTO details_temp (product_id, cycle, release, latest, latest_release, eol) 
						VALUES (?, ?, ?, ?, ?, ?)`,
					prodData.Name,
					release.Name,
					release.ReleaseDate,
					release.LatestName,
					release.LatestDate,
					release.EolFrom,
				)
				if err != nil {
					log.Error().Err(err).Msgf("Error inserting release data for %s", productName)
				}
			}
			p.Send(productProcessedMsg(productName))
		}
	}()

	// Run the program and wait for completion
	if _, err := p.Run(); err != nil {
		log.Error().Err(err).Msg("Error running progress display")
		return err
	}

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
	log.Info().Msg("Created and populated \"details\" table")

	return nil
}

// createProductsTable creates the 'products' table and inserts product information
func createProductsTable(cmd *cobra.Command, db *sql.DB) error {
	// Create 'products' table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			label TEXT,
			category_id TEXT,
			uri TEXT
		)`)
	if err != nil {
		log.Error().Err(err).Msg("Error creating 'products' table")
		return err
	}

	// Get products from cache
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products path")
		return err
	}

	products, err := utilities.GetProductsWithCacheRefresh(cmd, productsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products from cache")
		return err
	}

	// Initialize the bubbletea model with progress bar
	m := model{
		progress:          progress.New(progress.WithWidth(maxWidth)),
		totalProducts:     len(products.Products),
		processed:         0,
		done:              false,
		progressMessage:   "Fetching product information from API...",
		completionMessage: "All product information fetched!",
	}

	// Start the TUI in a goroutine
	p := tea.NewProgram(m)

	// Process products in a goroutine
	go func() {
		for productName := range products.Products {
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
					Name     string `json:"name"`
					Label    string `json:"label"`
					Category string `json:"category"`
					Links    struct {
						Html string `json:"html"`
					} `json:"links"`
				} `json:"result"`
			}

			if err := json.Unmarshal(body, &apiResp); err != nil {
				log.Warn().Err(err).Msgf("Error decoding JSON for %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
				continue
			}

			// Insert product data
			_, err = db.Exec(`INSERT INTO products (id, label, category_id, uri) 
				VALUES (?, ?, ?, ?)`,
				apiResp.Result.Name,
				apiResp.Result.Label,
				apiResp.Result.Category,
				apiResp.Result.Links.Html,
			)
			if err != nil {
				log.Error().Err(err).Msgf("Error inserting product data for %s", productName)
			}

			p.Send(productProcessedMsg(productName))
		}
	}()

	// Run the program and wait for completion
	if _, err := p.Run(); err != nil {
		log.Error().Err(err).Msg("Error running progress display")
		return err
	}

	log.Info().Msg("Created and populated \"products\" table")

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

		// Create 'products' table and insert product information
		if err := createProductsTable(cmd, db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'products' table")
		}

		// Create 'details_temp' table and insert product details
		if err := createTempDetailsTable(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details_temp' table")
		}

		// Create 'details' table from 'details_temp' with proper date types
		if err := createDetailsTable(db); err != nil {
			log.Fatal().Err(err).Msg("Error creating and populating 'details' table")
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
