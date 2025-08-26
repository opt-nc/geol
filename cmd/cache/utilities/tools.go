package utilities

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Product represents a product with its name and aliases from the API.
type Product struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}
type ProductsFile struct {
	Products map[string][]string `json:"products"`
}

var ApiUrl = "https://endoflife.date/api/v1/"

// removeFileIfExists removes the file at the given path if it exists.
// Returns nil if the file does not exist or is successfully removed.
func removeFileIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

// getAPIResponse performs an HTTP GET request to the given URL and returns the response if status is 200.
// The caller is responsible for closing the response body.
func getAPIResponse(url string) (*http.Response, error) {
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
