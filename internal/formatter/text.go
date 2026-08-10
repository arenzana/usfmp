package formatter

import (
	"fmt"
	"strings"

	"github.com/arenzana/usfmp/pkg/usfm"
)

// FormatText formats USFM documents as human-readable plain text.
// The output includes document titles, chapter numbers, section headings,
// verse numbers, and footnotes in brackets.
//
// The text format uses:
//   - Document title with underline
//   - "Chapter N" headings with dashes
//   - Section titles with indentation based on level
//   - "N. verse text" format for verses
//   - "[caller:reference - footnote text]" for footnotes
//   - Cross-references in parentheses after section titles
//
// Multiple documents are separated by a line of equal signs.
func FormatText(documents []*usfm.Document) (string, error) {
	var result strings.Builder

	for i, doc := range documents {
		if i > 0 {
			result.WriteString("\n" + strings.Repeat("=", 80) + "\n\n")
		}

		// Document header
		if doc.MainTitle != "" {
			fmt.Fprintf(&result, "%s\n", doc.MainTitle)
			result.WriteString(strings.Repeat("-", len(doc.MainTitle)) + "\n\n")
		}

		if doc.ID != "" {
			fmt.Fprintf(&result, "Book: %s\n", doc.ID)
		}

		if doc.Header != "" {
			fmt.Fprintf(&result, "Header: %s\n", doc.Header)
		}

		result.WriteString("\n")

		// Chapters
		for _, chapter := range doc.Chapters {
			fmt.Fprintf(&result, "Chapter %d\n", chapter.Number)
			result.WriteString(strings.Repeat("-", 20) + "\n\n")

			// Sections
			for _, section := range chapter.Sections {
				if section.Title != "" {
					// Add indent based on section level
					indent := strings.Repeat("  ", section.Level-1)
					fmt.Fprintf(&result, "%s%s\n", indent, section.Title)

					if section.Reference != "" {
						fmt.Fprintf(&result, "%s(%s)\n", indent, section.Reference)
					}
					result.WriteString("\n")
				}

				// Verses
				for _, verse := range section.Verses {
					fmt.Fprintf(&result, "%d. %s", verse.Number, verse.Text)

					// Add footnotes
					if len(verse.Footnotes) > 0 {
						result.WriteString(" [")
						for j, footnote := range verse.Footnotes {
							if j > 0 {
								result.WriteString("; ")
							}
							fmt.Fprintf(&result, "%s:%s - %s",
								footnote.Caller, footnote.Reference, footnote.Text)
						}
						result.WriteString("]")
					}

					result.WriteString("\n")
				}

				result.WriteString("\n")
			}
		}
	}

	return result.String(), nil
}
