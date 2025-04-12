package troff

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Lexer tokenizes a TROFF file into a stream of tokens.
type Lexer struct {
	reader      *bufio.Reader
	line        int
	column      int
	currentLine string
	linePos     int
	lastCommand string
	tokens      []Token
	peekPos     int
}

// NewLexer creates a new lexer for the given reader.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		reader:  bufio.NewReader(r),
		line:    1,
		column:  0,
		linePos: 0,
		tokens:  make([]Token, 0),
	}
}

// Tokenize processes the entire input and returns all tokens.
func (l *Lexer) Tokenize() ([]Token, error) {
	for {
		token, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		l.tokens = append(l.tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}
	return l.tokens, nil
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() (Token, error) {
	// If we have tokens already tokenized, return the next one
	if l.peekPos < len(l.tokens) {
		token := l.tokens[l.peekPos]
		l.peekPos++
		return token, nil
	}

	// If we're at the end of the current line or haven't read a line yet, read the next line
	if l.currentLine == "" || l.linePos >= len(l.currentLine) {
		if err := l.readNextLine(); err != nil {
			if err == io.EOF {
				return Token{Type: TokenEOF, Line: l.line, Column: l.column}, nil
			}
			return Token{}, err
		}
	}

	// Skip whitespace at the beginning of the line
	l.skipWhitespace()

	// If we're at the end of the line after skipping whitespace, read the next line
	if l.linePos >= len(l.currentLine) {
		return l.NextToken()
	}

	// Check for a command (starts with a period at the beginning of a line)
	if l.linePos == 0 && l.currentLine[0] == '.' {
		return l.readCommand()
	}

	// Otherwise, it's a parameter
	return l.readParam()
}

// PeekToken returns the next token without consuming it.
func (l *Lexer) PeekToken() (Token, error) {
	token, err := l.NextToken()
	if err != nil {
		return Token{}, err
	}
	l.peekPos--
	return token, nil
}

// readNextLine reads the next line from the input.
func (l *Lexer) readNextLine() error {
	line, err := l.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	// Remove trailing newline
	line = strings.TrimRight(line, "\r\n")

	l.currentLine = line
	l.linePos = 0
	l.line++
	l.column = 0

	return nil
}

// skipWhitespace skips whitespace characters in the current line.
func (l *Lexer) skipWhitespace() {
	for l.linePos < len(l.currentLine) && unicode.IsSpace(rune(l.currentLine[l.linePos])) {
		l.linePos++
		l.column++
	}
}

// readCommand reads a command token.
func (l *Lexer) readCommand() (Token, error) {
	// Skip the period
	l.linePos++
	l.column++

	// Read the command (two alpha characters)
	start := l.linePos
	for l.linePos < len(l.currentLine) && l.linePos-start < 2 && unicode.IsLetter(rune(l.currentLine[l.linePos])) {
		l.linePos++
		l.column++
	}

	// If we didn't read exactly two characters, it's not a valid command
	if l.linePos-start != 2 {
		return Token{Type: TokenUnknown, Value: l.currentLine[start-1 : l.linePos], Line: l.line, Column: start}, nil
	}

	command := l.currentLine[start:l.linePos]
	l.lastCommand = command

	return Token{Type: TokenCommand, Value: command, Line: l.line, Column: start}, nil
}

// readParam reads a parameter token.
func (l *Lexer) readParam() (Token, error) {
	start := l.linePos
	inQuote := false
	quoteChar := rune(0)

	// Read until whitespace or end of line, handling quotes
	for l.linePos < len(l.currentLine) {
		ch := rune(l.currentLine[l.linePos])

		// Handle quotes
		if (ch == '"' || ch == '\'') && (l.linePos == 0 || l.currentLine[l.linePos-1] != '\\') {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
			}
		}

		// Break on whitespace if not in quotes
		if unicode.IsSpace(ch) && !inQuote {
			break
		}

		l.linePos++
		l.column++
	}

	// If we're at the end of the line and still in a quote, the quote is unterminated
	if inQuote && l.linePos >= len(l.currentLine) {
		// For simplicity, we'll just treat it as a regular param
		// In a more robust implementation, we might want to handle this differently
	}

	// If we didn't read anything, skip to the next line
	if l.linePos == start {
		if err := l.readNextLine(); err != nil {
			if err == io.EOF {
				return Token{Type: TokenEOF, Line: l.line, Column: l.column}, nil
			}
			return Token{}, err
		}
		return l.NextToken()
	}

	param := l.currentLine[start:l.linePos]

	// Remove quotes if the parameter is quoted
	if len(param) >= 2 && (param[0] == '"' && param[len(param)-1] == '"' || param[0] == '\'' && param[len(param)-1] == '\'') {
		param = param[1 : len(param)-1]
	}

	return Token{Type: TokenParam, Value: param, Command: l.lastCommand, Line: l.line, Column: start}, nil
}
