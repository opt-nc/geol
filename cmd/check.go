package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/opt-nc/geol/v2/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("file", "f", ".geol.yaml", "File to check (default .geol.yaml)")
	checkCmd.Flags().BoolP("strict", "s", false, "Exit with error if any product is EOL")
	checkCmd.Flags().Bool("json", false, "Output in JSON format")
	checkCmd.Flags().StringP("date", "d", "", "Reference date for EOL calculations (format YYYY-MM-DD, default: today)")
}

type stackItem struct {
	Name                 string `yaml:"name"`
	Version              string `yaml:"version"`
	IdEol                string `yaml:"id_eol"`
	Skip                 bool   `yaml:"skip,omitempty"`
	ShouldAlwaysBeLatest bool   `yaml:"always-latest,omitempty"`
	ManualEol            string `yaml:"manual_eol,omitempty"`
	LtsStrategy          string `yaml:"lts_strategy,omitempty"` // "any" or "latest"
}
type geolConfig struct {
	AppName string      `yaml:"app_name"`
	Stack   []stackItem `yaml:"stack"`
}

type stackTableRow struct {
	Software      string `json:"software"`
	Version       string `json:"version"`
	EolDate       string `json:"eol_date"`
	Status        string `json:"status"`
	Days          string `json:"days"`
	IsLatest      bool   `json:"is_latest"`
	LatestVersion string `json:"latest_version"`
	LtsStrategy   string `json:"lts_strategy,omitempty"`
}

// getStackTableRows returns a slice of StackTableRow for a given stack and today date
func getStackTableRows(stack []stackItem, today time.Time) ([]stackTableRow, bool, []string) {
	rows := []stackTableRow{}
	errorOut := false
	violations := []string{}

	for _, item := range stack {
		// Skip items marked with skip: true
		if item.Skip {
			log.Info().Msgf("Found skip:true for %s %s, product will be skipped", item.Name, item.Version)
			continue
		}

		// Handle items with manual_eol set (product not in eol.date API)
		if item.ManualEol != "" {
			// Check if product exists in the API cache
			productsPath, err := utilities.GetProductsPath()
			if err == nil {
				products, err := utilities.GetProductsWithCacheRefresh(nil, productsPath)
				if err == nil {
					prod := item.IdEol
					found := false
					for name, aliases := range products.Products {
						if strings.EqualFold(prod, name) {
							found = true
							break
						}
						for _, alias := range aliases {
							if strings.EqualFold(prod, alias) {
								found = true
								break
							}
						}
						if found {
							break
						}
					}
					if found {
						log.Warn().Msgf("Product %s is available in eol.date API but has manual_eol set. Consider removing manual_eol to use official EOL data", item.Name)
					}
				}
			}

			log.Info().Msgf("Using manual EOL date for %s %s: %s (product not available in eol.date API)", item.Name, item.Version, item.ManualEol)
			eolDate := item.ManualEol
			var status string
			var daysStr string
			var daysInt int
			eolT, parseErr := time.Parse("2006-01-02", eolDate)
			if parseErr != nil {
				log.Error().Msgf("Invalid manual_eol date format for %s %s: %s (expected YYYY-MM-DD)", item.Name, item.Version, item.ManualEol)
				violations = append(violations, fmt.Sprintf("%s %s has invalid manual_eol date format: %s (expected YYYY-MM-DD)", item.Name, item.Version, item.ManualEol))
				errorOut = true
				continue
			}
			daysInt = int(eolT.Sub(today).Hours() / 24)
			daysStr = fmt.Sprintf("%d", daysInt)
			if daysInt < 0 {
				status = "EOL"
				errorOut = true
				years := -daysInt / 365
				months := (-daysInt % 365) / 30
				days := (-daysInt % 365) % 30
				log.Error().Msgf(
					"%s %s (%s) is %dy %dm %dd past EOL (manual EOL: %s)",
					item.Name, item.Version, item.Name, years, months, days, eolDate,
				)
			} else if daysInt < 30 {
				status = "WARN"
				log.Warn().Msgf(
					"%s %s (%s) is nearing EOL in %dd (manual EOL: %s)",
					item.Name, item.Version, item.Name, daysInt, eolDate,
				)
			} else {
				status = "OK"
			}
			rows = append(rows, stackTableRow{
				Software:      item.Name,
				Version:       item.Version,
				EolDate:       eolDate,
				Status:        status,
				Days:          daysStr,
				IsLatest:      false,
				LatestVersion: "-",
			})
			continue
		}

		// Handle lts_strategy enforcement
		if item.LtsStrategy != "" {
			activeLts, latestLts, ltsErr := lookupLtsInfo(item.IdEol)
			if ltsErr != nil {
				log.Error().Msgf("LTS strategy check failed for %s: %v", item.Name, ltsErr)
				violations = append(violations, fmt.Sprintf("%s: LTS strategy check failed: %v", item.Name, ltsErr))
				errorOut = true
				continue
			}
			if len(activeLts) == 0 {
				log.Error().Msgf("%s (%s): lts_strategy is set to '%s' but no active LTS versions are available for this product", item.Name, item.IdEol, item.LtsStrategy)
				violations = append(violations, fmt.Sprintf("%s (%s): lts_strategy '%s' cannot be enforced — no active LTS versions found", item.Name, item.IdEol, item.LtsStrategy))
				errorOut = true
				continue
			}

			switch item.LtsStrategy {
			case "any":
				isLts := false
				for _, lts := range activeLts {
					if lts == item.Version {
						isLts = true
						break
					}
				}
				if !isLts {
					log.Error().Msgf("%s %s: lts_strategy 'any' requires an active LTS version, but %s is not LTS (active LTS: %s)", item.Name, item.Version, item.Version, strings.Join(activeLts, ", "))
					violations = append(violations, fmt.Sprintf("%s %s is not an active LTS version (lts_strategy: any, active LTS: %s)", item.Name, item.Version, strings.Join(activeLts, ", ")))
					errorOut = true
				}
			case "latest":
				if item.Version != latestLts {
					log.Error().Msgf("%s %s: lts_strategy 'latest' requires the latest LTS version (%s), but got %s", item.Name, item.Version, latestLts, item.Version)
					violations = append(violations, fmt.Sprintf("%s %s is not the latest LTS version (lts_strategy: latest, latest LTS: %s)", item.Name, item.Version, latestLts))
					errorOut = true
				}
			}
		}

		eolDate, isLatest, latestVersion := lookupEolDate(item.IdEol, item.Version, today)
		var status string
		var daysStr string
		var daysInt int
		var eolT time.Time
		if eolDate != "" {
			eolT, _ = time.Parse("2006-01-02", eolDate)
			daysInt = int(eolT.Sub(today).Hours() / 24)
			daysStr = fmt.Sprintf("%d", daysInt)
			if daysInt < 0 {
				status = "EOL"
				errorOut = true
				// Calculate the time elapsed since EOL
				years := -daysInt / 365
				months := (-daysInt % 365) / 30
				days := (-daysInt % 365) % 30
				log.Error().Msgf(
					"%s %s (%s) is %dy %dm %dd past EOL (EOL: %s)",
					item.Name, item.Version, item.Name, years, months, days, eolDate,
				)
			} else if daysInt < 30 {
				status = "WARN"
				log.Warn().Msgf(
					"%s %s (%s) is nearing EOL in %dd (EOL: %s)",
					item.Name, item.Version, item.Name, daysInt, eolDate,
				)
			} else {
				status = "OK"
			}
		} else {
			daysStr = "-"
			status = "OK"
		}
		rows = append(rows, stackTableRow{
			Software:      item.Name,
			Version:       item.Version,
			EolDate:       eolDate,
			Status:        status,
			Days:          daysStr,
			IsLatest:      isLatest,
			LatestVersion: latestVersion,
			LtsStrategy:   item.LtsStrategy,
		})

		// Check always-latest flag
		if item.ShouldAlwaysBeLatest && !isLatest {
			violations = append(violations, fmt.Sprintf("%s %s is not the latest version (latest: %s)", item.Name, item.Version, latestVersion))
			violations = append(violations, fmt.Sprintf("%s should be in the latest version (current: %s, latest: %s)", item.Name, item.Version, latestVersion))
		}
	}
	// Sort rows by Status: EOL, WARN, OK, INFO, then by Days (from smallest to largest)
	statusOrder := map[string]int{"EOL": 0, "WARN": 1, "OK": 2}
	sort.SliceStable(rows, func(i, j int) bool {
		orderI, okI := statusOrder[rows[i].Status]
		orderJ, okJ := statusOrder[rows[j].Status]
		if !okI {
			orderI = 99
		}
		if !okJ {
			orderJ = 99
		}
		if orderI != orderJ {
			return orderI < orderJ
		}
		// If status is identical, sort by Days ascending ("-" at the end), comparing as int
		if rows[i].Days == "-" && rows[j].Days != "-" {
			return false
		}
		if rows[i].Days != "-" && rows[j].Days == "-" {
			return true
		}
		if rows[i].Days == "-" && rows[j].Days == "-" {
			return false // equal, do not change order
		}
		// Both are int, compare as int
		var di, dj int
		_, erri := fmt.Sscanf(rows[i].Days, "%d", &di)
		_, errj := fmt.Sscanf(rows[j].Days, "%d", &dj)
		if erri == nil && errj == nil {
			return di < dj
		}
		// fallback to lexicographical if problem
		return rows[i].Days < rows[j].Days
	})
	return rows, errorOut, violations
}

// lookupEolDate returns the EOL date for a given id_eol and version, along with whether the
// version is the latest cycle available as of referenceDate, and the name of that latest cycle.
// Cycles released after referenceDate are excluded so that Latest/Is Latest reflect what was
// available at the reference point in time rather than the current API snapshot.
func lookupEolDate(idEol, version string, referenceDate time.Time) (string, bool, string) {
	// Try to get products cache path
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products path")
		return "", false, ""
	}

	// Get products from cache (refresh if needed)
	products, err := utilities.GetProductsWithCacheRefresh(nil, productsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products from cache")
		return "", false, ""
	}

	prod := idEol

	found := false
	for name, aliases := range products.Products {
		if strings.EqualFold(prod, name) {
			found = true
			prod = name
			break
		}
		for _, alias := range aliases {
			if strings.EqualFold(prod, alias) {
				found = true
				prod = name
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		log.Error().Msgf("Product with id_eol %s not found in the API", idEol)
		os.Exit(1)
	}

	if len(prod) > 0 {
		url := utilities.APIUrl + "products/" + prod + "/releases/" + version
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

		url = utilities.APIUrl + "products/" + prod
		resp, err = http.Get(url)
		if err != nil {
			log.Error().Err(err).Msgf("Error requesting %s", prod)
			os.Exit(1)
		}
		body, err = io.ReadAll(resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			log.Error().Err(cerr).Msgf("Error closing HTTP body for %s", prod)
			os.Exit(1)
		}
		if err != nil {
			log.Error().Err(err).Msgf("Error reading response for %s", prod)
			os.Exit(1)
		}
		if resp.StatusCode != 200 {
			log.Error().Msgf("Product %s not found (status %d)", prod, resp.StatusCode)
			os.Exit(1)
		}
		var apiRespProd struct {
			Result struct {
				Releases []struct {
					Name        string `json:"name"`
					ReleaseDate string `json:"releaseDate"`
				} `json:"releases"`
			} `json:"result"`
		}

		if err := json.Unmarshal(body, &apiRespProd); err != nil {
			log.Error().Err(err).Msgf("Error decoding JSON for %s", prod)
			os.Exit(1)
		}

		// Determine latest cycle available as of referenceDate by excluding cycles
		// whose releaseDate is after the reference date.
		isLatest := false
		latestVersion := ""
		for _, rel := range apiRespProd.Result.Releases {
			if rel.ReleaseDate != "" {
				relDate, parseErr := time.Parse("2006-01-02", rel.ReleaseDate)
				if parseErr == nil && relDate.After(referenceDate) {
					continue
				}
			}
			// API returns releases newest-first; the first one that passes the
			// date filter is the latest cycle available at referenceDate.
			latestVersion = rel.Name
			break
		}
		if latestVersion != "" && latestVersion == version {
			isLatest = true
		}

		return apiResp.Result.EolFrom, isLatest, latestVersion
	}
	return "", false, ""
}

// lookupLtsInfo returns the currently active LTS release names (isLts=true, isEol=false) for a product,
// ordered from latest to oldest, and the name of the latest active LTS release.
// Returns an error if the product is not found or the API call fails.
func lookupLtsInfo(idEol string) ([]string, string, error) {
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving products path: %w", err)
	}
	products, err := utilities.GetProductsWithCacheRefresh(nil, productsPath)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving products from cache: %w", err)
	}

	prod := idEol
	found := false
	for name, aliases := range products.Products {
		if strings.EqualFold(prod, name) {
			found = true
			prod = name
			break
		}
		for _, alias := range aliases {
			if strings.EqualFold(prod, alias) {
				found = true
				prod = name
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return nil, "", fmt.Errorf("product with id_eol %s not found in the API", idEol)
	}

	url := utilities.APIUrl + "products/" + prod
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("error requesting %s: %w", prod, err)
	}
	body, err := io.ReadAll(resp.Body)
	if cerr := resp.Body.Close(); cerr != nil {
		return nil, "", fmt.Errorf("error closing HTTP body for %s: %w", prod, cerr)
	}
	if err != nil {
		return nil, "", fmt.Errorf("error reading response for %s: %w", prod, err)
	}
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("product %s not found (status %d)", prod, resp.StatusCode)
	}

	var apiRespProd struct {
		Result struct {
			Releases []struct {
				Name  string `json:"name"`
				IsLts bool   `json:"isLts"`
				IsEol bool   `json:"isEol"`
			} `json:"releases"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &apiRespProd); err != nil {
		return nil, "", fmt.Errorf("error decoding JSON for %s: %w", prod, err)
	}

	var activeLts []string
	for _, r := range apiRespProd.Result.Releases {
		if r.IsLts && !r.IsEol {
			activeLts = append(activeLts, r.Name)
		}
	}

	latestLts := ""
	if len(activeLts) > 0 {
		latestLts = activeLts[0]
	}
	return activeLts, latestLts, nil
}

// renderStackTable renders the stack table using lipgloss/table
func renderStackTable(rows []stackTableRow) string {
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	orange := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	t := table.New()
	t.Headers(
		"Software", "Version", "EOL Date", "Status", "Days", "Is Latest", "Latest",
	)
	for _, r := range rows {
		var daysStr string
		var statusStr string
		var latestStr string
		switch r.Status {
		case "EOL":
			statusStr = red.Render(r.Status)
			daysStr = red.Render(r.Days)
		case "OK":
			statusStr = green.Render(r.Status)
			daysStr = green.Render(r.Days)
		case "WARN":
			statusStr = orange.Render(r.Status)
			daysStr = orange.Render(r.Days)
		default:
			statusStr = r.Status
			daysStr = r.Days
		}
		if r.IsLatest {
			latestStr = green.Render("true")
		} else {
			latestStr = red.Render("false")
		}
		t.Row(
			r.Software,
			r.Version,
			r.EolDate,
			statusStr,
			daysStr,
			latestStr,
			r.LatestVersion,
		)
	}
	if term.IsTerminal(int(os.Stdout.Fd())) {
		t.Border(lipgloss.RoundedBorder())
	} else {
		t.Border(lipgloss.MarkdownBorder())
	}
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

// validationResult holds validation errors categorized by type
type validationResult struct {
	missing    []string
	duplicates []string
}

// checkRequiredKeys validates required keys in geolConfig and returns categorized errors
func checkRequiredKeys(config geolConfig) validationResult {
	result := validationResult{
		missing:    []string{},
		duplicates: []string{},
	}

	if config.AppName == "" {
		result.missing = append(result.missing, "app_name")
	}
	if len(config.Stack) == 0 {
		result.missing = append(result.missing, "stack")
	}

	// Check for duplicate names
	namesSeen := make(map[string]int)
	for i, item := range config.Stack {
		if item.Name == "" {
			result.missing = append(result.missing, fmt.Sprintf("stack[%d].name", i))
		} else {
			// Check for duplicate name
			if prevIdx, exists := namesSeen[item.Name]; exists {
				result.duplicates = append(result.duplicates, fmt.Sprintf("duplicate name '%s' at positions %d and %d", item.Name, prevIdx, i))
			}
			namesSeen[item.Name] = i
		}
		// version is always required
		if item.Version == "" {
			result.missing = append(result.missing, fmt.Sprintf("stack[%d].version", i))
		}
		if item.IdEol == "" {
			result.missing = append(result.missing, fmt.Sprintf("stack[%d].id_eol", i))
		}
		if item.LtsStrategy != "" && item.LtsStrategy != "any" && item.LtsStrategy != "latest" {
			result.missing = append(result.missing, fmt.Sprintf("stack[%d].lts_strategy must be 'any' or 'latest', got '%s'", i, item.LtsStrategy))
		}
	}
	return result
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"chk"},
	Short:   "Analyzes a stack from a YAML file, checks each component’s EOL status.",
	Long: `The 'check' command analyzes each software component listed in your stack YAML file (default: .geol.yaml), retrieves End-of-Life (EOL) information, and displays the EOL status report. Great to identify outdated software in a given stack.
Try using 'geol check init' to generate a sample stack YAML file. See https://opt-nc.github.io/geol/docs/tutorial-basics/check-command for more`,
	Example: `geol check
geol check --file stack.yaml
geol check --json`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		strict, _ := cmd.Flags().GetBool("strict")
		jsonOutput, _ := cmd.Flags().GetBool("json")
		_, err := os.Stat(file)
		if err != nil {
			log.Fatal().Msg("Error: the file does not exist: " + file)
		}

		// Read the YAML file
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal().Msg("Error reading file: " + err.Error())
		}

		var config geolConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Fatal().Msg("YAML format error: " + err.Error())
		}

		validation := checkRequiredKeys(config)
		hasErrors := false

		// Log missing fields
		if len(validation.missing) > 0 {
			for _, missing := range validation.missing {
				log.Error().Msgf("Missing or empty key: %s", missing)
			}
			hasErrors = true
		}

		// Log duplicate names
		if len(validation.duplicates) > 0 {
			for _, duplicate := range validation.duplicates {
				log.Error().Msg(duplicate)
			}
			hasErrors = true
		}

		if hasErrors {
			log.Fatal().Msg("Validation failed: please fix the errors above")
		}

		utilities.AnalyzeCacheProductsValidity(cmd)
		today := time.Now()
		if dateStr, _ := cmd.Flags().GetString("date"); dateStr != "" {
			parsed, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				log.Fatal().Msgf("Invalid --date format: %q (expected YYYY-MM-DD)", dateStr)
			}
			today = parsed
			log.Info().Msgf("Using reference date: %s", dateStr)
		}
		rows, errorOut, violations := getStackTableRows(config.Stack, today)

		if jsonOutput {
			output := struct {
				Title              string          `json:"title"`
				SoftwareComponents []stackTableRow `json:"software_components"`
			}{
				Title:              config.AppName,
				SoftwareComponents: rows,
			}
			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				log.Fatal().Msg("Error generating JSON output: " + err.Error())
			}
			fmt.Println(string(jsonData))
		} else {
			tableStr := renderStackTable(rows)
			styledTitle := lipgloss.NewStyle().
				Bold(true).Foreground(lipgloss.Color("#FFFF88")).
				Background(lipgloss.Color("#5F5FFF")).
				Render("## " + config.AppName)
			_, _ = lipgloss.Println(styledTitle)
			_, _ = lipgloss.Println(tableStr)
		}

		if len(violations) > 0 {
			for _, violation := range violations {
				log.Error().Msg(violation)
			}
		}

		if errorOut && strict {
			log.Fatal().Msg("One or more products are past EOL or not in latest version. Exiting with error due to strict mode.")
		}
	},
}
