package troff

import (
	"regexp"
	"strings"
)

// ConvertCodes replaces TROFF formatting codes with HTML/CSS equivalents.
// This is a Go implementation of the Python _convertCodes function.
func ConvertCodes(content string) string {
	// First convert multiline codes
	content = ConvertMultilineCodes(content)

	// Handle special cases for line-based commands first
	switch {
	case strings.HasPrefix(content, ".PP"):
		return "<p>"
	case strings.HasPrefix(content, ".br"):
		return "<br>"
	case strings.HasPrefix(content, ".nf"):
		return "{{pre}}"
	case strings.HasPrefix(content, ".fi"):
		return "{{endpre}}"
	case strings.HasPrefix(content, ".RS"):
		return `<div class="troff-RS">`
	case strings.HasPrefix(content, ".RE"):
		return "</div>"
	}

	// Define simple string replacements
	stringReplacements := map[string]string{
		"``": "&ldquo;",
		"''": "&rdquo;",
		"\t": " ",
	}

	// Apply string replacements
	for pattern, replacement := range stringReplacements {
		content = strings.ReplaceAll(content, pattern, replacement)
	}

	// Define regex replacements for more complex patterns
	regexReplacements := []struct {
		pattern     string
		replacement string
	}{
		{`\\\(18`, " 1/8"},
		{`\\\(14`, " 1/4"},
		{`\\\(13`, " 1/3"},
		{`\\\(12`, " 1/2"},
		{`\\\(34`, " 3/4"},
		{`\\\(mu`, "&times;"},
		{`\\\(em`, "&mdash;"},
		{`\\-`, "-"},
		{`\\z\\\(aae`, "&eacute;"},
		{`\\z\\\(aao`, "&oacute;"},
		{`\\z\\\(gaa`, "&agrave;"},
		{`\\o'e\\\(aa`, "&eacute;"},
		{`\\\*\:o`, "&ouml;"},
		{`\\o'o\"'`, "&ouml;"},
		{`\\s\-2`, `<span class="troff-s-2">`},
		{`\\s0`, "</span>"},
	}

	// Apply regex replacements
	for _, r := range regexReplacements {
		re := regexp.MustCompile(r.pattern)
		content = re.ReplaceAllString(content, r.replacement)
	}

	return strings.TrimSpace(content)
}

// ConvertMultilineCodes handles multiline formatting codes like \fI...\fP for italics.
// This is a Go implementation of the Python _convertMultilineCodes function.
func ConvertMultilineCodes(content string) string {
	// Define multiline replacements
	patterns := []struct {
		regex       string
		replacement string
	}{
		{`\\fI(.+?)\\f[PR]`, "<i>$1</i>"},
		{`\\fB(.+?)\\f[PR]`, "<b>$1</b>"},
	}

	// Apply replacements
	for _, p := range patterns {
		re := regexp.MustCompile(p.regex)
		content = re.ReplaceAllString(content, p.replacement)
	}

	return strings.TrimSpace(content)
}

// CollapseLines collapses multiple lines into paragraphs, handling special cases for preformatted text.
// This is a Go implementation of the Python _collapseLines function.
func CollapseLines(lines []string) []string {
	var result []string
	var lineBuf string
	pre := false

	for _, nextLine := range lines {
		if nextLine == "{{pre}}" {
			pre = true

			if lineBuf != "" {
				result = append(result, lineBuf)
				lineBuf = ""
			}

			result = append(result, nextLine)
			continue
		} else if nextLine == "{{endpre}}" {
			pre = false

			if lineBuf != "" {
				result = append(result, lineBuf)
				lineBuf = ""
			}

			result = append(result, nextLine)
			continue
		}

		switch {
		case pre:
			result = append(result, nextLine)
		case lineBuf == "":
			lineBuf = nextLine
		default:
			lineBuf += " " + nextLine
		}
	}

	if lineBuf != "" {
		result = append(result, lineBuf)
	}

	return result
}

// ProcessText processes a text string by converting TROFF codes to HTML.
// This is a convenience function that combines ConvertCodes and CollapseLines.
func ProcessText(text string) string {
	// First convert the entire text with multiline codes
	text = ConvertMultilineCodes(text)

	// Split the text into lines
	lines := strings.Split(text, "\n")

	// Convert each line for other codes
	convertedLines := make([]string, len(lines))
	for i, line := range lines {
		convertedLines[i] = ConvertCodes(line)
	}

	// Collapse the lines
	collapsed := CollapseLines(convertedLines)

	// Join the collapsed lines with newlines
	return strings.Join(collapsed, "\n")
}
