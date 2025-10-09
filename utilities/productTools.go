package utilities

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

// Product represents a product with its name and aliases from the API.
type Product struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}
type ProductsFile struct {
	Products map[string][]string `json:"products"`
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

func CheckCacheProductsTimeAndUpdate(cmd *cobra.Command, modTime time.Time) {
	CheckCacheTimeAndUpdateGeneric(modTime, 24*time.Hour, func() error {
		return FetchAndSaveProducts(cmd)
	})
}

func AnalyzeCacheProductsValidity(cmd *cobra.Command) {
	productsPath, err := GetProductsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving products path")
		os.Exit(1)
	}
	// Ensure cache exists, create if missing
	info, err := ensureCacheExistsGeneric(productsPath, func() error {
		return FetchAndSaveProducts(cmd)
	})
	if err != nil {
		log.Error().Err(err).Msg("Error ensuring cache exists")
		os.Exit(1)
	}

	modTime := info.ModTime()
	CheckCacheProductsTimeAndUpdate(cmd, modTime)
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
	type apiResponse struct {
		Result []Product `json:"result"`
	}
	var apiResp apiResponse

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

	// Remove the file if it exists
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
