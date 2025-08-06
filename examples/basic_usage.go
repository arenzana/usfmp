package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arenzana/usfmp/internal/formatter"
	"github.com/arenzana/usfmp/pkg/usfm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <usfm-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../bsb_usfm/01GENBSB.SFM\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Example 1: Basic parsing
	fmt.Println("=== Basic USFM Parsing Example ===")
	doc, err := parseUSFMFile(filename)
	if err != nil {
		log.Fatalf("Failed to parse USFM file: %v", err)
	}

	// Display basic information
	fmt.Printf("Book ID: %s\n", doc.ID)
	fmt.Printf("Title: %s\n", doc.MainTitle)
	fmt.Printf("Header: %s\n", doc.Header)
	fmt.Printf("Chapters: %d\n", len(doc.Chapters))
	fmt.Printf("Source: %s\n", doc.SourceFile)
	fmt.Printf("Parsed: %s\n\n", doc.ParsedAt.Format("2006-01-02 15:04:05"))

	// Example 2: Navigate structure
	fmt.Println("=== Structure Navigation Example ===")
	navigateStructure(doc)

	// Example 3: Extract verses by chapter
	fmt.Println("=== Extract Chapter Verses Example ===")
	if len(doc.Chapters) > 0 {
		verses := extractChapterVerses(doc, 1)
		fmt.Printf("Chapter 1 has %d verses:\n", len(verses))
		for i, verse := range verses {
			if i >= 3 { // Show only first 3 verses
				fmt.Printf("  ... and %d more verses\n", len(verses)-3)
				break
			}
			fmt.Printf("  %d: %s\n", verse.Number, limitText(verse.Text, 80))
		}
		fmt.Println()
	}

	// Example 4: Find verses with footnotes
	fmt.Println("=== Footnotes Example ===")
	findVersesWithFootnotes(doc, 5) // Show first 5

	// Example 5: Different output formats
	fmt.Println("=== Output Formats Example ===")
	demonstrateFormats(doc)
}

// parseUSFMFile demonstrates basic parsing
func parseUSFMFile(filename string) (*usfm.Document, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", err)
		}
	}()

	// Create parser with default options
	parser := usfm.NewParser(usfm.DefaultParseOptions())

	// Parse the document
	doc, err := parser.Parse(file, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse USFM: %w", err)
	}

	return doc, nil
}

// navigateStructure demonstrates how to navigate the document structure
func navigateStructure(doc *usfm.Document) {
	for i, chapter := range doc.Chapters {
		if i >= 2 { // Show only first 2 chapters
			fmt.Printf("... and %d more chapters\n", len(doc.Chapters)-2)
			break
		}

		fmt.Printf("Chapter %d (%d sections):\n", chapter.Number, len(chapter.Sections))

		for j, section := range chapter.Sections {
			if j >= 2 { // Show only first 2 sections per chapter
				fmt.Printf("    ... and %d more sections\n", len(chapter.Sections)-2)
				break
			}

			fmt.Printf("  Section %d: \"%s\" (%d verses)\n",
				section.Level, limitText(section.Title, 50), len(section.Verses))

			if section.Reference != "" {
				fmt.Printf("    Reference: %s\n", section.Reference)
			}
		}
	}
	fmt.Println()
}

// extractChapterVerses demonstrates extracting verses from a specific chapter
func extractChapterVerses(doc *usfm.Document, chapterNum int) []usfm.Verse {
	for _, chapter := range doc.Chapters {
		if chapter.Number == chapterNum {
			var allVerses []usfm.Verse
			for _, section := range chapter.Sections {
				allVerses = append(allVerses, section.Verses...)
			}
			return allVerses
		}
	}
	return nil
}

// findVersesWithFootnotes demonstrates finding verses that have footnotes
func findVersesWithFootnotes(doc *usfm.Document, limit int) {
	found := 0

	for _, chapter := range doc.Chapters {
		if found >= limit {
			break
		}

		for _, section := range chapter.Sections {
			if found >= limit {
				break
			}

			for _, verse := range section.Verses {
				if len(verse.Footnotes) > 0 {
					fmt.Printf("Chapter %d, Verse %d (%d footnotes):\n",
						chapter.Number, verse.Number, len(verse.Footnotes))
					fmt.Printf("  Text: %s\n", limitText(verse.Text, 100))

					for _, footnote := range verse.Footnotes {
						fmt.Printf("  Footnote [%s:%s]: %s\n",
							footnote.Caller, footnote.Reference,
							limitText(footnote.Text, 80))
					}
					fmt.Println()

					found++
					if found >= limit {
						break
					}
				}
			}
		}
	}

	if found == 0 {
		fmt.Println("No verses with footnotes found in this document.")
	}
}

// demonstrateFormats shows different output formats
func demonstrateFormats(doc *usfm.Document) {
	documents := []*usfm.Document{doc}

	// JSON format (first 500 characters)
	fmt.Println("JSON Format (excerpt):")
	jsonOutput, err := formatter.FormatJSON(documents)
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
	} else {
		fmt.Printf("%s...\n\n", limitText(jsonOutput, 500))
	}

	// Text format (first 800 characters)
	fmt.Println("Text Format (excerpt):")
	textOutput, err := formatter.FormatText(documents)
	if err != nil {
		fmt.Printf("Error formatting text: %v\n", err)
	} else {
		fmt.Printf("%s...\n\n", limitText(textOutput, 800))
	}

	// TSV format (first 3 lines)
	fmt.Println("TSV Format (first 3 lines):")
	tsvOutput, err := formatter.FormatTSV(documents)
	if err != nil {
		fmt.Printf("Error formatting TSV: %v\n", err)
	} else {
		lines := splitLines(tsvOutput)
		for i, line := range lines {
			if i >= 3 {
				fmt.Printf("... (%d more lines)\n", len(lines)-3)
				break
			}
			fmt.Printf("%s\n", limitText(line, 120))
		}
	}
}

// Helper functions

func limitText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func splitLines(text string) []string {
	lines := []string{}
	current := ""

	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

// init sets up example
func init() {
	// Change to project root if run from examples directory
	if filepath.Base(os.Args[0]) == "basic_usage" || filepath.Base(os.Args[0]) == "basic_usage.exe" {
		if _, err := os.Stat("../bsb_usfm"); err == nil {
			if err := os.Chdir(".."); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to change directory: %v\n", err)
			}
		}
	}
}
