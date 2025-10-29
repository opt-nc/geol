package templates

import (
	_ "embed"
	"os"

	"github.com/phuslu/log"
)

//go:embed checkTemplate.yaml
var GeolTemplate string

func GenerateTemplate(outputPath string) {
	if outputPath == "" {
		outputPath = "stack.yaml"
	}

	if _, err := os.Stat(outputPath); err == nil {
		log.Error().Msgf("the file %s already exists", outputPath)
		os.Exit(1)
	}

	log.Info().Msgf("Generating template file at %s", outputPath)
	if err := os.WriteFile(outputPath, []byte(GeolTemplate), 0644); err != nil {
		log.Error().Msgf("failed to write template file: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Template file %s generated successfully", outputPath)
}
