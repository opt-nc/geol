package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/opt-nc/geol/cmd/cache"
	"github.com/opt-nc/geol/cmd/product"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geol",
	Short: "A brief description of your application", // TODO: Update this description
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`, // TODO: Update this description
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

func init() {
	rootCmd.AddCommand(cache.CacheCmd)
	rootCmd.AddCommand(product.ProductCmd)

}
