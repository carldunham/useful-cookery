package troff

import (
	"fmt"
	"strings"
)

// TokenType represents the type of token in the TROFF lexer.
type TokenType int

const (
	// TokenUnknown represents an unknown token.
	TokenUnknown TokenType = iota
	// TokenEOF represents the end of file.
	TokenEOF
	// TokenCommand represents a TROFF command (a period at the start of a line followed by two alpha characters).
	TokenCommand
	// TokenParam represents a parameter to a command (alphanumeric string, possibly quoted).
	TokenParam
)

// String returns a string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case TokenUnknown:
		return "Unknown"
	case TokenEOF:
		return "EOF"
	case TokenCommand:
		return "Command"
	case TokenParam:
		return "Param"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}

// Token represents a lexical token in the TROFF file.
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
	Command string // For TokenParam, this is the command it belongs to
}

// String returns a string representation of the token.
func (t Token) String() string {
	if t.Type == TokenCommand {
		return fmt.Sprintf("%s(%s) at line %d, col %d", t.Type, t.Value, t.Line, t.Column)
	} else if t.Type == TokenParam {
		return fmt.Sprintf("%s(%s) for command %s at line %d, col %d", t.Type, t.Value, t.Command, t.Line, t.Column)
	}
	return fmt.Sprintf("%s at line %d, col %d", t.Type, t.Line, t.Column)
}

// IsCommand checks if the token is a specific command.
func (t Token) IsCommand(cmd string) bool {
	return t.Type == TokenCommand && strings.EqualFold(t.Value, cmd)
}
