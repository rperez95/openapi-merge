// Package config provides configuration types and loading for openapi-merge.
package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mitchellh/mapstructure"
)

// Config represents the main configuration for the merge operation.
type Config struct {
	// Inputs is the list of OpenAPI files to merge
	Inputs []InputConfig `mapstructure:"inputs" json:"inputs" yaml:"inputs"`

	// Output is the path to save the merged file
	Output string `mapstructure:"output" json:"output" yaml:"output"`

	// BasePath is a global prefix prepended to all paths after individual processing
	BasePath string `mapstructure:"basePath" json:"basePath,omitempty" yaml:"basePath,omitempty"`

	// Info contains metadata to override in the final file
	Info *InfoConfig `mapstructure:"info" json:"info,omitempty" yaml:"info,omitempty"`

	// Servers is the list of servers to replace in the final file
	Servers []ServerConfig `mapstructure:"servers" json:"servers,omitempty" yaml:"servers,omitempty"`

	// SecuritySchemes defines authentication methods (OAS3 components.securitySchemes)
	SecuritySchemes map[string]SecuritySchemeConfig `mapstructure:"securitySchemes" json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`

	// Security contains global security requirements
	Security []map[string][]string `mapstructure:"security" json:"security,omitempty" yaml:"security,omitempty"`

	// TagOrder defines the order of tags in the output
	TagOrder []string `mapstructure:"tagOrder" json:"tagOrder,omitempty" yaml:"tagOrder,omitempty"`

	// PathsOrder defines high-priority paths that should appear first
	PathsOrder []string `mapstructure:"pathsOrder" json:"pathsOrder,omitempty" yaml:"pathsOrder,omitempty"`
}

// InfoConfig represents the info section override configuration.
type InfoConfig struct {
	Title          string         `mapstructure:"title" json:"title,omitempty" yaml:"title,omitempty"`
	Description    string         `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`
	Version        string         `mapstructure:"version" json:"version,omitempty" yaml:"version,omitempty"`
	TermsOfService string         `mapstructure:"termsOfService" json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *ContactConfig `mapstructure:"contact" json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *LicenseConfig `mapstructure:"license" json:"license,omitempty" yaml:"license,omitempty"`
}

// ContactConfig represents contact information.
type ContactConfig struct {
	Name  string `mapstructure:"name" json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `mapstructure:"url" json:"url,omitempty" yaml:"url,omitempty"`
	Email string `mapstructure:"email" json:"email,omitempty" yaml:"email,omitempty"`
}

// LicenseConfig represents license information.
type LicenseConfig struct {
	Name string `mapstructure:"name" json:"name,omitempty" yaml:"name,omitempty"`
	URL  string `mapstructure:"url" json:"url,omitempty" yaml:"url,omitempty"`
}

// ServerConfig represents a server configuration.
type ServerConfig struct {
	URL         string                          `mapstructure:"url" json:"url" yaml:"url"`
	Description string                          `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariableConfig `mapstructure:"variables" json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariableConfig represents a server variable.
type ServerVariableConfig struct {
	Enum        []string `mapstructure:"enum" json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `mapstructure:"default" json:"default" yaml:"default"`
	Description string   `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`
}

// SecuritySchemeConfig represents an OAS3 security scheme definition.
// Supports: apiKey, http (basic/bearer), oauth2, openIdConnect
type SecuritySchemeConfig struct {
	// Type is the security scheme type: apiKey, http, oauth2, openIdConnect
	Type string `mapstructure:"type" json:"type" yaml:"type"`

	// Description of the security scheme
	Description string `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`

	// Name is the name of the header, query or cookie parameter (for apiKey type)
	Name string `mapstructure:"name" json:"name,omitempty" yaml:"name,omitempty"`

	// In is the location of the API key: header, query, or cookie (for apiKey type)
	In string `mapstructure:"in" json:"in,omitempty" yaml:"in,omitempty"`

	// Scheme is the HTTP auth scheme: basic, bearer, etc. (for http type)
	Scheme string `mapstructure:"scheme" json:"scheme,omitempty" yaml:"scheme,omitempty"`

	// BearerFormat is the format of the bearer token (for http bearer type)
	BearerFormat string `mapstructure:"bearerFormat" json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`

	// Flows contains OAuth2 flow configurations (for oauth2 type)
	Flows *OAuthFlowsConfig `mapstructure:"flows" json:"flows,omitempty" yaml:"flows,omitempty"`

	// OpenIdConnectUrl is the URL for OpenID Connect discovery (for openIdConnect type)
	OpenIdConnectUrl string `mapstructure:"openIdConnectUrl" json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlowsConfig represents OAuth2 flow configurations.
type OAuthFlowsConfig struct {
	Implicit          *OAuthFlowConfig `mapstructure:"implicit" json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlowConfig `mapstructure:"password" json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlowConfig `mapstructure:"clientCredentials" json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlowConfig `mapstructure:"authorizationCode" json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlowConfig represents a single OAuth2 flow configuration.
type OAuthFlowConfig struct {
	AuthorizationURL string            `mapstructure:"authorizationUrl" json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `mapstructure:"tokenUrl" json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `mapstructure:"refreshUrl" json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `mapstructure:"scopes" json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

// InputConfig represents a single input file configuration.
type InputConfig struct {
	// InputFile is the path to the source file (JSON or YAML)
	InputFile string `mapstructure:"inputFile" json:"inputFile" yaml:"inputFile"`

	// Dispute defines conflict resolution with prefix
	Dispute *DisputeConfig `mapstructure:"dispute" json:"dispute,omitempty" yaml:"dispute,omitempty"`

	// PathModification defines path transformation rules
	PathModification *PathModificationConfig `mapstructure:"pathModification" json:"pathModification,omitempty" yaml:"pathModification,omitempty"`

	// OperationSelection defines which operations to include/exclude
	OperationSelection *OperationSelectionConfig `mapstructure:"operationSelection" json:"operationSelection,omitempty" yaml:"operationSelection,omitempty"`

	// IncludeExtraParameters are parameters to inject into every operation
	IncludeExtraParameters []ParameterConfig `mapstructure:"includeExtraParameters" json:"includeExtraParameters,omitempty" yaml:"includeExtraParameters,omitempty"`

	// ExcludeParameters are parameter filters to remove from operations
	ExcludeParameters []ParamFilter `mapstructure:"excludeParameters" json:"excludeParameters,omitempty" yaml:"excludeParameters,omitempty"`

	// Description defines how to merge the input's description
	Description *DescriptionConfig `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`
}

// DisputeConfig defines conflict resolution configuration.
type DisputeConfig struct {
	// Prefix to add to component names on collision
	Prefix string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
}

// PathModificationConfig defines path transformation rules.
type PathModificationConfig struct {
	// StripStart is a string to remove from the beginning of paths
	StripStart string `mapstructure:"stripStart" json:"stripStart,omitempty" yaml:"stripStart,omitempty"`

	// Prepend is a string to add to the start of paths
	Prepend string `mapstructure:"prepend" json:"prepend,omitempty" yaml:"prepend,omitempty"`
}

// OperationSelectionConfig defines operation filtering rules.
type OperationSelectionConfig struct {
	// IncludeTags - only include operations with these tags
	IncludeTags []string `mapstructure:"includeTags" json:"includeTags,omitempty" yaml:"includeTags,omitempty"`

	// ExcludeTags - exclude operations with these tags
	ExcludeTags []string `mapstructure:"excludeTags" json:"excludeTags,omitempty" yaml:"excludeTags,omitempty"`

	// IncludePaths - whitelist specific paths/methods
	IncludePaths []PathFilter `mapstructure:"includePaths" json:"includePaths,omitempty" yaml:"includePaths,omitempty"`

	// ExcludePaths - blacklist specific paths/methods
	ExcludePaths []PathFilter `mapstructure:"excludePaths" json:"excludePaths,omitempty" yaml:"excludePaths,omitempty"`
}

// PathFilter represents a path/method filter with glob support.
type PathFilter struct {
	// Path supports glob matching (e.g., /api/*)
	Path string `mapstructure:"path" json:"path" yaml:"path"`

	// Method is the HTTP verb (GET, POST, etc.) or empty for all methods
	Method string `mapstructure:"method" json:"method,omitempty" yaml:"method,omitempty"`
}

// ParamFilter represents a parameter filter.
type ParamFilter struct {
	// Name is the parameter name to match
	Name string `mapstructure:"name" json:"name" yaml:"name"`

	// In is the parameter location (query, header, path, cookie)
	In string `mapstructure:"in" json:"in,omitempty" yaml:"in,omitempty"`
}

// ParameterConfig represents a parameter to inject.
type ParameterConfig struct {
	Name            string      `mapstructure:"name" json:"name" yaml:"name"`
	In              string      `mapstructure:"in" json:"in" yaml:"in"`
	Description     string      `mapstructure:"description" json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool        `mapstructure:"required" json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool        `mapstructure:"deprecated" json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool        `mapstructure:"allowEmptyValue" json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Schema          interface{} `mapstructure:"schema" json:"schema,omitempty" yaml:"schema,omitempty"`
}

// DescriptionConfig defines description merging logic.
type DescriptionConfig struct {
	// Append indicates whether to append the input's description
	Append bool `mapstructure:"append" json:"append,omitempty" yaml:"append,omitempty"`

	// Title configuration for the description section
	Title *DescriptionTitleConfig `mapstructure:"title" json:"title,omitempty" yaml:"title,omitempty"`
}

// DescriptionTitleConfig defines the title for description sections.
type DescriptionTitleConfig struct {
	// Value is the title text
	Value string `mapstructure:"value" json:"value" yaml:"value"`

	// HeadingLevel is the markdown heading level (1-6)
	HeadingLevel int `mapstructure:"headingLevel" json:"headingLevel,omitempty" yaml:"headingLevel,omitempty"`
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if len(c.Inputs) == 0 {
		return fmt.Errorf("at least one input file is required")
	}

	if c.Output == "" {
		return fmt.Errorf("output file path is required")
	}

	for i, input := range c.Inputs {
		if input.InputFile == "" {
			return fmt.Errorf("input[%d]: inputFile is required", i)
		}
	}

	return nil
}

// IsURL checks if a path is an HTTP/HTTPS URL.
func IsURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// ResolveRelativePaths resolves relative paths based on the config directory.
// URLs (http:// or https://) are left unchanged.
func (c *Config) ResolveRelativePaths(configDir string) {
	for i := range c.Inputs {
		// Skip URLs - they don't need path resolution
		if IsURL(c.Inputs[i].InputFile) {
			continue
		}
		if !filepath.IsAbs(c.Inputs[i].InputFile) {
			c.Inputs[i].InputFile = filepath.Join(configDir, c.Inputs[i].InputFile)
		}
	}

	if !filepath.IsAbs(c.Output) {
		c.Output = filepath.Join(configDir, c.Output)
	}
}

// ToOpenAPI3Info converts InfoConfig to openapi3.Info.
func (c *InfoConfig) ToOpenAPI3Info() *openapi3.Info {
	if c == nil {
		return nil
	}

	info := &openapi3.Info{
		Title:          c.Title,
		Description:    c.Description,
		Version:        c.Version,
		TermsOfService: c.TermsOfService,
	}

	if c.Contact != nil {
		info.Contact = &openapi3.Contact{
			Name:  c.Contact.Name,
			URL:   c.Contact.URL,
			Email: c.Contact.Email,
		}
	}

	if c.License != nil {
		info.License = &openapi3.License{
			Name: c.License.Name,
			URL:  c.License.URL,
		}
	}

	return info
}

// ToOpenAPI3Servers converts ServerConfig slice to openapi3.Servers.
func ToOpenAPI3Servers(servers []ServerConfig) openapi3.Servers {
	if len(servers) == 0 {
		return nil
	}

	result := make(openapi3.Servers, len(servers))
	for i, s := range servers {
		server := &openapi3.Server{
			URL:         s.URL,
			Description: s.Description,
		}

		if len(s.Variables) > 0 {
			server.Variables = make(map[string]*openapi3.ServerVariable)
			for name, v := range s.Variables {
				server.Variables[name] = &openapi3.ServerVariable{
					Enum:        v.Enum,
					Default:     v.Default,
					Description: v.Description,
				}
			}
		}

		result[i] = server
	}

	return result
}

// ToOpenAPI3Security converts security config to openapi3.SecurityRequirements.
func ToOpenAPI3Security(security []map[string][]string) openapi3.SecurityRequirements {
	if len(security) == 0 {
		return nil
	}

	result := make(openapi3.SecurityRequirements, len(security))
	for i, s := range security {
		result[i] = openapi3.SecurityRequirement(s)
	}

	return result
}

// ToOpenAPI3SecuritySchemes converts SecuritySchemeConfig map to openapi3.SecuritySchemes.
func ToOpenAPI3SecuritySchemes(schemes map[string]SecuritySchemeConfig) openapi3.SecuritySchemes {
	if len(schemes) == 0 {
		return nil
	}

	result := make(openapi3.SecuritySchemes)
	for name, cfg := range schemes {
		scheme := &openapi3.SecurityScheme{
			Type:             cfg.Type,
			Description:      cfg.Description,
			Name:             cfg.Name,
			In:               cfg.In,
			Scheme:           cfg.Scheme,
			BearerFormat:     cfg.BearerFormat,
			OpenIdConnectUrl: cfg.OpenIdConnectUrl,
		}

		// Convert OAuth2 flows if present
		if cfg.Flows != nil {
			scheme.Flows = &openapi3.OAuthFlows{}
			if cfg.Flows.Implicit != nil {
				scheme.Flows.Implicit = convertOAuthFlow(cfg.Flows.Implicit)
			}
			if cfg.Flows.Password != nil {
				scheme.Flows.Password = convertOAuthFlow(cfg.Flows.Password)
			}
			if cfg.Flows.ClientCredentials != nil {
				scheme.Flows.ClientCredentials = convertOAuthFlow(cfg.Flows.ClientCredentials)
			}
			if cfg.Flows.AuthorizationCode != nil {
				scheme.Flows.AuthorizationCode = convertOAuthFlow(cfg.Flows.AuthorizationCode)
			}
		}

		result[name] = &openapi3.SecuritySchemeRef{
			Value: scheme,
		}
	}

	return result
}

func convertOAuthFlow(cfg *OAuthFlowConfig) *openapi3.OAuthFlow {
	if cfg == nil {
		return nil
	}
	return &openapi3.OAuthFlow{
		AuthorizationURL: cfg.AuthorizationURL,
		TokenURL:         cfg.TokenURL,
		RefreshURL:       cfg.RefreshURL,
		Scopes:           cfg.Scopes,
	}
}

// ToOpenAPI3Parameter converts ParameterConfig to openapi3.Parameter.
func (p *ParameterConfig) ToOpenAPI3Parameter() *openapi3.Parameter {
	param := &openapi3.Parameter{
		Name:            p.Name,
		In:              p.In,
		Description:     p.Description,
		Required:        p.Required,
		Deprecated:      p.Deprecated,
		AllowEmptyValue: p.AllowEmptyValue,
	}

	// Handle schema conversion
	if p.Schema != nil {
		param.Schema = convertToSchemaRef(p.Schema)
	} else {
		// Default to string schema
		param.Schema = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"string"},
			},
		}
	}

	return param
}

func convertToSchemaRef(schema interface{}) *openapi3.SchemaRef {
	switch s := schema.(type) {
	case map[string]interface{}:
		schemaVal := &openapi3.Schema{}
		if typeVal, ok := s["type"].(string); ok {
			schemaVal.Type = &openapi3.Types{typeVal}
		}
		if format, ok := s["format"].(string); ok {
			schemaVal.Format = format
		}
		if desc, ok := s["description"].(string); ok {
			schemaVal.Description = desc
		}
		return &openapi3.SchemaRef{Value: schemaVal}
	default:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"string"},
			},
		}
	}
}

// DecodeHook returns a mapstructure decode hook for custom types.
func DecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		interfaceToMapHookFunc(),
	)
}

func interfaceToMapHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		return data, nil
	}
}
