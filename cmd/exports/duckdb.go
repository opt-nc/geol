package exports

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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
	progress      progress.Model
	totalProducts int
	processed     int
	done          bool
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
		content += pad + lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓") + " All products processed!\n"
	} else {
		content += pad + helpStyle("Inserting product details...")
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

	// Initialize the bubbletea model with progress bar
	m := model{
		progress:      progress.New(progress.WithWidth(maxWidth)),
		totalProducts: len(products.Products),
		processed:     0,
		done:          false,
	}

	// Start the TUI in a goroutine
	p := tea.NewProgram(m)

	// Process products in a goroutine
	go func() {
		for productName := range products.Products {
			prodData, err := product.FetchProductData(productName)
			if err != nil {
				log.Warn().Err(err).Msgf("Error fetching product data for %s, skipping", productName)
				p.Send(productProcessedMsg(productName))
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
					log.Error().Err(err).Msgf("Error inserting release data for %s", productName)
				}
			}
			p.Send(productProcessedMsg(productName))
		}
	}()

	// Run the program and wait for completion
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running progress display: %w", err)
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
