package check

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/opt-nc/geol/v2/utilities"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	CheckCmd.AddCommand(InitCmd)
	CheckCmd.Flags().StringP("file", "f", ".geol.yaml", "File to check (default .geol.yaml)")
	CheckCmd.Flags().BoolP("strict", "s", false, "Exit with error if any product is EOL")
	CheckCmd.Flags().Bool("json", false, "Output in JSON format")
	CheckCmd.Flags().StringP("date", "d", "", "Reference date for EOL calculations (format YYYY-MM-DD, default: today)")
}

type stackItem struct {
	Name                 string `yaml:"name"`
	Version              string `yaml:"version"`
	IdEol                string `yaml:"id_eol"`
	Skip                 bool   `yaml:"skip,omitempty"`
	ShouldAlwaysBeLatest bool   `yaml:"always-latest,omitempty"`
	ManualEol            string `yaml:"manual_eol,omitempty"`
	LtsStrategy          string `yaml:"lts_strategy,omitempty"`   // "any" or "latest"
	LtsGraceDays         int    `yaml:"lts_grace_days,omitempty"` // grace period (days) before failing when a newer LTS exists; only applies to lts_strategy: "latest"
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
	// EolHealthScore is a 0-100 "technical debt" score computed by an EOL scoring function
	// (see standardEolScore). Named eol_health_score (rather than score) so it isn't confused
	// with the overall stack score exposed at the top level of the JSON output.
	EolHealthScore int `json:"eol_health_score"`
}

// riskThresholdDays is the number of days before EOL at which a component is considered an
// "upcoming risk". This mirrors RISK_THRESHOLD_DAYS in the geol-check-report.qmd notebook.
const riskThresholdDays = 180

// standardEolScore is geol's default, built-in EOL scoring formula, mirroring the
// compute_health_score() logic from assets/_templates/notebooks/check/geol-check-report.qmd.
// It returns a score between 0 (fully past EOL) and 100 (up to date), based on:
//   - 0 when the component is already past its EOL date
//   - 30/35/45 when the component is nearing EOL (less than riskThresholdDays remaining),
//     with higher scores awarded to LTS versions (and the highest to the latest LTS version)
//   - 100 when the component is on the latest available version (no lag)
//   - 60/75/95 when a newer major version is available ("Major Lag"), with higher scores
//     awarded to LTS versions (and the highest to the latest LTS version)
//   - 80 when only a newer minor/patch version is available ("Minor Lag")
func standardEolScore(eolDate string, referenceDate time.Time, isLatest bool, version, latestVersion string, isLts, isLatestLts bool) int {
	if eolDate != "" {
		if eolT, err := time.Parse("2006-01-02", eolDate); err == nil {
			daysUntilEol := int(eolT.Sub(referenceDate).Hours() / 24)
			switch {
			case daysUntilEol < 0:
				return 0
			case daysUntilEol < riskThresholdDays:
				switch {
				case isLts && isLatestLts:
					return 45
				case isLts:
					return 35
				default:
					return 30
				}
			}
		}
	}

	if isLatest || latestVersion == "" || version == latestVersion {
		return 100
	}

	// Component is behind the latest known version: determine whether the lag is a major
	// version bump ("Major Lag") or just a minor/patch bump ("Minor Lag").
	isMajorLag := false
	v, errV := semver.NewVersion(version)
	lv, errL := semver.NewVersion(latestVersion)
	if errV == nil && errL == nil {
		isMajorLag = v.Major() != lv.Major()
	} else {
		// Fall back to a naive string comparison, matching the qmd notebook's fallback.
		isMajorLag = strings.SplitN(version, ".", 2)[0] != strings.SplitN(latestVersion, ".", 2)[0]
	}
	if isMajorLag {
		switch {
		case isLts && isLatestLts:
			return 95
		case isLts:
			return 75
		default:
			return 60
		}
	}
	return 80
}

// isVersionLts reports whether version matches (or is a sub-version of, e.g. "24.04.1" for
// cycle "24.04") one of the given active LTS cycles, and whether it matches the latest LTS cycle.
func isVersionLts(version string, activeLts []string, latestLts string) (isLts, isLatestLts bool) {
	for _, cycle := range activeLts {
		if version == cycle || strings.HasPrefix(version, cycle+".") {
			isLts = true
			break
		}
	}
	if isLts && latestLts != "" && (version == latestLts || strings.HasPrefix(version, latestLts+".")) {
		isLatestLts = true
	}
	return isLts, isLatestLts
}

// stackScore holds the overall stack debt score, along with a color and a human-readable
// message meant for display purposes (e.g. JSON "score" field, dashboards, badges).
type stackScore struct {
	Value   int    `json:"value"`
	Color   string `json:"color"`
	Message string `json:"message"`
}

// computeStackScore returns the average debt score across all scored components, along
// with a color/message pair summarizing the overall stack health.
func computeStackScore(rows []stackTableRow) stackScore {
	if len(rows) == 0 {
		return stackScore{Value: 100, Color: "green", Message: "Healthy — No software components to evaluate"}
	}

	total := 0
	for _, r := range rows {
		total += r.EolHealthScore
	}
	avg := int(math.Round(float64(total) / float64(len(rows))))

	switch {
	case avg >= 80:
		return stackScore{Value: avg, Color: "green", Message: "Healthy — All software components are up to date"}
	case avg >= 50:
		return stackScore{Value: avg, Color: "orange", Message: "Needs Attention — Some software components are not up to date"}
	default:
		return stackScore{Value: avg, Color: "red", Message: "Critical — Several software components are past end-of-life or severely outdated"}
	}
}

// renderScoreValue colorizes a per-component debt score for terminal/markdown table display.
func renderScoreValue(value int) string {
	switch {
	case value >= 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(fmt.Sprintf("%d", value))
	case value >= 50:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render(fmt.Sprintf("%d", value))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("%d", value))
	}
}

// renderStackScore renders a one-line summary of the overall stack debt score.
func renderStackScore(score stackScore) string {
	colorCode := map[string]string{"green": "46", "orange": "208", "red": "196"}[score.Color]
	if colorCode == "" {
		colorCode = "252"
	}
	valueStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorCode)).Render(fmt.Sprintf("%d/100", score.Value))
	return fmt.Sprintf("Stack Debt Score: %s — %s", valueStr, score.Message)
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
			manualIsLts, manualIsLatestLts := false, false
			if activeLts, latestLts, _, ltsErr := lookupLtsInfo(item.IdEol); ltsErr == nil {
				manualIsLts, manualIsLatestLts = isVersionLts(item.Version, activeLts, latestLts)
			}
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
				Software:       item.Name,
				Version:        item.Version,
				EolDate:        eolDate,
				Status:         status,
				Days:           daysStr,
				IsLatest:       false,
				LatestVersion:  "-",
				EolHealthScore: standardEolScore(eolDate, today, false, item.Version, "", manualIsLts, manualIsLatestLts),
			})
			continue
		}

		// Handle lts_strategy enforcement
		var lookedUpActiveLts []string
		var lookedUpLatestLts string
		if item.LtsStrategy != "" {
			activeLts, latestLts, latestLtsDate, ltsErr := lookupLtsInfo(item.IdEol)
			if ltsErr != nil {
				log.Fatal().Msgf("LTS strategy check failed for %s: %v", item.Name, ltsErr)
			}
			lookedUpActiveLts, lookedUpLatestLts = activeLts, latestLts
			if len(activeLts) == 0 {
				log.Fatal().Msgf("%s (%s): lts_strategy is set to '%s' but no active LTS versions are available for this product", item.Name, item.IdEol, item.LtsStrategy)
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
					log.Fatal().Msgf("%s %s: lts_strategy 'any' requires an active LTS version, but %s is not LTS (active LTS: %s)", item.Name, item.Version, item.Version, strings.Join(activeLts, ", "))
				}
			case "latest":
				if item.Version != latestLts {
					// Apply grace period: if lts_grace_days > 0 and the latest LTS was released
					// less than lts_grace_days days ago, warn instead of failing.
					withinGrace := false
					if item.LtsGraceDays > 0 && latestLtsDate != "" {
						if ltsRelDate, parseErr := time.Parse("2006-01-02", latestLtsDate); parseErr == nil {
							daysSinceLatestLts := int(today.Sub(ltsRelDate).Hours() / 24)
							if daysSinceLatestLts < item.LtsGraceDays {
								withinGrace = true
								log.Warn().Msgf(
									"%s %s: lts_strategy 'latest' — newer LTS %s was released %dd ago (grace period: %dd). Update before grace period expires.",
									item.Name, item.Version, latestLts, daysSinceLatestLts, item.LtsGraceDays,
								)
							}
						}
					}
					if !withinGrace {
						log.Error().Msgf("%s %s: lts_strategy 'latest' requires the latest LTS version (%s), but got %s", item.Name, item.Version, latestLts, item.Version)
						violations = append(violations, fmt.Sprintf("%s %s is not the latest LTS version (lts_strategy: latest, latest LTS: %s)", item.Name, item.Version, latestLts))
						errorOut = true
					}
				}
			}
		}

		eolDate, isLatest, latestVersion, eolLookupErr := lookupEolDate(item.IdEol, item.Version, today)
		if eolLookupErr != nil {
			log.Fatal().Msgf("%s %s: %v", item.Name, item.Version, eolLookupErr)
		}

		// Determine LTS status for score computation, reusing the lookup already performed
		// above for lts_strategy enforcement when available, to avoid a duplicate API call.
		activeLtsForScore, latestLtsForScore := lookedUpActiveLts, lookedUpLatestLts
		if item.LtsStrategy == "" {
			if fetchedActiveLts, fetchedLatestLts, _, ltsErr := lookupLtsInfo(item.IdEol); ltsErr == nil {
				activeLtsForScore, latestLtsForScore = fetchedActiveLts, fetchedLatestLts
			}
		}
		isLts, isLatestLts := isVersionLts(item.Version, activeLtsForScore, latestLtsForScore)
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
			Software:       item.Name,
			Version:        item.Version,
			EolDate:        eolDate,
			Status:         status,
			Days:           daysStr,
			IsLatest:       isLatest,
			LatestVersion:  latestVersion,
			LtsStrategy:    item.LtsStrategy,
			EolHealthScore: standardEolScore(eolDate, today, isLatest, item.Version, latestVersion, isLts, isLatestLts),
		})

		// Check always-latest flag
		if item.ShouldAlwaysBeLatest && !isLatest {
			violations = append(violations, fmt.Sprintf("%s %s is not the latest version (latest: %s)", item.Name, item.Version, latestVersion))
			violations = append(violations, fmt.Sprintf("%s should be in the latest version (current: %s, latest: %s)", item.Name, item.Version, latestVersion))
			errorOut = true
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

// findVersionSuggestion fetches all releases for a product and uses semver to suggest
// a valid release name that best matches the given version (by major.minor, then major).
// Returns an empty string if no suggestion is found or the version is already valid.
func findVersionSuggestion(prod, version string) string {
	v, err := semver.NewVersion(version)
	if err != nil {
		return ""
	}

	url := utilities.APIUrl + "products/" + prod
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if cerr := resp.Body.Close(); cerr != nil || err != nil || resp.StatusCode != 200 {
		return ""
	}

	var apiRespProd struct {
		Result struct {
			Releases []struct {
				Name string `json:"name"`
			} `json:"releases"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &apiRespProd); err != nil {
		return ""
	}

	releases := apiRespProd.Result.Releases

	// First pass: match on major.minor
	for _, rel := range releases {
		rv, err := semver.NewVersion(rel.Name)
		if err != nil {
			continue
		}
		if rv.Major() == v.Major() && rv.Minor() == v.Minor() {
			return rel.Name
		}
	}
	// Second pass: match on major only
	for _, rel := range releases {
		rv, err := semver.NewVersion(rel.Name)
		if err != nil {
			continue
		}
		if rv.Major() == v.Major() {
			return rel.Name
		}
	}
	return ""
}

// lookupEolDate returns the EOL date for a given id_eol and version, along with whether the
// version is the latest cycle available as of referenceDate, and the name of that latest cycle.
// Cycles released after referenceDate are excluded so that Latest/Is Latest reflect what was
// available at the reference point in time rather than the current API snapshot.
func lookupEolDate(idEol, version string, referenceDate time.Time) (string, bool, string, error) {
	// Try to get products cache path
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		return "", false, "", fmt.Errorf("error retrieving products path: %w", err)
	}

	// Get products from cache (refresh if needed)
	products, err := utilities.GetProductsWithCacheRefresh(nil, productsPath)
	if err != nil {
		return "", false, "", fmt.Errorf("error retrieving products from cache: %w", err)
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
		return "", false, "", fmt.Errorf("product with id_eol %s not found in the API", idEol)
	}

	if len(prod) > 0 {
		url := utilities.APIUrl + "products/" + prod + "/releases/" + version
		resp, err := http.Get(url)
		if err != nil {
			return "", false, "", fmt.Errorf("error requesting %s: %w", prod, err)
		}
		body, err := io.ReadAll(resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			return "", false, "", fmt.Errorf("error closing HTTP body for %s: %w", prod, cerr)
		}
		if err != nil {
			return "", false, "", fmt.Errorf("error reading response for %s: %w", prod, err)
		}
		if resp.StatusCode != 200 {
			if suggestion := findVersionSuggestion(prod, version); suggestion != "" {
				log.Info().Msgf("Version %q not found for product %q in endoflife.date. Did you mean %q? Consider updating your .geol.yaml to: version: \"%s\"", version, prod, suggestion, suggestion)
			}
			return "", false, "", fmt.Errorf("product %s version %s not found (status %d)", prod, version, resp.StatusCode)
		}
		var apiResp struct {
			Result struct {
				Name    string `json:"name"`
				EolFrom string `json:"eolFrom"`
				IsEol   bool   `json:"isEol"`
			} `json:"result"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			return "", false, "", fmt.Errorf("error decoding JSON for %s: %w", prod, err)
		}

		url = utilities.APIUrl + "products/" + prod
		resp, err = http.Get(url)
		if err != nil {
			return "", false, "", fmt.Errorf("error requesting %s: %w", prod, err)
		}
		body, err = io.ReadAll(resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			return "", false, "", fmt.Errorf("error closing HTTP body for %s: %w", prod, cerr)
		}
		if err != nil {
			return "", false, "", fmt.Errorf("error reading response for %s: %w", prod, err)
		}
		if resp.StatusCode != 200 {
			return "", false, "", fmt.Errorf("product %s not found (status %d)", prod, resp.StatusCode)
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
			return "", false, "", fmt.Errorf("error decoding JSON for %s: %w", prod, err)
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

		return apiResp.Result.EolFrom, isLatest, latestVersion, nil
	}
	return "", false, "", nil
}

// lookupLtsInfo returns the currently active LTS release names (isLts=true, isEol=false) for a product,
// ordered from latest to oldest, the name of the latest active LTS release, and its release date.
// Returns an error if the product is not found or the API call fails.
// lookupLtsInfo returns:
//   - activeLts: slice of active LTS release names (isLts=true, isEol=false), newest first
//   - latestLts: name of the latest active LTS release
//   - latestLtsReleaseDate: release date of the latest active LTS (YYYY-MM-DD), empty if unknown
//   - error
func lookupLtsInfo(idEol string) ([]string, string, string, error) {
	productsPath, err := utilities.GetProductsPath()
	if err != nil {
		return nil, "", "", fmt.Errorf("error retrieving products path: %w", err)
	}
	products, err := utilities.GetProductsWithCacheRefresh(nil, productsPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("error retrieving products from cache: %w", err)
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
		return nil, "", "", fmt.Errorf("product with id_eol %s not found in the API", idEol)
	}

	url := utilities.APIUrl + "products/" + prod
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", "", fmt.Errorf("error requesting %s: %w", prod, err)
	}
	body, err := io.ReadAll(resp.Body)
	if cerr := resp.Body.Close(); cerr != nil {
		return nil, "", "", fmt.Errorf("error closing HTTP body for %s: %w", prod, cerr)
	}
	if err != nil {
		return nil, "", "", fmt.Errorf("error reading response for %s: %w", prod, err)
	}
	if resp.StatusCode != 200 {
		return nil, "", "", fmt.Errorf("product %s not found (status %d)", prod, resp.StatusCode)
	}

	var apiRespProd struct {
		Result struct {
			Releases []struct {
				Name        string `json:"name"`
				ReleaseDate string `json:"releaseDate"`
				IsLts       bool   `json:"isLts"`
				IsEol       bool   `json:"isEol"`
			} `json:"releases"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &apiRespProd); err != nil {
		return nil, "", "", fmt.Errorf("error decoding JSON for %s: %w", prod, err)
	}

	var activeLts []string
	latestLtsDate := ""
	for _, r := range apiRespProd.Result.Releases {
		if r.IsLts && !r.IsEol {
			activeLts = append(activeLts, r.Name)
			if latestLtsDate == "" {
				latestLtsDate = r.ReleaseDate
			}
		}
	}

	latestLts := ""
	if len(activeLts) > 0 {
		latestLts = activeLts[0]
	}
	return activeLts, latestLts, latestLtsDate, nil
}

// renderStackTable renders the stack table using lipgloss/table
func renderStackTable(rows []stackTableRow) string {
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	orange := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	t := table.New()
	t.Headers(
		"Software", "Version", "EOL Date", "Status", "Days", "Is Latest", "Latest", "EOL Health Score",
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
			renderScoreValue(r.EolHealthScore),
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
	constraint []string
}

// checkRequiredKeys validates required keys in geolConfig and returns categorized errors
func checkRequiredKeys(config geolConfig) validationResult {
	result := validationResult{
		missing:    []string{},
		duplicates: []string{},
		constraint: []string{},
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
		if item.LtsGraceDays != 0 && item.LtsStrategy != "latest" {
			result.missing = append(result.missing, fmt.Sprintf("stack[%d].lts_grace_days is only applicable when lts_strategy is 'latest'", i))
		}
		if item.LtsGraceDays < 0 {
			result.missing = append(result.missing, fmt.Sprintf("stack[%d].lts_grace_days must be >= 0, got %d", i, item.LtsGraceDays))
		}
		if item.ShouldAlwaysBeLatest && item.LtsStrategy != "" {
			result.constraint = append(result.constraint, fmt.Sprintf("stack[%d] cannot define both always-latest and lts_strategy for the same product", i))
		}
	}
	return result
}

// checkCmd represents the check command
var CheckCmd = &cobra.Command{
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

		// Log constraint errors
		if len(validation.constraint) > 0 {
			for _, constraint := range validation.constraint {
				log.Error().Msgf("Constraint error: %s", constraint)
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
		score := computeStackScore(rows)

		if jsonOutput {
			output := struct {
				Title              string          `json:"title"`
				Score              []stackScore    `json:"score"`
				SoftwareComponents []stackTableRow `json:"software_components"`
			}{
				Title:              config.AppName,
				Score:              []stackScore{score},
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
			_, _ = lipgloss.Println(renderStackScore(score))
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
