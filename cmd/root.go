package cmd

import (
	"context"
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"
	"github.com/charmbracelet/fang"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/cmd/cache"
	"github.com/opt-nc/geol/cmd/exports"
	"github.com/opt-nc/geol/cmd/list"
	"github.com/opt-nc/geol/cmd/product"
	"github.com/opt-nc/geol/utilities"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geol",
	Short: "Show end-of-life dates for products",
	Long:  `Efficiently display product end-of-life dates in your terminal using the endoflife.date API.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		utilities.InitLogger(logLevel)
		checkGeolFile()
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help for this command")
	if err := rootCmd.PersistentFlags().MarkHidden("help"); err != nil {
		// Option 1: log or handle gracefully
		fmt.Fprintf(os.Stderr, "failed to hide help flag: %v\n", err)
		os.Exit(1)
	}

	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutVersion()); err != nil {
		os.Exit(1)
	}
}

func checkGeolFile() {
	exist, _ := os.Stat(".geol.yaml")
	if exist != nil {
		if err := validateWithCue(".geol.yaml"); err == nil {
			log.Debug().Msg("a valid .geol.yaml file exists in the current directory, run geol check to analyze it")
		} else {
			log.Debug().Str("error", err.Error()).Msg("a .geol.yaml file exists but it is not valid, use geol check init to create a new one")
		}
	}
}

func validateWithCue(yamlFile string) error {
	ctx := cuecontext.New()

	// Load the CUE schema
	cueSchemaData, err := os.ReadFile("geol_stack.cue")
	if err != nil {
		return fmt.Errorf("failed to read geol_stack.cue: %w", err)
	}

	cueSchema := ctx.CompileBytes(cueSchemaData)
	if cueSchema.Err() != nil {
		return fmt.Errorf("CUE schema compilation error: %w", cueSchema.Err())
	}

	// Convert YAML to CUE value
	yamlExpr, err := yaml.Extract(yamlFile, nil)
	if err != nil {
		return fmt.Errorf("YAML extraction error: %w", err)
	}

	yamlValue := ctx.BuildFile(yamlExpr)
	if yamlValue.Err() != nil {
		return fmt.Errorf("value construction error: %w", yamlValue.Err())
	}

	// Unify and validate
	unified := cueSchema.Unify(yamlValue)
	if err := unified.Validate(cue.Concrete(true)); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cache.CacheCmd)
	rootCmd.AddCommand(product.ProductCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(exports.ExportCmd)

	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Logging level, default info (debug, info, warn, error)")
}
