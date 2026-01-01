package merger

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gobwas/glob"
	"github.com/rperez95/openapi-merge/internal/config"
)

// getOperationsMap returns a map of HTTP method to operation.
func getOperationsMap(pathItem *openapi3.PathItem) map[string]*openapi3.Operation {
	return map[string]*openapi3.Operation{
		"GET":     pathItem.Get,
		"POST":    pathItem.Post,
		"PUT":     pathItem.Put,
		"DELETE":  pathItem.Delete,
		"PATCH":   pathItem.Patch,
		"HEAD":    pathItem.Head,
		"OPTIONS": pathItem.Options,
		"TRACE":   pathItem.Trace,
	}
}

// removeOperation removes an operation from a path item.
func removeOperation(pathItem *openapi3.PathItem, method string) {
	switch strings.ToUpper(method) {
	case "GET":
		pathItem.Get = nil
	case "POST":
		pathItem.Post = nil
	case "PUT":
		pathItem.Put = nil
	case "DELETE":
		pathItem.Delete = nil
	case "PATCH":
		pathItem.Patch = nil
	case "HEAD":
		pathItem.Head = nil
	case "OPTIONS":
		pathItem.Options = nil
	case "TRACE":
		pathItem.Trace = nil
	}
}

// isPathItemEmpty checks if a path item has no operations.
func isPathItemEmpty(pathItem *openapi3.PathItem) bool {
	return pathItem.Get == nil &&
		pathItem.Post == nil &&
		pathItem.Put == nil &&
		pathItem.Delete == nil &&
		pathItem.Patch == nil &&
		pathItem.Head == nil &&
		pathItem.Options == nil &&
		pathItem.Trace == nil
}

// matchPathFilter checks if a path/method matches a filter.
func matchPathFilter(path, method string, filter config.PathFilter) bool {
	// Check method first (if specified)
	if filter.Method != "" && !strings.EqualFold(method, filter.Method) {
		return false
	}

	// Check path with glob matching
	return matchGlob(filter.Path, path)
}

// matchGlob performs glob matching on a path.
func matchGlob(pattern, path string) bool {
	// Handle exact match
	if pattern == path {
		return true
	}

	// Use glob library for pattern matching
	g, err := glob.Compile(pattern)
	if err != nil {
		// Fallback to exact match if pattern is invalid
		return pattern == path
	}

	return g.Match(path)
}

// mergePathItem merges operations from source into destination.
func mergePathItem(dest, src *openapi3.PathItem) {
	if src.Get != nil && dest.Get == nil {
		dest.Get = src.Get
	}
	if src.Post != nil && dest.Post == nil {
		dest.Post = src.Post
	}
	if src.Put != nil && dest.Put == nil {
		dest.Put = src.Put
	}
	if src.Delete != nil && dest.Delete == nil {
		dest.Delete = src.Delete
	}
	if src.Patch != nil && dest.Patch == nil {
		dest.Patch = src.Patch
	}
	if src.Head != nil && dest.Head == nil {
		dest.Head = src.Head
	}
	if src.Options != nil && dest.Options == nil {
		dest.Options = src.Options
	}
	if src.Trace != nil && dest.Trace == nil {
		dest.Trace = src.Trace
	}

	// Merge parameters
	if len(src.Parameters) > 0 {
		for _, param := range src.Parameters {
			exists := false
			for _, existingParam := range dest.Parameters {
				if existingParam.Value != nil && param.Value != nil &&
					existingParam.Value.Name == param.Value.Name &&
					existingParam.Value.In == param.Value.In {
					exists = true
					break
				}
			}
			if !exists {
				dest.Parameters = append(dest.Parameters, param)
			}
		}
	}
}
