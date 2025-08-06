package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/arenzana/usfmp/pkg/usfm"
)

// createTestDocument creates a test document for formatting tests
func createTestDocument() *usfm.Document {
	return &usfm.Document{
		ID:        "GEN",
		Header:    "Genesis",
		MainTitle: "Genesis",
		TableOfContents: []usfm.TOCEntry{
			{Level: 1, Text: "Genesis"},
			{Level: 2, Text: "Gen"},
		},
		Chapters: []usfm.Chapter{
			{
				Number: 1,
				Sections: []usfm.Section{
					{
						Level:     1,
						Title:     "The Creation",
						Reference: "(John 1:1–5)",
						Verses: []usfm.Verse{
							{
								Number: 1,
								Text:   "In the beginning God created the heavens and the earth.",
								Footnotes: []usfm.Footnote{
									{
										Caller:    "+",
										Reference: "1:1",
										Text:      "Hebrew: Elohim",
									},
								},
							},
							{
								Number: 2,
								Text:   "Now the earth was formless and void.",
							},
						},
					},
				},
			},
		},
		ParsedAt:   time.Now(),
		SourceFile: "test.sfm",
	}
}

// TestFormatJSON tests JSON formatting
func TestFormatJSON(t *testing.T) {
	doc := createTestDocument()
	documents := []*usfm.Document{doc}

	result, err := FormatJSON(documents)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	// Check that result contains expected JSON structure
	if !strings.Contains(result, `"id": "GEN"`) {
		t.Error("JSON output should contain document ID")
	}

	if !strings.Contains(result, `"main_title": "Genesis"`) {
		t.Error("JSON output should contain main title")
	}

	if !strings.Contains(result, `"number": 1`) {
		t.Error("JSON output should contain chapter number")
	}

	if !strings.Contains(result, `"The Creation"`) {
		t.Error("JSON output should contain section title")
	}
}

// TestFormatJSONMultipleDocuments tests JSON formatting with multiple documents
func TestFormatJSONMultipleDocuments(t *testing.T) {
	doc1 := createTestDocument()
	doc2 := createTestDocument()
	doc2.ID = "EXO"
	doc2.MainTitle = "Exodus"

	documents := []*usfm.Document{doc1, doc2}

	result, err := FormatJSON(documents)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	// Should be an array when multiple documents
	if !strings.HasPrefix(result, "[") {
		t.Error("Multiple documents should be formatted as JSON array")
	}

	if !strings.Contains(result, `"id": "GEN"`) {
		t.Error("Should contain first document")
	}

	if !strings.Contains(result, `"id": "EXO"`) {
		t.Error("Should contain second document")
	}
}

// TestFormatJSONEmptyDocuments tests JSON formatting with empty input
func TestFormatJSONEmptyDocuments(t *testing.T) {
	result, err := FormatJSON([]*usfm.Document{})
	if err != nil {
		t.Fatalf("FormatJSON failed on empty input: %v", err)
	}

	if result != "[]" {
		t.Errorf("Expected empty array '[]', got '%s'", result)
	}
}

// TestFormatText tests text formatting
func TestFormatText(t *testing.T) {
	doc := createTestDocument()
	documents := []*usfm.Document{doc}

	result, err := FormatText(documents)
	if err != nil {
		t.Fatalf("FormatText failed: %v", err)
	}

	// Check for expected text content
	if !strings.Contains(result, "Genesis") {
		t.Error("Text output should contain main title")
	}

	if !strings.Contains(result, "Book: GEN") {
		t.Error("Text output should contain book ID")
	}

	if !strings.Contains(result, "Chapter 1") {
		t.Error("Text output should contain chapter number")
	}

	if !strings.Contains(result, "The Creation") {
		t.Error("Text output should contain section title")
	}

	if !strings.Contains(result, "1. In the beginning") {
		t.Error("Text output should contain verse text with number")
	}

	if !strings.Contains(result, "[+:1:1 - Hebrew: Elohim]") {
		t.Error("Text output should contain footnotes")
	}
}

// TestFormatTSV tests TSV formatting
func TestFormatTSV(t *testing.T) {
	doc := createTestDocument()
	documents := []*usfm.Document{doc}

	result, err := FormatTSV(documents)
	if err != nil {
		t.Fatalf("FormatTSV failed: %v", err)
	}

	lines := strings.Split(result, "\n")

	// Should have header row + data rows + empty line at end
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines (header + data), got %d", len(lines))
	}

	// Check header
	expectedHeader := "Book\tChapter\tVerse\tSection_Title\tSection_Level\tVerse_Text\tFootnotes\tReferences"
	if lines[0] != expectedHeader {
		t.Errorf("Expected header '%s', got '%s'", expectedHeader, lines[0])
	}

	// Check first data row
	if !strings.HasPrefix(lines[1], "GEN\t1\t1\t") {
		t.Error("First data row should start with 'GEN\\t1\\t1\\t'")
	}

	if !strings.Contains(lines[1], "The Creation") {
		t.Error("TSV should contain section title")
	}

	if !strings.Contains(lines[1], "In the beginning God created") {
		t.Error("TSV should contain verse text")
	}

	if !strings.Contains(lines[1], "+:1:1=Hebrew: Elohim") {
		t.Error("TSV should contain formatted footnotes")
	}
}

// TestCleanTSVField tests the TSV field cleaning function
func TestCleanTSVField(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"text\twith\ttabs", "text with tabs"},
		{"text\nwith\nnewlines", "text with newlines"},
		{"text\r\nwith\r\nCRLF", "text with CRLF"},
		{"text  with   multiple    spaces", "text with multiple spaces"},
		{"\t\n  whitespace  \r\n\t", "whitespace"},
		{"", ""},
	}

	for i, tc := range testCases {
		result := cleanTSVField(tc.input)
		if result != tc.expected {
			t.Errorf("Test case %d: expected '%s', got '%s'", i+1, tc.expected, result)
		}
	}
}
