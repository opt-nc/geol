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
	Short: "Show end-of-life dates for products",
	Long:  `Efficiently display product end-of-life dates in your terminal using the endoflife.date API.`}

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
