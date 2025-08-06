package formatter

import (
	"fmt"
	"strings"

	"github.com/arenzana/usfmp/pkg/usfm"
)

// FormatTSV formats USFM documents as Tab-Separated Values for data analysis.
// Each verse becomes a row with columns for book, chapter, verse, section info,
// verse text, footnotes, and references.
//
// Column format:
//
//	Book	Chapter	Verse	Section_Title	Section_Level	Verse_Text	Footnotes	References
//
// Footnotes are formatted as "caller:reference=text" and separated by semicolons.
// Text is cleaned of tabs and newlines to ensure proper TSV format.
// Multiple documents are included in the same output with their respective book IDs.
func FormatTSV(documents []*usfm.Document) (string, error) {
	var result strings.Builder

	// Header row
	result.WriteString("Book\tChapter\tVerse\tSection_Title\tSection_Level\tVerse_Text\tFootnotes\tReferences\n")

	for _, doc := range documents {
		bookID := doc.ID
		if bookID == "" {
			bookID = "UNKNOWN"
		}

		for _, chapter := range doc.Chapters {
			for _, section := range chapter.Sections {
				sectionTitle := cleanTSVField(section.Title)
				sectionLevel := fmt.Sprintf("%d", section.Level)
				references := cleanTSVField(section.Reference)

				for _, verse := range section.Verses {
					verseText := cleanTSVField(verse.Text)

					// Format footnotes
					var footnotes []string
					for _, footnote := range verse.Footnotes {
						footnoteStr := fmt.Sprintf("%s:%s=%s",
							footnote.Caller, footnote.Reference, footnote.Text)
						footnotes = append(footnotes, footnoteStr)
					}
					footnotesField := cleanTSVField(strings.Join(footnotes, "; "))

					// Write TSV row
					result.WriteString(fmt.Sprintf("%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\n",
						bookID,
						chapter.Number,
						verse.Number,
						sectionTitle,
						sectionLevel,
						verseText,
						footnotesField,
						references,
					))
				}
			}
		}
	}

	return result.String(), nil
}

// cleanTSVField cleans text for safe TSV output by removing tabs and newlines.
// It replaces tabs, newlines, and carriage returns with spaces, then collapses
// multiple consecutive spaces into single spaces. This ensures the text can be
// safely used in TSV format without breaking the column structure.
func cleanTSVField(text string) string {
	// Replace tabs and newlines with spaces
	cleaned := strings.ReplaceAll(text, "\t", " ")
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\r", " ")

	// Collapse multiple spaces into single space
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	return strings.TrimSpace(cleaned)
}
