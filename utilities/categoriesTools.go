package utilities

import (
	"os"
	"path/filepath"
)

func GetCategoriesPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	categoriesPath := filepath.Join(configDir, "geol", "categories.json")
	return categoriesPath, nil
}

// Category represents a category with its name and aliases from the API.
type Category struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}
type CategoriesFile struct {
	Categories map[string][]string `json:"categories"`
}
