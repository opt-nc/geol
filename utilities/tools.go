package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func InitLogger() {
	log.DefaultLogger = log.Logger{
		TimeField:  "time",
		TimeFormat: "15:04",
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}

// TodayDateString retourne la date du jour au format YYYY-MM-DD
func TodayDateString() string {
	return time.Now().Format("2006-01-02")
}

// EnsureCacheExists checks if the cache file exists, creates it if missing, and returns its FileInfo or an error.
func EnsureCacheExists(cmd *cobra.Command, productsPath string) (os.FileInfo, error) {
	info, err := os.Stat(productsPath)
	if err != nil {
		log.Warn().Err(err).Msg("cache not found, creating the cache...")
		if ferr := FetchAndSaveProducts(cmd); ferr != nil {
			log.Error().Err(ferr).Msg("Failed to create cache")
			return nil, ferr
		}
		// Try stat again after creating
		info, err = os.Stat(productsPath)
		if err != nil {
			log.Error().Err(err).Str("path", productsPath).Msg("Cache still not found after creation attempt")
			return nil, err
		}
	}
	return info, nil
}

// Product represents a product with its name and aliases from the API.
type Product struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}
type ProductsFile struct {
	Products map[string][]string `json:"products"`
}

var ApiUrl = "https://endoflife.date/api/v1/"

// RemoveFileIfExists removes the file at the given path if it exists.
// Returns nil if the file does not exist or is successfully removed.
func RemoveFileIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

// GetAPIResponse performs an HTTP GET request to the given URL and returns the response if status is 200.
// The caller is responsible for closing the response body.
func GetAPIResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}
	if resp.StatusCode != 200 {
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("unexpected HTTP status: %s (error closing body: %w)", resp.Status, err)
		}
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	return resp, nil
}

// GetProductsPath returns the path to the products.json file in the user's config directory.
func GetProductsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	productsPath := filepath.Join(configDir, "geol", "products.json")
	return productsPath, nil
}

func CheckCacheTimeAndUpdate(cmd *cobra.Command, modTime time.Time) {
	log.Info().Msg("Cache last updated " + modTime.Format("2006-01-02 15:04:05"))
	// Check if the cache is older than 24 hours
	if modTime.Before(time.Now().Add(-24 * time.Hour)) {
		log.Warn().Msg("The cache is older than 24 hours. Updating the cache...")
		if err := FetchAndSaveProducts(cmd); err != nil {
			log.Error().Err(err).Msg("Error updating cache")
		}
	}
}

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func FetchAndSaveProducts(cmd *cobra.Command) error {
	start := time.Now()
	productsPath, err := GetProductsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products path")
		return err
	}

	// Ensure the directory exists
	if err := createDirectoryIfNotExists(productsPath); err != nil {
		log.Error().Err(err).Msg("Error ensuring directory exists")
		return err
	}

	// HTTP GET request (extracted)
	resp, err := GetAPIResponse(ApiUrl + "products")
	if err != nil {
		log.Error().Err(err).Msg("Error during HTTP request")
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing response body")
		}
	}()

	// Define structure to parse the response
	type ApiResponse struct {
		Result []Product `json:"result"`
	}
	var apiResp ApiResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		log.Error().Err(err).Msg("Error decoding JSON")
		return err
	}

	products := ProductsFile{Products: make(map[string][]string)}
	for _, p := range apiResp.Result {
		aliases := []string{p.Name}
		aliases = append(aliases, p.Aliases...)
		products.Products[p.Name] = aliases
	}

	// Marshal the data to JSON
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error serializing JSON")
		return err
	}

	// Remove the file if it exists (extracted)
	if err := RemoveFileIfExists(productsPath); err != nil {
		log.Error().Err(err).Msg("Error removing old file")
		return err
	}
	// Save to file
	if err := os.WriteFile(productsPath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Error writing file")
		return err
	}
	// Print the number of products written and elapsed time
	elapsed := time.Since(start).Milliseconds()
	//log.Info().Msg(fmt.Sprintf("%d Products retrieved from endoflife.date \n(elapsed time: %d ms)\n", len(products.Products), elapsed))
	log.Info().Int("Number of products", len(products.Products)).Int64("elapsed time (ms)", elapsed).Msg("")
	return nil
}

// GetProductsWithCacheRefresh tries to unmarshal products from file, refreshes cache if needed, and returns the products.
func GetProductsWithCacheRefresh(cmd *cobra.Command, productsPath string) (ProductsFile, error) {
	var products ProductsFile
	if err := readAndUnmarshalProducts(productsPath, &products); err != nil {
		log.Error().Err(err).Msg("Error parsing JSON")
		log.Warn().Msg("Trying to refresh the cache now...")
		if err := FetchAndSaveProducts(cmd); err != nil {
			log.Error().Err(err).Msg("Error refreshing cache")
			return products, err
		}
		log.Info().Msg("Cache refreshed successfully. Now getting the products...")
		if err := readAndUnmarshalProducts(productsPath, &products); err != nil {
			log.Error().Err(err).Msg("Error parsing JSON after refresh")
			return products, err
		}
	}
	return products, nil
}

// readAndUnmarshalProducts lit le fichier et fait l'unmarshal JSON dans products.
func readAndUnmarshalProducts(productsPath string, products *ProductsFile) error {
	data, err := os.ReadFile(productsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error reading products file")
		return err
	}
	if err := json.Unmarshal(data, products); err != nil {
		return err
	}
	return nil
}
