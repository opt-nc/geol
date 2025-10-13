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

type TagsFile map[string]string

func GetTagsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	tagsPath := filepath.Join(configDir, "geol", "tags.json")
	return tagsPath, nil
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

	tags := make(TagsFile)
	for _, tag := range apiResp.Result {
		tags[tag.Name] = tag.Uri
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
	log.Info().Int("Number of tags", len(tags)).Int64("elapsed time (ms)", elapsed).Msg("")
	return nil
}

// GetTagsWithCacheRefresh tries to unmarshal tags from file, refreshes cache if needed, and returns the tags.
func GetTagsWithCacheRefresh(cmd *cobra.Command, tagsPath string) (TagsFile, error) {
	tags := make(TagsFile)
	if err := readAndUnmarshalTags(tagsPath, &tags); err != nil {
		log.Error().Err(err).Msg("Error parsing JSON")
		log.Warn().Msg("Trying to refresh the cache now...")
		if err := FetchAndSaveTags(cmd); err != nil {
			log.Error().Err(err).Msg("Error refreshing cache")
			return tags, err
		}
		log.Info().Msg("Cache refreshed successfully. Now getting the tags...")
		if err := readAndUnmarshalTags(tagsPath, &tags); err != nil {
			log.Error().Err(err).Msg("Error parsing JSON after refresh")
			return tags, err
		}
	}
	return tags, nil
}

// readAndUnmarshalTags lit le fichier et fait l'unmarshal JSON dans tags.
func readAndUnmarshalTags(tagsPath string, tags *TagsFile) error {
	data, err := os.ReadFile(tagsPath)
	if err != nil {
		log.Error().Err(err).Msg("Error reading tags file")
		return err
	}
	if err := json.Unmarshal(data, tags); err != nil {
		return err
	}
	return nil
}
