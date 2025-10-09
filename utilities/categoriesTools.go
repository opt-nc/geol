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

type CategoriesFile struct {
	Categories map[string]string `json:"categories"`
}

func CheckCacheCategoriesTimeAndUpdate(cmd *cobra.Command, modTime time.Time) {
	CheckCacheTimeAndUpdateGeneric(modTime, 24*time.Hour, func() error {
		return FetchAndSaveCategories(cmd)
	})
}

func AnalyzeCacheCategoriesValidity(cmd *cobra.Command) {
	categoriesPath, err := GetCategoriesPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving categories path")
		os.Exit(1)
	}
	// Ensure cache exists, create if missing
	info, err := ensureCacheExistsGeneric(categoriesPath, func() error {
		return FetchAndSaveCategories(cmd)
	})
	if err != nil {
		log.Error().Err(err).Msg("Error ensuring cache exists")
		os.Exit(1)
	}

	modTime := info.ModTime()
	CheckCacheCategoriesTimeAndUpdate(cmd, modTime)
}

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

	categories := CategoriesFile{Categories: make(map[string]string)}
	for _, category := range apiResp.Result {
		categories.Categories[category.Name] = category.Uri
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
	log.Info().Int("Number of categories", len(categories.Categories)).Int64("elapsed time (ms)", elapsed).Msg("")
	return nil
}
