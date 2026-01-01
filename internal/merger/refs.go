package merger

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// updateRefs updates all $ref references in the spec according to the rename map.
func updateRefs(spec *openapi3.T, renames map[string]string) {
	if len(renames) == 0 {
		return
	}

	// Update refs in paths
	if spec.Paths != nil {
		for _, pathItem := range spec.Paths.Map() {
			updatePathItemRefs(pathItem, renames)
		}
	}

	// Update refs in components
	if spec.Components != nil {
		updateComponentsRefs(spec.Components, renames)
	}
}

// updatePathItemRefs updates refs in a path item.
func updatePathItemRefs(pathItem *openapi3.PathItem, renames map[string]string) {
	if pathItem == nil {
		return
	}

	// Update refs in operations
	operations := []*openapi3.Operation{
		pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Delete,
		pathItem.Patch, pathItem.Head, pathItem.Options, pathItem.Trace,
	}

	for _, op := range operations {
		if op != nil {
			updateOperationRefs(op, renames)
		}
	}

	// Update refs in parameters
	for _, param := range pathItem.Parameters {
		updateParameterRefRefs(param, renames)
	}
}

// updateOperationRefs updates refs in an operation.
func updateOperationRefs(op *openapi3.Operation, renames map[string]string) {
	// Update parameters
	for _, param := range op.Parameters {
		updateParameterRefRefs(param, renames)
	}

	// Update request body
	if op.RequestBody != nil {
		updateRequestBodyRefRefs(op.RequestBody, renames)
	}

	// Update responses
	for _, resp := range op.Responses.Map() {
		updateResponseRefRefs(resp, renames)
	}

	// Update callbacks
	for _, callback := range op.Callbacks {
		updateCallbackRefRefs(callback, renames)
	}
}

// updateParameterRefRefs updates refs in a parameter ref.
func updateParameterRefRefs(paramRef *openapi3.ParameterRef, renames map[string]string) {
	if paramRef == nil {
		return
	}

	// Update the ref itself
	if paramRef.Ref != "" {
		if newRef, ok := renames[paramRef.Ref]; ok {
			paramRef.Ref = newRef
		}
	}

	// Update schema refs
	if paramRef.Value != nil && paramRef.Value.Schema != nil {
		updateSchemaRefRefs(paramRef.Value.Schema, renames)
	}
}

// updateSchemaRefRefs updates refs in a schema ref.
func updateSchemaRefRefs(schemaRef *openapi3.SchemaRef, renames map[string]string) {
	if schemaRef == nil {
		return
	}

	// Update the ref itself
	if schemaRef.Ref != "" {
		if newRef, ok := renames[schemaRef.Ref]; ok {
			schemaRef.Ref = newRef
		}
	}

	// Update nested schemas
	if schemaRef.Value != nil {
		schema := schemaRef.Value

		// Update items
		if schema.Items != nil {
			updateSchemaRefRefs(schema.Items, renames)
		}

		// Update properties
		for _, prop := range schema.Properties {
			updateSchemaRefRefs(prop, renames)
		}

		// Update additionalProperties
		if schema.AdditionalProperties.Schema != nil {
			updateSchemaRefRefs(schema.AdditionalProperties.Schema, renames)
		}

		// Update allOf
		for _, s := range schema.AllOf {
			updateSchemaRefRefs(s, renames)
		}

		// Update oneOf
		for _, s := range schema.OneOf {
			updateSchemaRefRefs(s, renames)
		}

		// Update anyOf
		for _, s := range schema.AnyOf {
			updateSchemaRefRefs(s, renames)
		}

		// Update not
		if schema.Not != nil {
			updateSchemaRefRefs(schema.Not, renames)
		}
	}
}

// updateRequestBodyRefRefs updates refs in a request body ref.
func updateRequestBodyRefRefs(bodyRef *openapi3.RequestBodyRef, renames map[string]string) {
	if bodyRef == nil {
		return
	}

	// Update the ref itself
	if bodyRef.Ref != "" {
		if newRef, ok := renames[bodyRef.Ref]; ok {
			bodyRef.Ref = newRef
		}
	}

	// Update content schemas
	if bodyRef.Value != nil && bodyRef.Value.Content != nil {
		for _, mediaType := range bodyRef.Value.Content {
			if mediaType.Schema != nil {
				updateSchemaRefRefs(mediaType.Schema, renames)
			}
		}
	}
}

// updateResponseRefRefs updates refs in a response ref.
func updateResponseRefRefs(respRef *openapi3.ResponseRef, renames map[string]string) {
	if respRef == nil {
		return
	}

	// Update the ref itself
	if respRef.Ref != "" {
		if newRef, ok := renames[respRef.Ref]; ok {
			respRef.Ref = newRef
		}
	}

	// Update content schemas
	if respRef.Value != nil {
		if respRef.Value.Content != nil {
			for _, mediaType := range respRef.Value.Content {
				if mediaType.Schema != nil {
					updateSchemaRefRefs(mediaType.Schema, renames)
				}
			}
		}

		// Update headers
		for _, header := range respRef.Value.Headers {
			updateHeaderRefRefs(header, renames)
		}
	}
}

// updateHeaderRefRefs updates refs in a header ref.
func updateHeaderRefRefs(headerRef *openapi3.HeaderRef, renames map[string]string) {
	if headerRef == nil {
		return
	}

	// Update the ref itself
	if headerRef.Ref != "" {
		if newRef, ok := renames[headerRef.Ref]; ok {
			headerRef.Ref = newRef
		}
	}

	// Update schema
	if headerRef.Value != nil && headerRef.Value.Schema != nil {
		updateSchemaRefRefs(headerRef.Value.Schema, renames)
	}
}

// updateCallbackRefRefs updates refs in a callback ref.
func updateCallbackRefRefs(callbackRef *openapi3.CallbackRef, renames map[string]string) {
	if callbackRef == nil {
		return
	}

	// Update the ref itself
	if callbackRef.Ref != "" {
		if newRef, ok := renames[callbackRef.Ref]; ok {
			callbackRef.Ref = newRef
		}
	}

	// Update path items in callback
	if callbackRef.Value != nil {
		for _, pathItem := range callbackRef.Value.Map() {
			updatePathItemRefs(pathItem, renames)
		}
	}
}

// updateComponentsRefs updates refs in components.
func updateComponentsRefs(components *openapi3.Components, renames map[string]string) {
	// Update schemas
	for _, schema := range components.Schemas {
		updateSchemaRefRefs(schema, renames)
	}

	// Update parameters
	for _, param := range components.Parameters {
		updateParameterRefRefs(param, renames)
	}

	// Update responses
	for _, resp := range components.Responses {
		updateResponseRefRefs(resp, renames)
	}

	// Update request bodies
	for _, body := range components.RequestBodies {
		updateRequestBodyRefRefs(body, renames)
	}

	// Update headers
	for _, header := range components.Headers {
		updateHeaderRefRefs(header, renames)
	}

	// Update callbacks
	for _, callback := range components.Callbacks {
		updateCallbackRefRefs(callback, renames)
	}
}
