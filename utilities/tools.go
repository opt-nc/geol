package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// EnsureCacheExists checks if the cache file exists, creates it if missing, and returns its FileInfo or an error.
func EnsureCacheExists(cmd *cobra.Command, productsPath string) (os.FileInfo, error) {
	info, err := os.Stat(productsPath)
	if err != nil {
		cmd.PrintErrln("cache not found, creating the cache...")
		if ferr := FetchAndSaveProducts(cmd); ferr != nil {
			cmd.PrintErrln("Failed to create cache:", ferr)
			return nil, ferr
		}
		// Try stat again after creating
		info, err = os.Stat(productsPath)
		if err != nil {
			cmd.PrintErrln("Cache still not found after creation attempt:", productsPath)
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
	cmd.Printf("Cache last updated: %s\n", modTime.Format("2006-01-02 15:04:05"))
	// Check if the cache is older than 24 hours
	if modTime.Before(time.Now().Add(-24 * time.Hour)) {
		cmd.Printf("Warning: The cache is older than 24 hours. Updating the cache...\n")
		if err := FetchAndSaveProducts(cmd); err != nil {
			cmd.PrintErrln("Error updating cache:", err)
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
		cmd.PrintErrln("Error retrieving products path:", err)
		return err
	}

	// Ensure the directory exists
	if err := createDirectoryIfNotExists(productsPath); err != nil {
		cmd.PrintErrln("Error ensuring directory exists:", err)
		return err
	}

	// HTTP GET request (extracted)
	resp, err := GetAPIResponse(ApiUrl + "products")
	if err != nil {
		cmd.PrintErrln("Error during HTTP request:", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			cmd.PrintErrln("Error closing response body:", err)
		}
	}()

	// Define structure to parse the response
	type ApiResponse struct {
		Result []Product `json:"result"`
	}
	var apiResp ApiResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		cmd.PrintErrln("Error decoding JSON:", err)
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
		cmd.PrintErrln("Error serializing JSON:", err)
		return err
	}

	// Remove the file if it exists (extracted)
	if err := RemoveFileIfExists(productsPath); err != nil {
		cmd.PrintErrln("Error removing old file:", err)
		return err
	}
	// Save to file
	if err := os.WriteFile(productsPath, data, 0644); err != nil {
		cmd.PrintErrln("Error writing file:", err)
		return err
	}
	// Print the number of products written and elapsed time
	elapsed := time.Since(start).Milliseconds()
	cmd.Printf("%d Products retrieved from endoflife.date \n(elapsed time: %d ms)\n", len(products.Products), elapsed)
	return nil
}
