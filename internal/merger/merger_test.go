package merger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rperez95/openapi-merge/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerger_BasicMerge(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "openapi-merge-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	// Create test OpenAPI files
	spec1 := `{
		"openapi": "3.0.0",
		"info": {
			"title": "API 1",
			"version": "1.0.0"
		},
		"paths": {
			"/users": {
				"get": {
					"summary": "Get users",
					"tags": ["Users"],
					"responses": {
						"200": {
							"description": "Success"
						}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"User": {
					"type": "object",
					"properties": {
						"id": {"type": "string"},
						"name": {"type": "string"}
					}
				}
			}
		}
	}`

	spec2 := `{
		"openapi": "3.0.0",
		"info": {
			"title": "API 2",
			"version": "1.0.0"
		},
		"paths": {
			"/products": {
				"get": {
					"summary": "Get products",
					"tags": ["Products"],
					"responses": {
						"200": {
							"description": "Success"
						}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"Product": {
					"type": "object",
					"properties": {
						"id": {"type": "string"},
						"name": {"type": "string"}
					}
				}
			}
		}
	}`

	// Write test files
	spec1Path := filepath.Join(tempDir, "spec1.json")
	spec2Path := filepath.Join(tempDir, "spec2.json")
	outputPath := filepath.Join(tempDir, "merged.json")

	require.NoError(t, os.WriteFile(spec1Path, []byte(spec1), 0644))
	require.NoError(t, os.WriteFile(spec2Path, []byte(spec2), 0644))

	// Create config
	cfg := &config.Config{
		Inputs: []config.InputConfig{
			{InputFile: spec1Path},
			{InputFile: spec2Path},
		},
		Output: outputPath,
		Info: &config.InfoConfig{
			Title:       "Merged API",
			Description: "A merged API specification",
			Version:     "1.0.0",
		},
	}

	// Run merge
	m := New(cfg, false)
	err = m.Merge()
	require.NoError(t, err)

	// Verify output exists
	_, err = os.Stat(outputPath)
	require.NoError(t, err)

	// Read and verify output
	outputData, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(outputData), "/users")
	assert.Contains(t, string(outputData), "/products")
	assert.Contains(t, string(outputData), "User")
	assert.Contains(t, string(outputData), "Product")
}

func TestMerger_PathModification(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "openapi-merge-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	spec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "API",
			"version": "1.0.0"
		},
		"paths": {
			"/v1/users": {
				"get": {
					"summary": "Get users",
					"responses": {
						"200": {"description": "Success"}
					}
				}
			}
		}
	}`

	specPath := filepath.Join(tempDir, "spec.json")
	outputPath := filepath.Join(tempDir, "merged.json")

	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))

	cfg := &config.Config{
		Inputs: []config.InputConfig{
			{
				InputFile: specPath,
				PathModification: &config.PathModificationConfig{
					StripStart: "/v1",
					Prepend:    "/api",
				},
			},
		},
		Output: outputPath,
	}

	m := New(cfg, false)
	err = m.Merge()
	require.NoError(t, err)

	outputData, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(outputData), "/api/users")
	assert.NotContains(t, string(outputData), "/v1/users")
}

func TestMerger_DisputePrefix(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "openapi-merge-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	spec1 := `{
		"openapi": "3.0.0",
		"info": {"title": "API 1", "version": "1.0.0"},
		"paths": {
			"/items": {
				"get": {
					"summary": "Get items",
					"responses": {
						"200": {
							"description": "Success",
							"content": {
								"application/json": {
									"schema": {"$ref": "#/components/schemas/Item"}
								}
							}
						}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"Item": {
					"type": "object",
					"properties": {"id": {"type": "string"}}
				}
			}
		}
	}`

	spec2 := `{
		"openapi": "3.0.0",
		"info": {"title": "API 2", "version": "1.0.0"},
		"paths": {
			"/other-items": {
				"get": {
					"summary": "Get other items",
					"responses": {
						"200": {
							"description": "Success",
							"content": {
								"application/json": {
									"schema": {"$ref": "#/components/schemas/Item"}
								}
							}
						}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"Item": {
					"type": "object",
					"properties": {"name": {"type": "string"}}
				}
			}
		}
	}`

	spec1Path := filepath.Join(tempDir, "spec1.json")
	spec2Path := filepath.Join(tempDir, "spec2.json")
	outputPath := filepath.Join(tempDir, "merged.json")

	require.NoError(t, os.WriteFile(spec1Path, []byte(spec1), 0644))
	require.NoError(t, os.WriteFile(spec2Path, []byte(spec2), 0644))

	cfg := &config.Config{
		Inputs: []config.InputConfig{
			{InputFile: spec1Path},
			{
				InputFile: spec2Path,
				Dispute:   &config.DisputeConfig{Prefix: "API2_"},
			},
		},
		Output: outputPath,
	}

	m := New(cfg, false)
	err = m.Merge()
	require.NoError(t, err)

	outputData, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(outputData), "Item")
	assert.Contains(t, string(outputData), "API2_Item")
}

func TestMerger_OperationSelection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "openapi-merge-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	spec := `{
		"openapi": "3.0.0",
		"info": {"title": "API", "version": "1.0.0"},
		"paths": {
			"/users": {
				"get": {
					"summary": "Get users",
					"tags": ["Users"],
					"responses": {"200": {"description": "Success"}}
				}
			},
			"/admin": {
				"get": {
					"summary": "Admin endpoint",
					"tags": ["Admin"],
					"responses": {"200": {"description": "Success"}}
				}
			}
		}
	}`

	specPath := filepath.Join(tempDir, "spec.json")
	outputPath := filepath.Join(tempDir, "merged.json")

	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))

	cfg := &config.Config{
		Inputs: []config.InputConfig{
			{
				InputFile: specPath,
				OperationSelection: &config.OperationSelectionConfig{
					ExcludeTags: []string{"Admin"},
				},
			},
		},
		Output: outputPath,
	}

	m := New(cfg, false)
	err = m.Merge()
	require.NoError(t, err)

	outputData, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(outputData), "/users")
	assert.NotContains(t, string(outputData), "/admin")
}

func TestMatchGlob(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"/api/*", "/api/users", true},
		{"/api/*", "/api/users/123", true}, // gobwas/glob * matches any characters including /
		{"/api/**", "/api/users/123", true},
		{"/api/*/items", "/api/v1/items", true},
		{"/users", "/users", true},
		{"/users", "/products", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.path, func(t *testing.T) {
			got := matchGlob(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &config.Config{
				Inputs: []config.InputConfig{{InputFile: "test.json"}},
				Output: "output.json",
			},
			wantErr: false,
		},
		{
			name: "missing inputs",
			cfg: &config.Config{
				Inputs: []config.InputConfig{},
				Output: "output.json",
			},
			wantErr: true,
		},
		{
			name: "missing output",
			cfg: &config.Config{
				Inputs: []config.InputConfig{{InputFile: "test.json"}},
				Output: "",
			},
			wantErr: true,
		},
		{
			name: "missing inputFile",
			cfg: &config.Config{
				Inputs: []config.InputConfig{{InputFile: ""}},
				Output: "output.json",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
