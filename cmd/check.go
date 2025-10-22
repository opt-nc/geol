package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/x/term"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("file", "f", ".geol.yaml", "File to check (default .geol.yaml)")
}

type stackItem struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	IdEol    string `yaml:"id_eol"`
	Critical bool   `yaml:"critical"`
}
type geolConfig struct {
	AppName string      `yaml:"app_name"`
	Stack   []stackItem `yaml:"stack"`
}

// checkRequiredKeys validates required keys in geolConfig and returns a slice of missing keys
func checkRequiredKeys(config geolConfig) []string {
	missing := []string{}

	if config.AppName == "" {
		missing = append(missing, "app_name")
	}
	if len(config.Stack) == 0 {
		missing = append(missing, "stack")
	}

	for i, item := range config.Stack {
		if item.Name == "" {
			missing = append(missing, fmt.Sprintf("stack[%d].name", i))
		}
		if item.Version == "" {
			missing = append(missing, fmt.Sprintf("stack[%d].version", i))
		}
		if item.IdEol == "" {
			missing = append(missing, fmt.Sprintf("stack[%d].id_eol", i))
		}
		// Check if 'critical' key is present (must be true or false, not omitted)
		if fmt.Sprintf("%v", item.Critical) != "true" && fmt.Sprintf("%v", item.Critical) != "false" {
			missing = append(missing, fmt.Sprintf("stack[%d].critical", i))
		}
	}
	return missing
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"chk"},
	Short:   "TODO",
	Long:    `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		_, err := os.Stat(file)
		if err != nil {
			log.Error().Msg("Error: the file does not exist: " + file)
			return
		}

		// Read the YAML file
		data, err := os.ReadFile(file)
		if err != nil {
			log.Error().Msg("Error reading file: " + err.Error())
			return
		}

		var config geolConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Error().Msg("YAML format error: " + err.Error())
			return
		}

		missing := checkRequiredKeys(config)
		if len(missing) > 0 {
			log.Error().Msg("Missing or empty keys: " + fmt.Sprintf("%v", missing))
			os.Exit(1)
		}

		if term.IsTerminal(os.Stdout.Fd()) { // detect if output is not a terminal
			// TODO
		} else {
			// TODO
		}
	},
}
