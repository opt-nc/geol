package utilities

import (
	"os"
	"path/filepath"
)

type Tag struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type TagsFile struct {
	Tags map[string][]string `json:"tags"`
}

func GetTagsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	tagsPath := filepath.Join(configDir, "geol", "tags.json")
	return tagsPath, nil
}
