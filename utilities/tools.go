package utilities

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/phuslu/log"
)

// ensureCacheExistsGeneric checks if the cache file exists, creates it if missing using the provided fetchAndSave function, and returns its FileInfo or an error.
func ensureCacheExistsGeneric(cachePath string, fetchAndSave func() error) (os.FileInfo, error) {
	info, err := os.Stat(cachePath)
	if err != nil {
		log.Warn().Err(err).Str("path", cachePath).Msg("cache not found, creating the cache...")
		if ferr := fetchAndSave(); ferr != nil {
			log.Error().Err(ferr).Msg("Failed to create cache")
			return nil, ferr
		}
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
