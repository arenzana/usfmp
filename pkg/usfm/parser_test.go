package usfm

import (
	"strings"
	"testing"
)

// TestParseBasicUSFM tests parsing a basic USFM document
func TestParseBasicUSFM(t *testing.T) {
	input := `\id GEN - Test Bible
\h Genesis
\toc1 Genesis
\toc2 Genesis
\mt1 Genesis
\c 1
\s1 The Creation
\r (John 1:1–5; Hebrews 11:1–3)
\v 1 In the beginning God created the heavens and the earth.
\v 2 Now the earth was formless and void, and darkness was over the surface of the deep.
\c 2
\s1 The Garden of Eden
\v 1 Thus the heavens and the earth were completed in all their vast array.`

	parser := NewParser(DefaultParseOptions())
	doc, err := parser.Parse(strings.NewReader(input), "test.sfm")
	
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	// Test document metadata
	if doc.ID != "GEN - Test Bible" {
		t.Errorf("Expected ID 'GEN - Test Bible', got '%s'", doc.ID)
	}
	
	if doc.Header != "Genesis" {
		t.Errorf("Expected header 'Genesis', got '%s'", doc.Header)
	}
	
	if doc.MainTitle != "Genesis" {
		t.Errorf("Expected main title 'Genesis', got '%s'", doc.MainTitle)
	}
	
	// Test table of contents
	if len(doc.TableOfContents) != 2 {
		t.Errorf("Expected 2 TOC entries, got %d", len(doc.TableOfContents))
	}
	
	// Test chapters
	if len(doc.Chapters) != 2 {
		t.Errorf("Expected 2 chapters, got %d", len(doc.Chapters))
	}
	
	// Test first chapter
	chapter1 := doc.Chapters[0]
	if chapter1.Number != 1 {
		t.Errorf("Expected chapter number 1, got %d", chapter1.Number)
	}
	
	if len(chapter1.Sections) != 1 {
		t.Errorf("Expected 1 section in chapter 1, got %d", len(chapter1.Sections))
	}
	
	// Test section
	section := chapter1.Sections[0]
	if section.Title != "The Creation" {
		t.Errorf("Expected section title 'The Creation', got '%s'", section.Title)
	}
	
	if section.Reference != "(John 1:1–5; Hebrews 11:1–3)" {
		t.Errorf("Expected section reference '(John 1:1–5; Hebrews 11:1–3)', got '%s'", section.Reference)
	}
	
	// Test verses
	if len(section.Verses) != 2 {
		t.Errorf("Expected 2 verses in section, got %d", len(section.Verses))
	}
	
	verse1 := section.Verses[0]
	if verse1.Number != 1 {
		t.Errorf("Expected verse number 1, got %d", verse1.Number)
	}
	
	expectedText := "In the beginning God created the heavens and the earth."
	if verse1.Text != expectedText {
		t.Errorf("Expected verse text '%s', got '%s'", expectedText, verse1.Text)
	}
}

// TestParseFootnotes tests parsing verses with footnotes
func TestParseFootnotes(t *testing.T) {
	input := `\id GEN - Test Bible
\c 1
\v 3 And God said, "Let there be light," \f + \fr 1:3 \ft Cited in 2 Corinthians 4:6\f* and there was light.
\v 5 God called the light "day," and the darkness He called "night." And there was evening, and there was morning—the first day.\f + \fr 1:5 \ft Literally day one\f*`

	parser := NewParser(DefaultParseOptions())
	doc, err := parser.Parse(strings.NewReader(input), "test.sfm")
	
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	// Should have one chapter with verses containing footnotes
	if len(doc.Chapters) != 1 {
		t.Fatalf("Expected 1 chapter, got %d", len(doc.Chapters))
	}
	
	chapter := doc.Chapters[0]
	if len(chapter.Sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(chapter.Sections))
	}
	
	section := chapter.Sections[0]
	if len(section.Verses) != 2 {
		t.Fatalf("Expected 2 verses, got %d", len(section.Verses))
	}
	
	// Test verse 3 footnote
	verse3 := section.Verses[0]
	if len(verse3.Footnotes) != 1 {
		t.Errorf("Expected 1 footnote in verse 3, got %d", len(verse3.Footnotes))
	} else {
		footnote := verse3.Footnotes[0]
		if footnote.Caller != "+" {
			t.Errorf("Expected footnote caller '+', got '%s'", footnote.Caller)
		}
		if footnote.Reference != "1:3" {
			t.Errorf("Expected footnote reference '1:3', got '%s'", footnote.Reference)
		}
		if footnote.Text != "Cited in 2 Corinthians 4:6" {
			t.Errorf("Expected footnote text 'Cited in 2 Corinthians 4:6', got '%s'", footnote.Text)
		}
	}
	
	// Test that footnote markers are removed from main text
	expectedText := `And God said, "Let there be light," and there was light.`
	if verse3.Text != expectedText {
		t.Errorf("Expected cleaned verse text '%s', got '%s'", expectedText, verse3.Text)
	}
}

// TestParseMultipleSections tests parsing multiple section levels
func TestParseMultipleSections(t *testing.T) {
	input := `\id GEN - Test Bible
\c 1
\s1 Major Section
\v 1 First verse.
\s2 Minor Section
\v 2 Second verse.
\s3 Sub-section
\v 3 Third verse.`

	parser := NewParser(DefaultParseOptions())
	doc, err := parser.Parse(strings.NewReader(input), "test.sfm")
	
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	chapter := doc.Chapters[0]
	if len(chapter.Sections) != 3 {
		t.Fatalf("Expected 3 sections, got %d", len(chapter.Sections))
	}
	
	// Test section levels
	levels := []int{1, 2, 3}
	titles := []string{"Major Section", "Minor Section", "Sub-section"}
	
	for i, expectedLevel := range levels {
		if chapter.Sections[i].Level != expectedLevel {
			t.Errorf("Expected section %d level %d, got %d", i, expectedLevel, chapter.Sections[i].Level)
		}
		if chapter.Sections[i].Title != titles[i] {
			t.Errorf("Expected section %d title '%s', got '%s'", i, titles[i], chapter.Sections[i].Title)
		}
	}
}

// TestParseStrictMode tests strict parsing mode
func TestParseStrictMode(t *testing.T) {
	input := `\id GEN - Test Bible
\unknownmarker Some content
\c 1
\v 1 First verse.`

	// Test strict mode - should fail
	strictOptions := ParseOptions{
		StrictMode:        true,
		IncludeFootnotes:  true,
		IncludeReferences: true,
	}
	
	strictParser := NewParser(strictOptions)
	_, err := strictParser.Parse(strings.NewReader(input), "test.sfm")
	
	if err == nil {
		t.Error("Expected error in strict mode with unknown marker, but parsing succeeded")
	}
	
	// Test non-strict mode - should succeed
	lenientParser := NewParser(DefaultParseOptions())
	doc, err := lenientParser.Parse(strings.NewReader(input), "test.sfm")
	
	if err != nil {
		t.Fatalf("Expected success in non-strict mode, but got error: %v", err)
	}
	
	// Should still parse the valid parts
	if doc.ID != "GEN - Test Bible" {
		t.Errorf("Expected ID to be parsed correctly even with unknown marker")
	}
}

// TestParseEmptyInput tests parsing empty input
func TestParseEmptyInput(t *testing.T) {
	parser := NewParser(DefaultParseOptions())
	doc, err := parser.Parse(strings.NewReader(""), "empty.sfm")
	
	if err != nil {
		t.Fatalf("Parse failed on empty input: %v", err)
	}
	
	// Should return valid but empty document
	if len(doc.Chapters) != 0 {
		t.Errorf("Expected 0 chapters for empty input, got %d", len(doc.Chapters))
	}
}

// TestMarkerParsing tests the internal marker parsing functionality
func TestMarkerParsing(t *testing.T) {
	parser := NewParser(DefaultParseOptions())
	
	testCases := []struct {
		input       string
		expectedTag string
		expectedContent string
		shouldError bool
	}{
		{`\id GEN - Test Bible`, "id", "GEN - Test Bible", false},
		{`\c 1`, "c", "1", false},
		{`\v 1 In the beginning...`, "v", "1 In the beginning...", false},
		{`\s1 Section Title`, "s1", "Section Title", false},
		{`Not a marker`, "", "", true},
		{`\`, "", "", true},
	}
	
	for i, tc := range testCases {
		marker, err := parser.parseMarker(tc.input, i+1)
		
		if tc.shouldError {
			if err == nil {
				t.Errorf("Test case %d: expected error but got none", i+1)
			}
			continue
		}
		
		if err != nil {
			t.Errorf("Test case %d: unexpected error: %v", i+1, err)
			continue
		}
		
		if marker.Tag != tc.expectedTag {
			t.Errorf("Test case %d: expected tag '%s', got '%s'", i+1, tc.expectedTag, marker.Tag)
		}
		
		if marker.Content != tc.expectedContent {
			t.Errorf("Test case %d: expected content '%s', got '%s'", i+1, tc.expectedContent, marker.Content)
		}
	}
}