// Package formatter provides output formatters for USFM documents.
// It supports multiple output formats including JSON, plain text, and TSV.
package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/arenzana/usfmp/pkg/usfm"
)

// FormatJSON formats one or more USFM documents as JSON output.
// If a single document is provided, it returns the document directly.
// If multiple documents are provided, they are returned as a JSON array.
// The output is formatted with proper indentation for readability.
//
// Example output for single document:
//
//	{
//	  "id": "GEN",
//	  "main_title": "Genesis",
//	  "chapters": [...]
//	}
//
// Example output for multiple documents:
//
//	[
//	  {"id": "GEN", ...},
//	  {"id": "EXO", ...}
//	]
func FormatJSON(documents []*usfm.Document) (string, error) {
	if len(documents) == 0 {
		return "[]", nil
	}
	
	// If single document, return it directly (not wrapped in array)
	if len(documents) == 1 {
		data, err := json.MarshalIndent(documents[0], "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal document to JSON: %w", err)
		}
		return string(data), nil
	}
	
	// Multiple documents - return as array
	data, err := json.MarshalIndent(documents, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal documents to JSON: %w", err)
	}
	
	return string(data), nil
}