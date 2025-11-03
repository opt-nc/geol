package utilities

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

// CheckCacheTimeAndUpdateGeneric logs the cache mod time and updates the cache if older than maxAge using RefreshAllCaches.
func CheckCacheTimeAndUpdateGeneric(modTime time.Time, maxAge time.Duration, cmd *cobra.Command) {
	if modTime.Before(time.Now().Add(-maxAge)) {
		log.Warn().Msg("Cache last updated " + modTime.Format("2006-01-02 15:04:05") + ", older than 24 hours. Updating the cache...")
		RefreshAllCaches(cmd)
	}
}

// ensureCacheExistsGeneric checks if the cache file exists, creates it if missing using RefreshAllCaches, and returns its FileInfo or an error.
func EnsureCacheExistsGeneric(cachePath string, cmd *cobra.Command) (os.FileInfo, error) {
	info, err := os.Stat(cachePath)
	if err != nil {
		log.Warn().Err(err).Str("path", cachePath).Msg("cache not found, creating the cache...")
		RefreshAllCaches(cmd)
		// Try stat again after creating
		info, err = os.Stat(cachePath)
		if err != nil {
			log.Error().Err(err).Str("path", cachePath).Msg("Cache still not found after creation attempt")
			return nil, err
		}
	}
	return info, nil
}

var ApiUrl = "https://endoflife.date/api/v1/"

func InitLogger(logLevel string) {
	var level log.Level
	switch logLevel {
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	default:
		level = log.ErrorLevel
	}
	log.DefaultLogger = log.Logger{
		TimeField:  "time",
		TimeFormat: "15:04",
		Level:      level,
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

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// RefreshAllCaches runs all cache refresh functions and exits with code 1 if any fail.
func RefreshAllCaches(cmd *cobra.Command) {
	if err := FetchAndSaveProducts(cmd); err != nil {
		os.Exit(1)
	}
	if err := FetchAndSaveTags(cmd); err != nil {
		os.Exit(1)
	}
	if err := FetchAndSaveCategories(cmd); err != nil {
		os.Exit(1)
	}

	if err := CreateDoNotEditFile(); err != nil {
		log.Error().Err(err).Msg("Error creating DO_NOT_EDIT_ANYTHING file")
		os.Exit(1)
	}
}

// CreateDoNotEditFile creates a DO_NOT_EDIT_ANYTHING file in the geol config directory to warn users not to edit anything there.
func CreateDoNotEditFile() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	file := filepath.Join(configDir, "geol", "DO_NOT_EDIT_ANYTHING")
	if err := os.WriteFile(file, []byte("This directory is managed by geol. Do not edit anything here."), 0644); err != nil {
		return err
	}
	return nil
}
