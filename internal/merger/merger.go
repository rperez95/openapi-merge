// Package merger provides the core logic for merging OpenAPI specifications.
package merger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rperez95/openapi-merge/internal/config"
	"gopkg.in/yaml.v3"
)

// Merger handles the merging of OpenAPI specifications.
type Merger struct {
	cfg     *config.Config
	verbose bool
	master  *openapi3.T
}

// New creates a new Merger instance.
func New(cfg *config.Config, verbose bool) *Merger {
	return &Merger{
		cfg:     cfg,
		verbose: verbose,
	}
}

// Merge executes the merge operation.
func (m *Merger) Merge() error {
	// Initialize master spec
	m.master = &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       "Merged API",
			Description: "",
			Version:     "1.0.0",
		},
		Paths: &openapi3.Paths{
			Extensions: make(map[string]interface{}),
		},
		Components: &openapi3.Components{
			Schemas:         make(openapi3.Schemas),
			Parameters:      make(openapi3.ParametersMap),
			Headers:         make(openapi3.Headers),
			RequestBodies:   make(openapi3.RequestBodies),
			Responses:       make(openapi3.ResponseBodies),
			SecuritySchemes: make(openapi3.SecuritySchemes),
			Examples:        make(openapi3.Examples),
			Links:           make(openapi3.Links),
			Callbacks:       make(openapi3.Callbacks),
		},
		Tags: make(openapi3.Tags, 0),
	}

	// Track merged descriptions for appending
	var mergedDescriptions []string

	// Process each input file
	for i, input := range m.cfg.Inputs {
		if m.verbose {
			fmt.Printf("Processing input %d: %s\n", i+1, input.InputFile)
		}

		// Load and parse the spec
		spec, err := m.loadSpec(input.InputFile)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", input.InputFile, err)
		}

		// Apply operation selection filters
		spec = m.filterOperations(spec, &input)

		// Apply path modifications
		spec = m.modifyPaths(spec, &input)

		// Apply parameter modifications
		spec = m.modifyParameters(spec, &input)

		// Handle conflicts with dispute prefix
		if input.Dispute != nil && input.Dispute.Prefix != "" {
			spec = m.applyDisputePrefix(spec, input.Dispute.Prefix)
		}

		// Merge into master
		if err := m.mergeSpec(spec, &input); err != nil {
			return fmt.Errorf("failed to merge %s: %w", input.InputFile, err)
		}

		// Handle description appending
		if input.Description != nil && input.Description.Append && spec.Info != nil {
			desc := m.formatDescription(spec.Info.Description, input.Description)
			if desc != "" {
				mergedDescriptions = append(mergedDescriptions, desc)
			}
		}
	}

	// Apply post-processing
	m.applyOverrides(mergedDescriptions)
	m.sortOutput()

	// Write output
	return m.writeOutput()
}

// loadSpec loads and parses an OpenAPI specification, converting OAS2 to OAS3 if needed.
// Supports both local files and HTTP/HTTPS URLs.
func (m *Merger) loadSpec(filePath string) (*openapi3.T, error) {
	var data []byte
	var err error
	var ext string

	if config.IsURL(filePath) {
		data, ext, err = m.fetchFromURL(filePath)
	} else {
		data, err = os.ReadFile(filePath)
		ext = strings.ToLower(filepath.Ext(filePath))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detect if it's Swagger 2.0 or OpenAPI 3.x
	var raw map[string]interface{}

	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Check for Swagger 2.0
	if swagger, ok := raw["swagger"].(string); ok && strings.HasPrefix(swagger, "2.") {
		if m.verbose {
			fmt.Printf("  Detected Swagger 2.0, converting to OpenAPI 3.0\n")
		}
		return m.convertSwagger2ToOpenAPI3(data, ext)
	}

	// Load as OpenAPI 3.x
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	spec, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// Validate the spec
	if err := spec.Validate(context.Background()); err != nil {
		if m.verbose {
			fmt.Printf("  Warning: Validation issues: %v\n", err)
		}
	}

	return spec, nil
}

// fetchFromURL fetches data from an HTTP/HTTPS URL.
// Automatically converts GitHub blob URLs to raw URLs.
// Uses GITHUB_TOKEN environment variable for authentication with GitHub URLs.
func (m *Merger) fetchFromURL(url string) ([]byte, string, error) {
	// Convert GitHub blob URLs to raw URLs
	url = convertGitHubURL(url)

	if m.verbose {
		fmt.Printf("  Fetching from URL: %s\n", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add GitHub token authentication if available and URL is GitHub
	if isGitHubURL(url) {
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			req.Header.Set("Authorization", "token "+token)
			if m.verbose {
				fmt.Printf("  Using GITHUB_TOKEN for authentication\n")
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Determine extension from URL
	ext := strings.ToLower(filepath.Ext(url))
	// Handle URLs with query params
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}

	return data, ext, nil
}

// isGitHubURL checks if a URL is a GitHub URL that can use token auth.
func isGitHubURL(url string) bool {
	return strings.Contains(url, "github.com") ||
		strings.Contains(url, "githubusercontent.com") ||
		strings.Contains(url, "github.io")
}

// convertGitHubURL converts GitHub blob/tree URLs to raw.githubusercontent.com URLs.
// Example: https://github.com/owner/repo/blob/branch/path/file.json
//       -> https://raw.githubusercontent.com/owner/repo/branch/path/file.json
func convertGitHubURL(url string) string {
	// Match GitHub blob URLs
	githubBlobRegex := regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/blob/(.+)$`)
	if matches := githubBlobRegex.FindStringSubmatch(url); matches != nil {
		owner := matches[1]
		repo := matches[2]
		pathWithBranch := matches[3]
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", owner, repo, pathWithBranch)
	}

	// Match GitHub tree URLs (for directories, though usually not used for single files)
	githubTreeRegex := regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/tree/(.+)$`)
	if matches := githubTreeRegex.FindStringSubmatch(url); matches != nil {
		owner := matches[1]
		repo := matches[2]
		pathWithBranch := matches[3]
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", owner, repo, pathWithBranch)
	}

	return url
}

// convertSwagger2ToOpenAPI3 converts a Swagger 2.0 spec to OpenAPI 3.0.
func (m *Merger) convertSwagger2ToOpenAPI3(data []byte, ext string) (*openapi3.T, error) {
	// Parse Swagger 2.0 spec
	var swagger2Doc openapi2.T

	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &swagger2Doc); err != nil {
			return nil, fmt.Errorf("failed to parse Swagger 2.0 YAML: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &swagger2Doc); err != nil {
			return nil, fmt.Errorf("failed to parse Swagger 2.0 JSON: %w", err)
		}
	}

	// Convert to OpenAPI 3.0
	spec, err := openapi2conv.ToV3(&swagger2Doc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Swagger 2.0 to OpenAPI 3.0: %w", err)
	}

	// Ensure OpenAPI version is set to 3.0
	spec.OpenAPI = "3.0.3"

	return spec, nil
}

// filterOperations applies operation selection filters.
func (m *Merger) filterOperations(spec *openapi3.T, input *config.InputConfig) *openapi3.T {
	if input.OperationSelection == nil {
		return spec
	}

	sel := input.OperationSelection
	if spec.Paths == nil {
		return spec
	}

	pathsToRemove := make([]string, 0)

	for path, pathItem := range spec.Paths.Map() {
		if pathItem == nil {
			continue
		}

		operations := getOperationsMap(pathItem)

		for method, op := range operations {
			if op == nil {
				continue
			}

			shouldInclude := m.shouldIncludeOperation(path, method, op, sel)

			if !shouldInclude {
				// Remove the operation
				removeOperation(pathItem, method)
			}
		}

		// Check if path item is now empty
		if isPathItemEmpty(pathItem) {
			pathsToRemove = append(pathsToRemove, path)
		}
	}

	// Remove empty paths
	for _, path := range pathsToRemove {
		spec.Paths.Delete(path)
	}

	return spec
}

// shouldIncludeOperation determines if an operation should be included based on filters.
func (m *Merger) shouldIncludeOperation(path, method string, op *openapi3.Operation, sel *config.OperationSelectionConfig) bool {
	// Check includeTags
	if len(sel.IncludeTags) > 0 {
		hasMatchingTag := false
		for _, opTag := range op.Tags {
			for _, includeTag := range sel.IncludeTags {
				if opTag == includeTag {
					hasMatchingTag = true
					break
				}
			}
			if hasMatchingTag {
				break
			}
		}
		if !hasMatchingTag {
			return false
		}
	}

	// Check excludeTags
	if len(sel.ExcludeTags) > 0 {
		for _, opTag := range op.Tags {
			for _, excludeTag := range sel.ExcludeTags {
				if opTag == excludeTag {
					return false
				}
			}
		}
	}

	// Check includePaths
	if len(sel.IncludePaths) > 0 {
		matched := false
		for _, filter := range sel.IncludePaths {
			if matchPathFilter(path, method, filter) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check excludePaths
	if len(sel.ExcludePaths) > 0 {
		for _, filter := range sel.ExcludePaths {
			if matchPathFilter(path, method, filter) {
				return false
			}
		}
	}

	return true
}

// modifyPaths applies path modifications (stripStart, prepend).
func (m *Merger) modifyPaths(spec *openapi3.T, input *config.InputConfig) *openapi3.T {
	if input.PathModification == nil {
		return spec
	}

	if spec.Paths == nil {
		return spec
	}

	mod := input.PathModification
	newPaths := openapi3.NewPaths()

	for path, pathItem := range spec.Paths.Map() {
		newPath := path

		// Apply stripStart
		if mod.StripStart != "" && strings.HasPrefix(newPath, mod.StripStart) {
			newPath = strings.TrimPrefix(newPath, mod.StripStart)
		}

		// Apply prepend
		if mod.Prepend != "" {
			newPath = mod.Prepend + newPath
		}

		// Ensure path starts with /
		if !strings.HasPrefix(newPath, "/") {
			newPath = "/" + newPath
		}

		// Update refs in pathItem to reflect path change
		newPaths.Set(newPath, pathItem)
	}

	spec.Paths = newPaths
	return spec
}

// modifyParameters applies parameter modifications (include/exclude).
func (m *Merger) modifyParameters(spec *openapi3.T, input *config.InputConfig) *openapi3.T {
	if spec.Paths == nil {
		return spec
	}

	for _, pathItem := range spec.Paths.Map() {
		if pathItem == nil {
			continue
		}

		operations := getOperationsMap(pathItem)

		for _, op := range operations {
			if op == nil {
				continue
			}

			// Inject extra parameters
			if len(input.IncludeExtraParameters) > 0 {
				for _, paramCfg := range input.IncludeExtraParameters {
					param := paramCfg.ToOpenAPI3Parameter()
					// Check if parameter already exists
					exists := false
					for _, existingParam := range op.Parameters {
						if existingParam.Value != nil &&
							existingParam.Value.Name == param.Name &&
							existingParam.Value.In == param.In {
							exists = true
							break
						}
					}
					if !exists {
						op.Parameters = append(op.Parameters, &openapi3.ParameterRef{
							Value: param,
						})
					}
				}
			}

			// Remove excluded parameters
			if len(input.ExcludeParameters) > 0 {
				filteredParams := make(openapi3.Parameters, 0)
				for _, paramRef := range op.Parameters {
					if paramRef.Value == nil {
						filteredParams = append(filteredParams, paramRef)
						continue
					}
					param := paramRef.Value
					excluded := false
					for _, filter := range input.ExcludeParameters {
						if filter.Name == param.Name {
							if filter.In == "" || filter.In == param.In {
								excluded = true
								break
							}
						}
					}
					if !excluded {
						filteredParams = append(filteredParams, paramRef)
					}
				}
				op.Parameters = filteredParams
			}
		}
	}

	return spec
}

// applyDisputePrefix applies prefix to all component names and updates refs.
func (m *Merger) applyDisputePrefix(spec *openapi3.T, prefix string) *openapi3.T {
	if spec.Components == nil {
		return spec
	}

	// Build rename map
	renames := make(map[string]string)

	// Rename schemas
	if len(spec.Components.Schemas) > 0 {
		newSchemas := make(openapi3.Schemas)
		for name, schema := range spec.Components.Schemas {
			newName := prefix + name
			renames["#/components/schemas/"+name] = "#/components/schemas/" + newName
			renames["#/definitions/"+name] = "#/components/schemas/" + newName
			newSchemas[newName] = schema
		}
		spec.Components.Schemas = newSchemas
	}

	// Rename responses
	if len(spec.Components.Responses) > 0 {
		newResponses := make(openapi3.ResponseBodies)
		for name, resp := range spec.Components.Responses {
			newName := prefix + name
			renames["#/components/responses/"+name] = "#/components/responses/" + newName
			newResponses[newName] = resp
		}
		spec.Components.Responses = newResponses
	}

	// Rename parameters
	if len(spec.Components.Parameters) > 0 {
		newParams := make(openapi3.ParametersMap)
		for name, param := range spec.Components.Parameters {
			newName := prefix + name
			renames["#/components/parameters/"+name] = "#/components/parameters/" + newName
			newParams[newName] = param
		}
		spec.Components.Parameters = newParams
	}

	// Rename security schemes
	if len(spec.Components.SecuritySchemes) > 0 {
		newSchemes := make(openapi3.SecuritySchemes)
		for name, scheme := range spec.Components.SecuritySchemes {
			newName := prefix + name
			renames["#/components/securitySchemes/"+name] = "#/components/securitySchemes/" + newName
			newSchemes[newName] = scheme
		}
		spec.Components.SecuritySchemes = newSchemes
	}

	// Rename request bodies
	if len(spec.Components.RequestBodies) > 0 {
		newBodies := make(openapi3.RequestBodies)
		for name, body := range spec.Components.RequestBodies {
			newName := prefix + name
			renames["#/components/requestBodies/"+name] = "#/components/requestBodies/" + newName
			newBodies[newName] = body
		}
		spec.Components.RequestBodies = newBodies
	}

	// Update all $ref references
	updateRefs(spec, renames)

	return spec
}

// mergeSpec merges a processed spec into the master spec.
func (m *Merger) mergeSpec(spec *openapi3.T, input *config.InputConfig) error {
	// Merge paths
	if spec.Paths != nil {
		for path, pathItem := range spec.Paths.Map() {
			existingPath := m.master.Paths.Find(path)
			if existingPath != nil {
				// Merge operations into existing path
				mergePathItem(existingPath, pathItem)
			} else {
				m.master.Paths.Set(path, pathItem)
			}
		}
	}

	// Merge components
	if spec.Components != nil {
		if err := m.mergeComponents(spec.Components, input); err != nil {
			return err
		}
	}

	// Merge tags
	if len(spec.Tags) > 0 {
		for _, tag := range spec.Tags {
			if !m.hasTag(tag.Name) {
				m.master.Tags = append(m.master.Tags, tag)
			}
		}
	}

	return nil
}

// mergeComponents merges components from spec into master.
func (m *Merger) mergeComponents(components *openapi3.Components, input *config.InputConfig) error {
	hasDisputePrefix := input.Dispute != nil && input.Dispute.Prefix != ""

	// Merge schemas
	for name, schema := range components.Schemas {
		if existing, ok := m.master.Components.Schemas[name]; ok {
			if !schemasEqual(existing, schema) && !hasDisputePrefix {
				return fmt.Errorf("schema collision for '%s' without dispute prefix", name)
			}
			// Skip if exact match or has dispute prefix (already renamed)
			continue
		}
		m.master.Components.Schemas[name] = schema
	}

	// Merge responses
	for name, resp := range components.Responses {
		if _, ok := m.master.Components.Responses[name]; !ok {
			m.master.Components.Responses[name] = resp
		}
	}

	// Merge parameters
	for name, param := range components.Parameters {
		if _, ok := m.master.Components.Parameters[name]; !ok {
			m.master.Components.Parameters[name] = param
		}
	}

	// Merge security schemes
	for name, scheme := range components.SecuritySchemes {
		if _, ok := m.master.Components.SecuritySchemes[name]; !ok {
			m.master.Components.SecuritySchemes[name] = scheme
		}
	}

	// Merge request bodies
	for name, body := range components.RequestBodies {
		if _, ok := m.master.Components.RequestBodies[name]; !ok {
			m.master.Components.RequestBodies[name] = body
		}
	}

	// Merge examples
	for name, example := range components.Examples {
		if _, ok := m.master.Components.Examples[name]; !ok {
			m.master.Components.Examples[name] = example
		}
	}

	// Merge headers
	for name, header := range components.Headers {
		if _, ok := m.master.Components.Headers[name]; !ok {
			m.master.Components.Headers[name] = header
		}
	}

	// Merge links
	for name, link := range components.Links {
		if _, ok := m.master.Components.Links[name]; !ok {
			m.master.Components.Links[name] = link
		}
	}

	// Merge callbacks
	for name, callback := range components.Callbacks {
		if _, ok := m.master.Components.Callbacks[name]; !ok {
			m.master.Components.Callbacks[name] = callback
		}
	}

	return nil
}

// applyOverrides applies configuration overrides to the master spec.
func (m *Merger) applyOverrides(mergedDescriptions []string) {
	// Apply global basePath to all paths
	if m.cfg.BasePath != "" {
		m.applyBasePath()
	}

	// Apply info override
	if m.cfg.Info != nil {
		info := m.cfg.Info.ToOpenAPI3Info()
		if info != nil {
			if info.Title != "" {
				m.master.Info.Title = info.Title
			}
			if info.Version != "" {
				m.master.Info.Version = info.Version
			}
			if info.Description != "" {
				m.master.Info.Description = info.Description
			}
			if info.TermsOfService != "" {
				m.master.Info.TermsOfService = info.TermsOfService
			}
			if info.Contact != nil {
				m.master.Info.Contact = info.Contact
			}
			if info.License != nil {
				m.master.Info.License = info.License
			}
		}
	}

	// Append merged descriptions
	if len(mergedDescriptions) > 0 {
		existingDesc := m.master.Info.Description
		if existingDesc != "" {
			existingDesc += "\n\n"
		}
		m.master.Info.Description = existingDesc + strings.Join(mergedDescriptions, "\n\n")
	}

	// Apply servers override
	if len(m.cfg.Servers) > 0 {
		m.master.Servers = config.ToOpenAPI3Servers(m.cfg.Servers)
	}

	// Apply security schemes (components.securitySchemes)
	if len(m.cfg.SecuritySchemes) > 0 {
		if m.master.Components == nil {
			m.master.Components = &openapi3.Components{}
		}
		if m.master.Components.SecuritySchemes == nil {
			m.master.Components.SecuritySchemes = make(openapi3.SecuritySchemes)
		}
		// Merge security schemes from config
		for name, schemeRef := range config.ToOpenAPI3SecuritySchemes(m.cfg.SecuritySchemes) {
			m.master.Components.SecuritySchemes[name] = schemeRef
		}
	}

	// Apply security requirements (global security)
	if len(m.cfg.Security) > 0 {
		m.master.Security = config.ToOpenAPI3Security(m.cfg.Security)
	}
}

// applyBasePath prepends the global basePath to all paths.
func (m *Merger) applyBasePath() {
	if m.master.Paths == nil {
		return
	}

	basePath := m.cfg.BasePath
	// Ensure basePath starts with / and doesn't end with /
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	basePath = strings.TrimSuffix(basePath, "/")

	newPaths := openapi3.NewPaths()
	for path, pathItem := range m.master.Paths.Map() {
		newPath := basePath + path
		newPaths.Set(newPath, pathItem)
	}
	m.master.Paths = newPaths

	if m.verbose {
		fmt.Printf("Applied global basePath: %s\n", basePath)
	}
}

// sortOutput sorts tags and paths according to configuration.
func (m *Merger) sortOutput() {
	// Sort tags
	if len(m.cfg.TagOrder) > 0 {
		m.sortTags()
	}

	// Paths are sorted during output since openapi3.Paths is a map
	// We'll handle this in writeOutput
}

// sortTags sorts the tags based on tagOrder configuration.
func (m *Merger) sortTags() {
	if len(m.master.Tags) == 0 {
		return
	}

	tagOrder := m.cfg.TagOrder
	tagMap := make(map[string]*openapi3.Tag)
	for _, tag := range m.master.Tags {
		tagMap[tag.Name] = tag
	}

	sortedTags := make(openapi3.Tags, 0, len(m.master.Tags))

	// Add tags in specified order
	for _, tagName := range tagOrder {
		if tag, ok := tagMap[tagName]; ok {
			sortedTags = append(sortedTags, tag)
			delete(tagMap, tagName)
		}
	}

	// Add remaining tags
	for _, tag := range m.master.Tags {
		if _, ok := tagMap[tag.Name]; ok {
			sortedTags = append(sortedTags, tag)
		}
	}

	m.master.Tags = sortedTags
}

// writeOutput serializes and writes the master spec to disk.
func (m *Merger) writeOutput() error {
	// Create output directory if needed
	outputDir := filepath.Dir(m.cfg.Output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Determine output format
	ext := strings.ToLower(filepath.Ext(m.cfg.Output))
	var data []byte
	var err error

	if ext == ".yaml" || ext == ".yml" {
		data, err = m.marshalYAML()
	} else {
		data, err = m.marshalJSON()
	}

	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if err := os.WriteFile(m.cfg.Output, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// marshalJSON marshals the spec to JSON with sorted paths.
func (m *Merger) marshalJSON() ([]byte, error) {
	// Sort paths for deterministic output
	sortedSpec := m.createSortedSpec()
	return json.MarshalIndent(sortedSpec, "", "  ")
}

// marshalYAML marshals the spec to YAML with sorted paths.
func (m *Merger) marshalYAML() ([]byte, error) {
	sortedSpec := m.createSortedSpec()
	return yaml.Marshal(sortedSpec)
}

// createSortedSpec creates a copy of the spec with sorted paths.
func (m *Merger) createSortedSpec() map[string]interface{} {
	// Convert to map for custom ordering
	data, _ := json.Marshal(m.master)
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)

	// Sort paths
	if paths, ok := result["paths"].(map[string]interface{}); ok {
		sortedPaths := m.sortPaths(paths)
		result["paths"] = sortedPaths
	}

	return result
}

// sortPaths sorts paths according to pathsOrder configuration.
func (m *Merger) sortPaths(paths map[string]interface{}) map[string]interface{} {
	// Create ordered map
	orderedPaths := make(map[string]interface{})

	// Get all path keys
	allPaths := make([]string, 0, len(paths))
	for path := range paths {
		allPaths = append(allPaths, path)
	}

	// Sort: priority paths first, then alphabetically
	sortedPaths := make([]string, 0, len(allPaths))

	// Add priority paths first
	for _, priorityPath := range m.cfg.PathsOrder {
		for _, path := range allPaths {
			if path == priorityPath {
				sortedPaths = append(sortedPaths, path)
			}
		}
	}

	// Add remaining paths alphabetically
	remainingPaths := make([]string, 0)
	for _, path := range allPaths {
		isPriority := false
		for _, priorityPath := range m.cfg.PathsOrder {
			if path == priorityPath {
				isPriority = true
				break
			}
		}
		if !isPriority {
			remainingPaths = append(remainingPaths, path)
		}
	}

	// Sort remaining paths
	for i := 0; i < len(remainingPaths); i++ {
		for j := i + 1; j < len(remainingPaths); j++ {
			if remainingPaths[i] > remainingPaths[j] {
				remainingPaths[i], remainingPaths[j] = remainingPaths[j], remainingPaths[i]
			}
		}
	}

	sortedPaths = append(sortedPaths, remainingPaths...)

	// Build ordered map
	for _, path := range sortedPaths {
		orderedPaths[path] = paths[path]
	}

	return orderedPaths
}

// formatDescription formats a description with optional title.
func (m *Merger) formatDescription(desc string, cfg *config.DescriptionConfig) string {
	if desc == "" {
		return ""
	}

	if cfg.Title != nil && cfg.Title.Value != "" {
		level := cfg.Title.HeadingLevel
		if level < 1 || level > 6 {
			level = 2
		}
		heading := strings.Repeat("#", level)
		return fmt.Sprintf("%s %s\n\n%s", heading, cfg.Title.Value, desc)
	}

	return desc
}

// hasTag checks if a tag with the given name already exists.
func (m *Merger) hasTag(name string) bool {
	for _, tag := range m.master.Tags {
		if tag.Name == name {
			return true
		}
	}
	return false
}

// schemasEqual compares two schema refs for equality (simple comparison).
func schemasEqual(a, b *openapi3.SchemaRef) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Simple reference comparison
	if a.Ref != "" && b.Ref != "" {
		return a.Ref == b.Ref
	}
	// For value comparison, we do a simple JSON comparison
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
