package utilities

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

type Tag struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type TagsFile struct {
	Tags map[string]string `json:"tags"`
}

func GetTagsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	tagsPath := filepath.Join(configDir, "geol", "tags.json")
	return tagsPath, nil
}

func CheckCacheTagsTimeAndUpdate(cmd *cobra.Command, modTime time.Time) {
	CheckCacheTimeAndUpdateGeneric(modTime, 24*time.Hour, func() error {
		return FetchAndSaveTags(cmd)
	})
}

func AnalyzeCacheTagsValidity(cmd *cobra.Command) {
	tagsPath, err := GetTagsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving tags path")
		os.Exit(1)
	}
	// Ensure cache exists, create if missing
	info, err := ensureCacheExistsGeneric(tagsPath, func() error {
		return FetchAndSaveTags(cmd)
	})
	if err != nil {
		log.Error().Err(err).Msg("Error ensuring cache exists")
		os.Exit(1)
	}

	modTime := info.ModTime()
	CheckCacheTagsTimeAndUpdate(cmd, modTime)
}

func FetchAndSaveTags(cmd *cobra.Command) error {
	start := time.Now()
	tagsPath, err := GetTagsPath()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving tags path")
		return err
	}
	// HTTP GET request (extracted)
	resp, err := GetAPIResponse(ApiUrl + "tags")
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
		Result []Tag `json:"result"`
	}
	var apiResp apiResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		log.Error().Err(err).Msg("Error decoding JSON response")
		return err
	}

	tags := TagsFile{Tags: make(map[string]string)}
	for _, tag := range apiResp.Result {
		tags.Tags[tag.Name] = tag.Uri
	}

	// Marshal the data to JSON
	data, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error serializing JSON")
		return err
	}

	// Remove the file if it exists
	if err := RemoveFileIfExists(tagsPath); err != nil {
		log.Error().Err(err).Msg("Error removing existing tags file")
		return err
	}

	// Save to file
	if err := os.WriteFile(tagsPath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Error writing tags file")
		return err
	}

	elapsed := time.Since(start).Milliseconds()
	log.Info().Int("Number of tags", len(tags.Tags)).Int64("elapsed time (ms)", elapsed).Msg("")
	return nil
}
