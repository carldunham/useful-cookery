package troff

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

// Lexer tokenizes a TROFF file into a stream of tokens.
type Lexer struct {
	reader           *bufio.Reader
	line             int
	column           int
	currentLine      string
	linePos          int
	lastCommand      string
	tokens           []Token
	peekPos          int
	lastParamToken   Token
	inMultilineParam bool
	multilineBuffer  string
}

// NewLexer creates a new lexer for the given reader.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		reader:           bufio.NewReader(r),
		line:             1,
		column:           0,
		linePos:          0,
		tokens:           make([]Token, 0),
		inMultilineParam: false,
		multilineBuffer:  "",
	}
}

// Tokenize processes the entire input and returns all tokens.
func (l *Lexer) Tokenize() ([]Token, error) {
	// Reset tokens and peekPos to ensure we start fresh
	l.tokens = make([]Token, 0)
	l.peekPos = 0

	for {
		// Get the next token directly from input, not from the tokens slice
		token, err := l.readNextToken()
		if err != nil {
			return nil, err
		}

		// If we hit EOF, break without adding it to the tokens
		if token.Type == TokenEOF {
			break
		}

		// Add non-EOF tokens to our list
		l.tokens = append(l.tokens, token)
	}
	return l.tokens, nil
}

// readNextToken reads the next token from input without using the tokens slice
func (l *Lexer) readNextToken() (Token, error) {
	for {
		// If we're at the end of the current line or haven't read a line yet, read the next line
		if l.currentLine == "" || l.linePos >= len(l.currentLine) {
			if err := l.readNextLine(); err != nil {
				if errors.Is(err, io.EOF) {
					// If we have a multiline buffer, return it as a token
					if l.multilineBuffer != "" {
						token := Token{
							Type:    TokenParam,
							Value:   l.multilineBuffer,
							Command: l.lastCommand,
							Line:    l.line,
							Column:  0,
						}
						l.multilineBuffer = ""
						l.inMultilineParam = false
						return token, nil
					}

					// Reset multiline param state at EOF
					l.inMultilineParam = false
					return Token{Type: TokenEOF, Line: l.line, Column: l.column}, nil
				}
				return Token{}, err
			}
		}

		// Skip whitespace at the beginning of the line
		l.skipWhitespace()

		// If we're still at the end of the line after skipping whitespace, continue to the next iteration
		// which will read the next line
		if l.linePos >= len(l.currentLine) {
			continue
		}

		// Check for a command (starts with a period at the beginning of a line)
		if l.linePos == 0 && l.currentLine[0] == '.' {
			// If we have a multiline buffer, return it as a token before processing the command
			if l.multilineBuffer != "" {
				token := Token{
					Type:    TokenParam,
					Value:   l.multilineBuffer,
					Command: l.lastCommand,
					Line:    l.line,
					Column:  0,
				}
				l.multilineBuffer = ""
				l.inMultilineParam = false
				return token, nil
			}

			// Otherwise, process the command
			return l.readCommand()
		}

		// If we're at the beginning of a line and not at a command, and we're in a multiline parameter context,
		// add the line to the multiline buffer
		if l.linePos == 0 && l.inMultilineParam {
			if l.multilineBuffer != "" {
				l.multilineBuffer += "\n"
			}
			l.multilineBuffer += l.currentLine
			l.linePos = len(l.currentLine) // Consume the entire line
			continue
		}

		// Otherwise, it's a parameter
		token, err := l.readParam()
		if err != nil {
			return Token{}, err
		}

		// If we got a special "try again" token from readParam, continue the loop
		if token.Type == TokenUnknown && token.Value == "" && token.Line == -1 && token.Column == -1 {
			continue
		}

		// If this is a parameter token, save it and set the multiline flag
		if token.Type == TokenParam {
			token.Command = l.lastCommand
			l.lastParamToken = token
			l.inMultilineParam = true
		}

		return token, nil
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() (Token, error) {
	// If we have tokens already tokenized, return the next one
	if l.peekPos < len(l.tokens) {
		token := l.tokens[l.peekPos]
		l.peekPos++
		return token, nil
	}

	// Otherwise, read a new token from the input
	token, err := l.readNextToken()
	if err != nil {
		return Token{}, err
	}

	// Add it to our tokens slice and increment peekPos
	l.tokens = append(l.tokens, token)
	l.peekPos = len(l.tokens)

	return token, nil
}

// PeekToken returns the next token without consuming it.
func (l *Lexer) PeekToken() (Token, error) {
	if l.peekPos < len(l.tokens) {
		return l.tokens[l.peekPos], nil
	}

	// Read the next token
	token, err := l.NextToken()
	if err != nil {
		return Token{}, err
	}

	// Adjust peekPos so the token isn't consumed
	l.peekPos--

	return token, nil
}

// readNextLine reads the next line from the input.
func (l *Lexer) readNextLine() error {
	line, err := l.reader.ReadString('\n')

	// If we hit EOF and there's no data, return EOF
	if err == io.EOF && line == "" {
		return err
	}

	// For other errors (except EOF with data), return the error
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
	// Check for quoted parameters
	if l.linePos < len(l.currentLine) && (l.currentLine[l.linePos] == '"' || l.currentLine[l.linePos] == '\'') {
		return l.readQuotedParam()
	}

	// Standard parameter handling for all cases
	return l.readStandardParam()
}

// readQuotedParam reads a quoted parameter
func (l *Lexer) readQuotedParam() (Token, error) {
	start := l.linePos
	quoteChar := rune(l.currentLine[l.linePos])

	// Skip the opening quote
	l.linePos++
	l.column++

	// Find the closing quote
	quoteEnd := -1
	for i := l.linePos; i < len(l.currentLine); i++ {
		if l.currentLine[i] == byte(quoteChar) && (i == 0 || l.currentLine[i-1] != '\\') {
			quoteEnd = i
			break
		}
	}

	var param string
	if quoteEnd == -1 {
		// Unterminated quote, treat the rest of the line as the parameter
		param = l.currentLine[l.linePos:]
		l.linePos = len(l.currentLine)
	} else {
		// Extract the parameter without the quotes
		param = l.currentLine[l.linePos:quoteEnd]
		l.linePos = quoteEnd + 1 // Skip the closing quote
	}

	l.column += len(param) + 1 // +1 for the closing quote

	return Token{Type: TokenParam, Value: param, Command: l.lastCommand, Line: l.line, Column: start}, nil
}

// readStandardParam reads a standard (non-quoted) parameter
func (l *Lexer) readStandardParam() (Token, error) {
	start := l.linePos
	inQuote := false
	quoteChar := rune(0)

	// Special case: if we're at the beginning of a line and not at a command,
	// and we're in a multiline parameter context, treat the entire line as a parameter
	if l.linePos == 0 && l.currentLine[0] != '.' && l.inMultilineParam {
		// Read the rest of the line as a single parameter
		param := l.currentLine[start:]
		l.linePos = len(l.currentLine) // Consume the entire line

		return Token{Type: TokenParam, Value: param, Command: l.lastCommand, Line: l.line, Column: start}, nil
	}

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

	// If we didn't read anything, return a special value to signal readNextToken to try again
	if l.linePos == start {
		// Just return an empty token with a special type that readNextToken will recognize
		return Token{Type: TokenUnknown, Value: "", Line: -1, Column: -1}, nil
	}

	param := l.currentLine[start:l.linePos]

	// Remove quotes if the parameter is quoted
	if len(param) >= 2 && (param[0] == '"' && param[len(param)-1] == '"' || param[0] == '\'' && param[len(param)-1] == '\'') {
		param = param[1 : len(param)-1]
	}

	return Token{Type: TokenParam, Value: param, Command: l.lastCommand, Line: l.line, Column: start}, nil
}

// skipWhitespacePos returns the position after skipping whitespace
func (l *Lexer) skipWhitespacePos() int {
	pos := 0
	for pos < len(l.currentLine) && unicode.IsSpace(rune(l.currentLine[pos])) {
		pos++
	}
	return pos
}
