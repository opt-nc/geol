package utilities

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func GetCategoriesPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	categoriesPath := filepath.Join(configDir, "geol", "categories.json")
	return categoriesPath, nil
}

type Category struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type CategoriesFile map[string]string

func FetchAndSaveCategories(cmd *cobra.Command) error {
	start := time.Now()
	categoriesPath, err := GetCategoriesPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving categories path")
		return err
	}
	// HTTP GET request (extracted)
	resp, err := GetAPIResponse(ApiUrl + "categories")
	if err != nil {
		log.Error().Err(err).Msg("Error during HTTP request")
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing response body")
		}
	}()

	//define structure to parse the response
	type apiResponse struct {
		Result []Category `json:"result"`
	}
	var apiResp apiResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		log.Error().Err(err).Msg("Error decoding JSON response")
		return err
	}

	categories := make(CategoriesFile)
	for _, category := range apiResp.Result {
		categories[category.Name] = category.Uri
	}

	// Marshal the data to JSON
	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error serializing JSON")
		return err
	}

	// Remove the file if it exists
	if err := RemoveFileIfExists(categoriesPath); err != nil {
		log.Error().Err(err).Msg("Error removing existing categories file")
		return err
	}

	// Save to file
	if err := os.WriteFile(categoriesPath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Error writing categories file")
		return err
	}

	elapsed := time.Since(start).Milliseconds()
	log.Info().Int("Number of categories", len(categories)).Int64("elapsed time (ms)", elapsed).Msg("")
	return nil
}

// GetCategoriesWithCacheRefresh tries to unmarshal categories from file, refreshes cache if needed, and returns the categories.
func GetCategoriesWithCacheRefresh(cmd *cobra.Command, categoriesPath string) (CategoriesFile, error) {
	categories := make(CategoriesFile)
	if err := readAndUnmarshalCategories(categoriesPath, &categories); err != nil {
		log.Error().Err(err).Msg("Error parsing JSON")
		log.Warn().Msg("Trying to refresh the cache now...")
		if err := FetchAndSaveCategories(cmd); err != nil {
			log.Error().Err(err).Msg("Error refreshing cache")
			return categories, err
		}
		log.Info().Msg("Cache refreshed successfully. Now getting the categories...")
		if err := readAndUnmarshalCategories(categoriesPath, &categories); err != nil {
			log.Error().Err(err).Msg("Error parsing JSON after refresh")
			return categories, err
		}
	}
	return categories, nil
}

// readAndUnmarshalCategories lit le fichier et fait l'unmarshal JSON dans categories.
func readAndUnmarshalCategories(categoriesPath string, categories *CategoriesFile) error {
	data, err := os.ReadFile(categoriesPath)
	if err != nil {
		log.Error().Err(err).Msg("Error reading categories file")
		return err
	}
	if err := json.Unmarshal(data, categories); err != nil {
		return err
	}
	return nil
}
