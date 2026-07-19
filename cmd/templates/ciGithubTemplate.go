package templates

import (
	_ "embed"
	"os"

	"github.com/phuslu/log"
)

//go:embed ci_githubTemplate.yaml
var CIGithubTemplate string

func GenerateCGithubTemplate(outputPath string, force bool) {
	if outputPath == "" {
		outputPath = ".github/workflows/geol-action.yml"
	}

	if _, err := os.Stat(outputPath); err == nil {
		if !force {
			log.Error().Msgf("the file %s already exists", outputPath)
			os.Exit(1)
		}
		log.Warn().Msgf("Overwriting existing file %s", outputPath)
	}

	log.Info().Msgf("Generating template file at %s", outputPath)
	content := CIGithubTemplate
	if err := os.WriteFile(outputPath, []byte(content), 0o644); err != nil {
		log.Error().Msgf("failed to write template file: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Template file %s generated successfully", outputPath)
	log.Info().Msg("You can now use or customize the GitHub Actions workflow " + outputPath + " to analyze your stack.")
}
