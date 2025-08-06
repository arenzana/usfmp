// Package usfm provides a comprehensive parser for USFM (Unified Standard Format Marker) files.
//
// USFM is a markup format used for biblical texts, defined at https://docs.usfm.bible/usfm/3.1/index.html.
// This package can parse USFM files into structured Go types and supports various output formats
// including JSON, plain text, and TSV.
//
// Basic usage:
//
//	parser := usfm.NewParser(usfm.DefaultParseOptions())
//	doc, err := parser.Parse(reader, "filename.sfm")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Access parsed content
//	fmt.Printf("Book: %s, Title: %s\n", doc.ID, doc.MainTitle)
//	for _, chapter := range doc.Chapters {
//		fmt.Printf("Chapter %d has %d sections\n", chapter.Number, len(chapter.Sections))
//	}
package usfm

import "time"

// Document represents a complete USFM document containing all parsed content.
// It includes metadata, structure, and the full hierarchy of chapters, sections, and verses.
type Document struct {
	// Identification and metadata
	ID           string            `json:"id"`           // Book identification from \id marker
	Header       string            `json:"header"`       // Running header from \h marker
	TableOfContents []TOCEntry     `json:"toc"`          // Table of contents entries from \toc1, \toc2, \toc3 markers
	MainTitle    string            `json:"main_title"`   // Main title from \mt1 marker
	
	// Content structure
	Chapters     []Chapter         `json:"chapters"`     // All chapters in the book
	
	// Metadata
	ParsedAt     time.Time         `json:"parsed_at"`    // When the document was parsed
	SourceFile   string            `json:"source_file"`  // Original file path
}

// TOCEntry represents a table of contents entry from \toc1, \toc2, or \toc3 markers.
type TOCEntry struct {
	Level int    `json:"level"` // TOC level: 1, 2, or 3
	Text  string `json:"text"`  // The table of contents text
}

// Chapter represents a chapter within a book, identified by a \c marker.
// Each chapter contains one or more sections which in turn contain verses.
type Chapter struct {
	Number   int       `json:"number"`   // Chapter number from \c marker
	Sections []Section `json:"sections"` // Sections within the chapter
}

// Section represents a section within a chapter, typically marked by \s1, \s2, or \s3.
// Sections group related verses together and may include cross-references.
type Section struct {
	Level   int     `json:"level"`   // Section level: 1 (\s1), 2 (\s2), or 3 (\s3)
	Title   string  `json:"title"`   // Section title text
	Reference string `json:"reference,omitempty"` // Cross-reference text from \r marker
	Verses  []Verse `json:"verses"`  // Verses contained in this section
}

// Verse represents a single verse from a \v marker.
// Verses contain the main biblical text and may include footnotes.
type Verse struct {
	Number    int        `json:"number"`    // Verse number from \v marker
	Text      string     `json:"text"`      // Main verse text with footnotes removed
	Footnotes []Footnote `json:"footnotes,omitempty"` // Footnotes extracted from the text
}

// Footnote represents a footnote within a verse, marked by \f...\f* tags.
// Footnotes provide additional information about the biblical text.
type Footnote struct {
	Caller    string `json:"caller"`    // Footnote caller symbol (usually "+")
	Reference string `json:"reference"` // Reference text from \fr marker
	Text      string `json:"text"`      // Footnote content from \ft marker
}

// Marker represents a parsed USFM marker with its content.
// This is used internally during parsing to represent any \marker found in the text.
type Marker struct {
	Tag     string `json:"tag"`     // The marker tag (without backslash, e.g., "c", "v", "s1")
	Content string `json:"content"` // Content following the marker
	Line    int    `json:"line"`    // Line number in source file (for error reporting)
}

// ParseOptions configures the behavior of the USFM parser.
// These options control how strictly the parser validates input and what content it includes.
type ParseOptions struct {
	StrictMode     bool // Whether to fail on unknown/unrecognized markers
	IncludeFootnotes bool // Whether to parse and extract footnotes from verse text
	IncludeReferences bool // Whether to parse cross-reference markers (\r)
}

// DefaultParseOptions returns sensible default parsing options.
// The defaults enable footnote and reference parsing while using lenient mode
// that ignores unknown markers rather than failing.
func DefaultParseOptions() ParseOptions {
	return ParseOptions{
		StrictMode:       false, // Lenient mode - ignore unknown markers
		IncludeFootnotes: true,  // Parse footnotes from verse text
		IncludeReferences: true, // Include cross-references from \r markers
	}
}