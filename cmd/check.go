package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
	"github.com/opt-nc/geol/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("file", "f", ".geol.yaml", "File to check (default .geol.yaml)")
}

type stackItem struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	IdEol    string `yaml:"id_eol"`
	Critical bool   `yaml:"critical"`
}
type geolConfig struct {
	AppName string      `yaml:"app_name"`
	Stack   []stackItem `yaml:"stack"`
}

type stackTableRow struct {
	Software string
	Version  string
	EolDate  string
	Status   string
	Critical bool
	Days     int
}

// getStackTableRows returns a slice of StackTableRow for a given stack and today date
func getStackTableRows(stack []stackItem, today time.Time) ([]stackTableRow, bool) {
	rows := []stackTableRow{}
	errorOut := false

	for _, item := range stack {
		// Lookup EOL date using id_eol and version (placeholder logic)
		eolDate := lookupEolDate(item.IdEol, item.Version)
		var status string
		var days int
		var eolT time.Time
		if eolDate != "" && item.Critical {
			eolT, _ = time.Parse("2006-01-02", eolDate)
			days = int(eolT.Sub(today).Hours() / 24)
			if days < 0 {
				status = "❌ EOL"
				errorOut = true
			} else if days < 180 {
				status = "⚠️ WARN"
			} else {
				status = "✅ OK"
			}
		} else {
			status = "➿ INFO"
		}
		rows = append(rows, stackTableRow{
			Software: item.Name,
			Version:  item.Version,
			EolDate:  eolDate,
			Status:   status,
			Critical: item.Critical,
			Days:     days,
		})
	}
	return rows, errorOut
}

// lookupEolDate should return the EOL date for a given id_eol and version (YYYY-MM-DD)
func lookupEolDate(idEol, version string) string {
	// Try to get products cache path
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products path")
		return ""
	}

	// Get products from cache (refresh if needed)
	products, err := utilities.GetProductsWithCacheRefresh(nil, productsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products from cache")
		return ""
	}

	// Find product by id_eol
	prod := products.Products[idEol]

	if len(prod) > 0 {
		url := utilities.ApiUrl + "products/" + prod[0] + "/releases/" + version
		// log.Info().Msgf("Requesting EOL date from URL: %s", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Error().Err(err).Msgf("Error requesting %s", prod)
			os.Exit(1)
		}
		body, err := io.ReadAll(resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			log.Error().Err(cerr).Msgf("Error closing HTTP body for %s", prod)
			os.Exit(1)
		}
		if err != nil {
			log.Error().Err(err).Msgf("Error reading response for %s", prod)
			os.Exit(1)
		}
		if resp.StatusCode != 200 {
			log.Error().Msgf("Product %s version %s not found (status %d)", prod, version, resp.StatusCode)
			os.Exit(1)
		}
		var apiResp struct {
			Result struct {
				Name    string `json:"name"`
				EolFrom string `json:"eolFrom"`
				IsEol   bool   `json:"isEol"`
			} `json:"result"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Error().Err(err).Msgf("Error decoding JSON for %s", prod)
			os.Exit(1)
		}
		return apiResp.Result.EolFrom
	}
	return ""
}

// renderStackTable renders the stack table using lipgloss/table
func renderStackTable(rows []stackTableRow) string {
	t := table.New()
	t.Headers(
		"SOFTWARE", "VERSION", "EOL DATE", "STATUS", "CRITICAL", "DAYS",
	)
	for _, r := range rows {
		t.Row(
			r.Software,
			r.Version,
			r.EolDate,
			r.Status,
			fmt.Sprintf("%v", r.Critical),
			fmt.Sprintf("%d", r.Days),
		)
	}
	t.Border(lipgloss.MarkdownBorder())
	t.BorderBottom(false)
	t.BorderTop(false)
	t.BorderLeft(false)
	t.BorderRight(false)
	t.BorderStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("63")))
	t.StyleFunc(func(row, col int) lipgloss.Style {
		padding := 1
		return lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Align(lipgloss.Left).Padding(0, padding)
	})
	return t.Render()
}

// checkRequiredKeys validates required keys in geolConfig and returns a slice of missing keys
func checkRequiredKeys(config geolConfig) []string {
	missing := []string{}

	if config.AppName == "" {
		missing = append(missing, "app_name")
	}
	if len(config.Stack) == 0 {
		missing = append(missing, "stack")
	}

	for i, item := range config.Stack {
		if item.Name == "" {
			missing = append(missing, fmt.Sprintf("stack[%d].name", i))
		}
		if item.Version == "" {
			missing = append(missing, fmt.Sprintf("stack[%d].version", i))
		}
		if item.IdEol == "" {
			missing = append(missing, fmt.Sprintf("stack[%d].id_eol", i))
		}
		// Check if 'critical' key is present (must be true or false, not omitted)
		if fmt.Sprintf("%v", item.Critical) != "true" && fmt.Sprintf("%v", item.Critical) != "false" {
			missing = append(missing, fmt.Sprintf("stack[%d].critical", i))
		}
	}
	return missing
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"chk"},
	Short:   "TODO",
	Long:    `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		_, err := os.Stat(file)
		if err != nil {
			log.Error().Msg("Error: the file does not exist: " + file)
			return
		}

		// Read the YAML file
		data, err := os.ReadFile(file)
		if err != nil {
			log.Error().Msg("Error reading file: " + err.Error())
			return
		}

		var config geolConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Error().Msg("YAML format error: " + err.Error())
			return
		}

		missing := checkRequiredKeys(config)
		if len(missing) > 0 {
			log.Error().Msg("Missing or empty keys: " + fmt.Sprintf("%v", missing))
			os.Exit(1)
		}

		utilities.AnalyzeCacheProductsValidity(cmd)
		today := time.Now()
		// log.Info().Msg(fmt.Sprintf("Checking stack in file: %+v", config.Stack))
		rows, errorOut := getStackTableRows(config.Stack, today)
		if term.IsTerminal(os.Stdout.Fd()) { // detect if output is a terminal
			tableStr := renderStackTable(rows)
			fmt.Println(tableStr)
		} else {
			// TODO
			log.Info().Msg("Non-terminal output not implemented yet")
		}
		if errorOut {
			os.Exit(1)
		}
	},
}
