package templates

import (
	_ "embed"
	"os"
	"strings"

	"github.com/phuslu/log"
)

//go:embed checkTemplate.yaml
var GeolTemplate string

func GenerateTemplate(outputPath string, force bool, appName string) {
	if outputPath == "" {
		outputPath = "stack.yaml"
	}

	if _, err := os.Stat(outputPath); err == nil {
		if !force {
			log.Error().Msgf("the file %s already exists", outputPath)
			os.Exit(1)
		}
		log.Warn().Msgf("Overwriting existing file %s", outputPath)
	}

	log.Info().Msgf("Generating template file at %s", outputPath)
	content := GeolTemplate
	if appName != "" {
		content = strings.ReplaceAll(content, "MySuperApp", appName)
	}
	if err := os.WriteFile(outputPath, []byte(content), 0o644); err != nil {
		log.Error().Msgf("failed to write template file: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Template file %s generated successfully", outputPath)
	log.Info().Msg("You can now analyze your stack by running 'geol check --file " + outputPath + "'")
}
