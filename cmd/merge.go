package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rperez95/openapi-merge/internal/config"
	"github.com/rperez95/openapi-merge/internal/merger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputFile string

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge OpenAPI specifications based on config",
	Long: `Merge multiple OpenAPI 2.0/3.0 specifications into a single OpenAPI 3.0 file.
	
The merge process:
1. Loads each input file (converting OAS 2.0 to 3.0 if needed)
2. Applies path modifications and filters
3. Handles component conflicts with dispute prefixes
4. Merges all specs into a single output file

Example:
  openapi-merge merge --config merge-config.yaml
  openapi-merge merge --config merge-config.yaml -o unified-api.json
  openapi-merge merge --config merge-config.yaml --output unified-api.yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if GetConfigFile() == "" {
			return fmt.Errorf("required flag \"config\" not set")
		}
		return nil
	},
	RunE: runMerge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Add output flag
	mergeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path (overrides config file)")
}

func runMerge(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override output if flag is provided
	if outputFile != "" {
		// Make absolute path if relative
		if !filepath.IsAbs(outputFile) {
			cwd, _ := os.Getwd()
			outputFile = filepath.Join(cwd, outputFile)
		}
		cfg.Output = outputFile
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create merger and execute
	m := merger.New(cfg, IsVerbose())

	if IsVerbose() {
		fmt.Printf("Starting merge with %d input files\n", len(cfg.Inputs))
		fmt.Printf("Output file: %s\n", cfg.Output)
	}

	if err := m.Merge(); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	fmt.Printf("Successfully merged %d specifications into %s\n", len(cfg.Inputs), cfg.Output)
	return nil
}

func loadConfig() (*config.Config, error) {
	var cfg config.Config

	// Set up decoder options to use mapstructure tags
	if err := viper.Unmarshal(&cfg, viper.DecodeHook(config.DecodeHook())); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Resolve relative paths based on config file location
	configDir := getConfigDir()
	cfg.ResolveRelativePaths(configDir)

	return &cfg, nil
}

func getConfigDir() string {
	cfgFile := GetConfigFile()
	if cfgFile == "" {
		cwd, _ := os.Getwd()
		return cwd
	}

	// Get directory from config file path
	for i := len(cfgFile) - 1; i >= 0; i-- {
		if cfgFile[i] == '/' || cfgFile[i] == '\\' {
			return cfgFile[:i]
		}
	}
	cwd, _ := os.Getwd()
	return cwd
}
