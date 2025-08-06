# USFM Parser (usfmp)

[![Go Reference](https://pkg.go.dev/badge/github.com/arenzana/usfmp.svg)](https://pkg.go.dev/github.com/arenzana/usfmp)
[![Go Report Card](https://goreportcard.com/badge/github.com/arenzana/usfmp)](https://goreportcard.com/report/github.com/arenzana/usfmp)

A comprehensive Go parser for USFM (Unified Standard Format Marker) files used in biblical texts. This tool provides both a command-line interface and a Go library for parsing USFM files into structured data with multiple output formats.

## Features

- 🔍 **Comprehensive USFM Support**: Parses all major USFM 3.1 markers including chapters, sections, verses, footnotes, and cross-references
- 📖 **Multiple Output Formats**: JSON, plain text, TSV, and PDF (planned)
- 🛠️ **CLI and Library**: Use as a standalone command-line tool or integrate as a Go library
- ⚡ **High Performance**: Efficient parsing with pre-compiled regular expressions
- 🔧 **Flexible Configuration**: Strict vs. lenient parsing modes, optional footnote/reference extraction
- ✅ **Well Tested**: Comprehensive test suite with real biblical text samples
- 📚 **Rich Documentation**: Complete GoDoc documentation with examples

## Quick Start

### Installation

```bash
# Install the CLI tool
go install github.com/arenzana/usfmp/cmd/usfmp@latest

# Or build from source
git clone https://github.com/arenzana/usfmp
cd usfmp
make build
```

### Command Line Usage

```bash
# Parse a single USFM file to JSON
usfmp -f json genesis.sfm

# Parse entire directory to readable text
usfmp -f txt biblical-texts/

# Generate TSV for data analysis
usfmp -f tsv --output analysis.tsv biblical-texts/

# Strict parsing mode (fail on unknown markers)
usfmp --strict -f json genesis.sfm

# Quiet mode (suppress info messages)
usfmp --quiet -f json genesis.sfm > output.json
```

### Library Usage

```go
package main

import (
    "fmt"
    "os"
    "github.com/arenzana/usfmp/pkg/usfm"
)

func main() {
    // Open USFM file
    file, err := os.Open("genesis.sfm")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Create parser with default options
    parser := usfm.NewParser(usfm.DefaultParseOptions())
    
    // Parse the document
    doc, err := parser.Parse(file, "genesis.sfm")
    if err != nil {
        panic(err)
    }

    // Access parsed content
    fmt.Printf("Book: %s\n", doc.ID)
    fmt.Printf("Title: %s\n", doc.MainTitle)
    fmt.Printf("Chapters: %d\n", len(doc.Chapters))
    
    // Iterate through structure
    for _, chapter := range doc.Chapters {
        fmt.Printf("Chapter %d: %d sections\n", chapter.Number, len(chapter.Sections))
        for _, section := range chapter.Sections {
            fmt.Printf("  Section: %s (%d verses)\n", section.Title, len(section.Verses))
            for _, verse := range section.Verses {
                fmt.Printf("    %d: %s\n", verse.Number, verse.Text)
                for _, footnote := range verse.Footnotes {
                    fmt.Printf("      Footnote: %s\n", footnote.Text)
                }
            }
        }
    }
}
```

## USFM Format Support

This parser supports the major USFM 3.1 markers:

### Document Structure
- `\id` - Book identification
- `\h` - Running header text  
- `\toc1`, `\toc2`, `\toc3` - Table of contents entries
- `\mt1` - Main title

### Content Structure  
- `\c` - Chapter numbers
- `\s1`, `\s2`, `\s3` - Section headings (multiple levels)
- `\r` - Cross-references
- `\v` - Verse numbers and text

### Footnotes
- `\f...\\f*` - Footnote blocks
- `\fr` - Footnote reference
- `\ft` - Footnote text

### Parsing Modes

The parser supports different modes for handling edge cases:

```go
// Strict mode - fails on unknown markers
options := usfm.ParseOptions{
    StrictMode:        true,
    IncludeFootnotes:  true,
    IncludeReferences: true,
}

// Lenient mode - ignores unknown markers (default)
options := usfm.DefaultParseOptions()
```

## Output Formats

### JSON Format
Structured JSON with full document hierarchy:

```json
{
  "id": "GEN - Berean Standard Bible",
  "header": "Genesis", 
  "main_title": "Genesis",
  "chapters": [
    {
      "number": 1,
      "sections": [
        {
          "level": 1,
          "title": "The Creation",
          "reference": "(John 1:1–5; Hebrews 11:1–3)",
          "verses": [
            {
              "number": 1,
              "text": "In the beginning God created the heavens and the earth.",
              "footnotes": []
            }
          ]
        }
      ]
    }
  ]
}
```

### Text Format
Human-readable text with proper formatting:

```
Genesis
-------

Book: GEN - Berean Standard Bible
Header: Genesis

Chapter 1
--------------------

The Creation
(John 1:1–5; Hebrews 11:1–3)

1. In the beginning God created the heavens and the earth.
2. Now the earth was formless and void...
```

### TSV Format
Tab-separated values for data analysis:

```
Book	Chapter	Verse	Section_Title	Section_Level	Verse_Text	Footnotes	References
GEN	1	1	The Creation	1	In the beginning God created...		(John 1:1–5)
GEN	1	2	The Creation	1	Now the earth was formless...		(John 1:1–5)
```

## Development

### Building

```bash
# Build the CLI binary
make build

# Run tests
make test

# Generate coverage report  
make coverage

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Development build (format, vet, test, build)
make dev
```

### Testing with Sample Data

The repository includes sample biblical texts for testing:

```bash
# Test with Genesis
make test-genesis

# Run sample data through parser
make run-sample
```

### Available Make Targets

- `make build` - Build the CLI binary
- `make test` - Run all unit tests
- `make coverage` - Generate test coverage report
- `make lint` - Run code linter
- `make fmt` - Format source code
- `make clean` - Remove build artifacts
- `make build-all` - Build for multiple platforms
- `make install` - Install binary to Go bin directory

## API Documentation

Complete API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/arenzana/usfmp).

### Key Types

- [`Document`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#Document) - Complete USFM document
- [`Chapter`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#Chapter) - Book chapter with sections
- [`Section`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#Section) - Thematic section with verses
- [`Verse`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#Verse) - Individual verse with footnotes
- [`ParseOptions`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#ParseOptions) - Parser configuration

### Key Functions

- [`NewParser(options)`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#NewParser) - Create new parser
- [`Parse(reader, filename)`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#Parser.Parse) - Parse USFM content
- [`DefaultParseOptions()`](https://pkg.go.dev/github.com/arenzana/usfmp/pkg/usfm#DefaultParseOptions) - Get default options

## Examples

### Parse Multiple Files

```go
func parseDirectory(dirPath string) ([]*usfm.Document, error) {
    var documents []*usfm.Document
    parser := usfm.NewParser(usfm.DefaultParseOptions())
    
    files, err := filepath.Glob(filepath.Join(dirPath, "*.sfm"))
    if err != nil {
        return nil, err
    }
    
    for _, filename := range files {
        file, err := os.Open(filename)
        if err != nil {
            continue
        }
        
        doc, err := parser.Parse(file, filename)
        file.Close()
        
        if err != nil {
            return nil, fmt.Errorf("parsing %s: %w", filename, err)
        }
        
        documents = append(documents, doc)
    }
    
    return documents, nil
}
```

### Extract Verses by Chapter

```go
func getChapterVerses(doc *usfm.Document, chapterNum int) []usfm.Verse {
    for _, chapter := range doc.Chapters {
        if chapter.Number == chapterNum {
            var verses []usfm.Verse
            for _, section := range chapter.Sections {
                verses = append(verses, section.Verses...)
            }
            return verses
        }
    }
    return nil
}
```

### Custom Formatting

```go
func formatVerseList(verses []usfm.Verse) string {
    var result strings.Builder
    for _, verse := range verses {
        result.WriteString(fmt.Sprintf("%d. %s\n", verse.Number, verse.Text))
        for _, footnote := range verse.Footnotes {
            result.WriteString(fmt.Sprintf("   Note: %s\n", footnote.Text))
        }
    }
    return result.String()
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Run the test suite (`make test`)
5. Format your code (`make fmt`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Write tests for new functionality
- Follow Go conventions and best practices
- Add comprehensive documentation
- Ensure all tests pass
- Run `make lint` before submitting

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [USFM Specification](https://docs.usfm.bible/usfm/3.1/index.html) - Official USFM documentation
- [Berean Standard Bible](https://bereanbible.com/) - Sample text for testing
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## Related Projects

- [usfm-grammar](https://github.com/Bridgeconn/usfm-grammar) - JavaScript USFM parser
- [python-usfm](https://github.com/unfoldingWord-dev/python-usfm) - Python USFM tools

---

**USFM Parser** - Making biblical text processing simple and powerful in Go.