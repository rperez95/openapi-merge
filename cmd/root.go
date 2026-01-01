// Package cmd provides the CLI commands for openapi-merge.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	// Version info set by main
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// SetVersionInfo sets the version information from main
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openapi-merge",
	Short: "Merge multiple OpenAPI specifications into one",
	Long: `openapi-merge is a CLI tool that merges multiple OpenAPI 2.0 (Swagger) 
and OpenAPI 3.0/3.1 specifications into a single valid OpenAPI 3.0 file.

This is primarily used for API Gateways where multiple microservices 
need to be exposed under a single unified schema.

Example:
  openapi-merge merge --config merge-config.yaml
  openapi-merge merge --config merge-config.json`,
	Version: "dev",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	updateVersion()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (required for merge)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	// Set version template
	rootCmd.SetVersionTemplate(`{{.Name}} {{.Version}}
`)
}

// updateVersion updates the root command version string
func updateVersion() {
	if commit != "unknown" && date != "unknown" {
		rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	} else {
		rootCmd.Version = version
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// IsVerbose returns whether verbose mode is enabled.
func IsVerbose() bool {
	return verbose
}

// GetConfigFile returns the config file path.
func GetConfigFile() string {
	return cfgFile
}
