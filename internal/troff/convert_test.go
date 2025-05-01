package troff_test

import (
	"reflect"
	"testing"

	"github.com/carldunham/useful-cookery/internal/troff"
)

func TestConvertMultilineCodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "italic text",
			input:    "This is \\fIitalic\\fP text",
			expected: "This is <i>italic</i> text",
		},
		{
			name:     "bold text",
			input:    "This is \\fBbold\\fP text",
			expected: "This is <b>bold</b> text",
		},
		{
			name:     "mixed formatting",
			input:    "\\fBBold\\fP and \\fIitalic\\fP text",
			expected: "<b>Bold</b> and <i>italic</i> text",
		},
		{
			name:     "with fR terminator",
			input:    "\\fIitalic\\fR text",
			expected: "<i>italic</i> text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := troff.ConvertMultilineCodes(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertMultilineCodes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertCodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "fractions",
			input:    "\\\\(14 cup and \\\\(12 teaspoon",
			expected: "\\ 1/4 cup and \\ 1/2 teaspoon",
		},
		{
			name:     "special characters",
			input:    "2 \\\\(mu 3 equals 6",
			expected: "2 \\&times; 3 equals 6",
		},
		{
			name:     "accented characters",
			input:    "caf\\z\\\\(aae",
			expected: "caf\\z\\\\(aae",
		},
		{
			name:     "quotes",
			input:    "``quoted text''",
			expected: "&ldquo;quoted text&rdquo;",
		},
		{
			name:     "paragraph command",
			input:    ".PP",
			expected: "<p>",
		},
		{
			name:     "line break command",
			input:    ".br",
			expected: "<br>",
		},
		{
			name:     "preformatted start command",
			input:    ".nf",
			expected: "{{pre}}",
		},
		{
			name:     "preformatted end command",
			input:    ".fi",
			expected: "{{endpre}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := troff.ConvertCodes(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertCodes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCollapseLines(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "simple paragraph",
			input:    []string{"Line 1", "Line 2", "Line 3"},
			expected: []string{"Line 1 Line 2 Line 3"},
		},
		{
			name:     "with preformatted text",
			input:    []string{"Before pre", "{{pre}}", "Pre line 1", "Pre line 2", "{{endpre}}", "After pre"},
			expected: []string{"Before pre", "{{pre}}", "Pre line 1", "Pre line 2", "{{endpre}}", "After pre"},
		},
		{
			name:     "multiple paragraphs with pre",
			input:    []string{"Para 1 line 1", "Para 1 line 2", "{{pre}}", "Pre text", "{{endpre}}", "Para 2"},
			expected: []string{"Para 1 line 1 Para 1 line 2", "{{pre}}", "Pre text", "{{endpre}}", "Para 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := troff.CollapseLines(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CollapseLines(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessText(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text with formatting",
			input:    "This is \\fBbold\\fP and \\fIitalic\\fP text",
			expected: "This is <b>bold</b> and <i>italic</i> text",
		},
		{
			name:     "multiline with formatting",
			input:    "Line 1 with \\fBbold\\fP\nLine 2 with \\fIitalic\\fP",
			expected: "Line 1 with <b>bold</b> Line 2 with <i>italic</i>",
		},
		{
			name:     "with preformatted text",
			input:    "Normal text\n.nf\nPreformatted \\fBbold\\fP\n.fi\nMore text",
			expected: "Normal text\n{{pre}}\nPreformatted <b>bold</b>\n{{endpre}}\nMore text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := troff.ProcessText(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
