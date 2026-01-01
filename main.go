// Package main is the entry point for openapi-merge CLI tool.
// This tool merges multiple OpenAPI 2.0/3.0 specifications into a single OpenAPI 3.0 file.
package main

import (
	"runtime/debug"

	"github.com/rperez95/openapi-merge/cmd"
)

// Build information, set via ldflags (or read from build info)
var (
	version = ""
	commit  = ""
	date    = ""
)

func main() {
	// If version not set via ldflags, try to get from build info
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if len(setting.Value) >= 7 {
						commit = setting.Value[:7]
					} else {
						commit = setting.Value
					}
				case "vcs.time":
					date = setting.Value
				}
			}
		}
	}

	// Fallback defaults
	if version == "" {
		version = "dev"
	}
	if commit == "" {
		commit = "unknown"
	}
	if date == "" {
		date = "unknown"
	}

	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
