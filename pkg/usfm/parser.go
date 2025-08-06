package usfm

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser handles the parsing of USFM content into structured Document objects.
// It uses regular expressions to identify and parse various USFM markers and their content.
// The parser behavior can be configured through ParseOptions.
type Parser struct {
	options ParseOptions

	// Regular expressions for parsing different USFM elements
	markerRegex   *regexp.Regexp // Matches any USFM marker (\tag content)
	chapterRegex  *regexp.Regexp // Matches chapter markers (\c number)
	verseRegex    *regexp.Regexp // Matches verse markers (\v number text)
	footnoteRegex *regexp.Regexp // Matches footnote content (\f...\f*)
}

// NewParser creates a new USFM parser with the specified options.
// The parser pre-compiles regular expressions for efficient parsing.
//
// Example:
//
//	options := usfm.ParseOptions{
//		StrictMode: true,
//		IncludeFootnotes: true,
//		IncludeReferences: true,
//	}
//	parser := usfm.NewParser(options)
func NewParser(options ParseOptions) *Parser {
	return &Parser{
		options:       options,
		markerRegex:   regexp.MustCompile(`^\\([a-z0-9]+)\*?\s*(.*)`),
		chapterRegex:  regexp.MustCompile(`^\\c\s+(\d+)`),
		verseRegex:    regexp.MustCompile(`^\\v\s+(\d+)\s*(.*)`),
		footnoteRegex: regexp.MustCompile(`\\f\s*([^\\]*?)\\fr\s*([^\\]*?)\\ft\s*([^\\]*?)\\f\*`),
	}
}

// Parse parses USFM content from an io.Reader and returns a structured Document.
// The sourceFile parameter is used for metadata and error reporting.
//
// The parser processes the input line by line, identifying USFM markers and
// building a hierarchical structure of chapters, sections, and verses.
// Footnotes and cross-references are extracted based on the parser options.
//
// Returns an error if:
//   - Invalid marker syntax is encountered in strict mode
//   - Unknown markers are found in strict mode
//   - Malformed chapter or verse numbers are found
//   - IO errors occur while reading
//
// Example:
//
//	file, err := os.Open("genesis.sfm")
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	parser := usfm.NewParser(usfm.DefaultParseOptions())
//	doc, err := parser.Parse(file, "genesis.sfm")
//	if err != nil {
//		return err
//	}
func (p *Parser) Parse(reader io.Reader, sourceFile string) (*Document, error) {
	doc := &Document{
		ParsedAt:   time.Now(),
		SourceFile: sourceFile,
		Chapters:   make([]Chapter, 0),
	}

	scanner := bufio.NewScanner(reader)
	lineNumber := 0

	var currentChapter *Chapter
	var currentSection *Section

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse marker
		marker, err := p.parseMarker(line, lineNumber)
		if err != nil {
			if p.options.StrictMode {
				return nil, fmt.Errorf("line %d: %w", lineNumber, err)
			}
			continue // Skip invalid markers in non-strict mode
		}

		// Handle different marker types
		switch marker.Tag {
		case "id":
			doc.ID = marker.Content
		case "h":
			doc.Header = marker.Content
		case "toc1":
			doc.TableOfContents = append(doc.TableOfContents, TOCEntry{Level: 1, Text: marker.Content})
		case "toc2":
			doc.TableOfContents = append(doc.TableOfContents, TOCEntry{Level: 2, Text: marker.Content})
		case "toc3":
			doc.TableOfContents = append(doc.TableOfContents, TOCEntry{Level: 3, Text: marker.Content})
		case "mt1":
			doc.MainTitle = marker.Content
		case "c":
			chapterNum, err := p.parseChapter(marker.Content)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNumber, err)
			}

			// Save current section to current chapter before switching chapters
			if currentSection != nil && currentChapter != nil {
				currentChapter.Sections = append(currentChapter.Sections, *currentSection)
			}

			// Save previous chapter if it exists
			if currentChapter != nil {
				doc.Chapters = append(doc.Chapters, *currentChapter)
			}

			// Start new chapter
			currentChapter = &Chapter{
				Number:   chapterNum,
				Sections: make([]Section, 0),
			}
			currentSection = nil
		case "s1", "s2", "s3":
			level := p.getSectionLevel(marker.Tag)
			section := Section{
				Level:  level,
				Title:  marker.Content,
				Verses: make([]Verse, 0),
			}

			// Add previous section to chapter if exists
			if currentSection != nil && currentChapter != nil {
				currentChapter.Sections = append(currentChapter.Sections, *currentSection)
			}

			currentSection = &section
		case "r":
			// Reference/cross-reference - attach to current section
			if p.options.IncludeReferences && currentSection != nil {
				currentSection.Reference = marker.Content
			}
		case "v":
			verse, err := p.parseVerse(marker.Content, p.options.IncludeFootnotes)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNumber, err)
			}

			// Ensure we have a section to add the verse to
			if currentSection == nil {
				currentSection = &Section{
					Level:  1,
					Title:  "",
					Verses: make([]Verse, 0),
				}
			}

			currentSection.Verses = append(currentSection.Verses, *verse)
		default:
			// Handle unknown markers in strict mode
			if p.options.StrictMode {
				return nil, fmt.Errorf("line %d: unknown marker '\\%s'", lineNumber, marker.Tag)
			}
			// In non-strict mode, ignore unknown markers
		}
	}

	// Add final section and chapter
	if currentSection != nil && currentChapter != nil {
		currentChapter.Sections = append(currentChapter.Sections, *currentSection)
	}
	if currentChapter != nil {
		doc.Chapters = append(doc.Chapters, *currentChapter)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return doc, nil
}

// parseMarker extracts marker information from a line
func (p *Parser) parseMarker(line string, lineNumber int) (*Marker, error) {
	if !strings.HasPrefix(line, "\\") {
		return nil, fmt.Errorf("line does not start with marker")
	}

	matches := p.markerRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid marker format")
	}

	return &Marker{
		Tag:     matches[1],
		Content: strings.TrimSpace(matches[2]),
		Line:    lineNumber,
	}, nil
}

// parseChapter extracts chapter number from chapter marker content
func (p *Parser) parseChapter(content string) (int, error) {
	num, err := strconv.Atoi(strings.TrimSpace(content))
	if err != nil {
		return 0, fmt.Errorf("invalid chapter number: %w", err)
	}
	return num, nil
}

// parseVerse extracts verse number and content, including footnotes if enabled
func (p *Parser) parseVerse(content string, includeFootnotes bool) (*Verse, error) {
	parts := strings.SplitN(content, " ", 2)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid verse format")
	}

	verseNum, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid verse number: %w", err)
	}

	verseText := ""
	if len(parts) > 1 {
		verseText = parts[1]
	}

	verse := &Verse{
		Number:    verseNum,
		Text:      verseText,
		Footnotes: make([]Footnote, 0),
	}

	// Parse footnotes if enabled
	if includeFootnotes {
		footnotes := p.extractFootnotes(verseText)
		verse.Footnotes = footnotes

		// Remove footnote markers from main text
		verse.Text = p.removeFootnoteMarkers(verseText)
	}

	return verse, nil
}

// getSectionLevel returns the numeric level for section markers
func (p *Parser) getSectionLevel(tag string) int {
	switch tag {
	case "s1":
		return 1
	case "s2":
		return 2
	case "s3":
		return 3
	default:
		return 1
	}
}

// extractFootnotes finds and extracts footnotes from verse text
func (p *Parser) extractFootnotes(text string) []Footnote {
	footnotes := make([]Footnote, 0)

	matches := p.footnoteRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			footnote := Footnote{
				Caller:    strings.TrimSpace(match[1]),
				Reference: strings.TrimSpace(match[2]),
				Text:      strings.TrimSpace(match[3]),
			}
			footnotes = append(footnotes, footnote)
		}
	}

	return footnotes
}

// removeFootnoteMarkers removes footnote markup from text, leaving clean readable text
func (p *Parser) removeFootnoteMarkers(text string) string {
	// Remove footnote markers but keep the main text clean
	cleaned := p.footnoteRegex.ReplaceAllString(text, "")

	// Clean up multiple spaces that may result from footnote removal
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	return strings.TrimSpace(cleaned)
}
