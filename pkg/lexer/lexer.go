package lexer

import (
	"unicode"
	"unicode/utf8"
	"github.com/26thavenue/json_parser/pkg/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	columnStart  int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, columnStart: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		var size int
		l.ch, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.position = l.readPosition
		l.readPosition += size
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	start := l.position - l.columnStart

	var tok token.Token
	tok.Line = l.line
	tok.Start = start

	switch l.ch {
	case '{':
		tok = l.newToken(token.LeftBrace, "{")
	case '}':
		tok = l.newToken(token.RightBrace, "}")
	case '[':
		tok = l.newToken(token.LeftBracket, "[")
	case ']':
		tok = l.newToken(token.RightBracket, "]")
	case ',':
		tok = l.newToken(token.Comma, ",")
	case ':':
		tok = l.newToken(token.Colon, ":")
	case '"':
		tok.Type = token.String
		tok.Literal = l.readString()
	case '\'':
		tok = l.newToken(token.QUOTE, "'")
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isDigit(l.ch) {
			tok.Type = token.Number
			tok.Literal = l.readNumber()
			tok.End = l.position - l.columnStart
			return tok
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			tok.End = l.position - l.columnStart
			return tok
		} else {
			tok = l.newToken(token.Illegal, string(l.ch))
		}
	}

	l.readChar()
	tok.End = l.position - l.columnStart
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.columnStart = l.readPosition
		}
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: l.line, Start: l.position - l.columnStart, End: l.position - l.columnStart + len(literal)}
}

func (l *Lexer) lookupIdent(ident string) token.TokenType {
	switch ident {
	case "true":
		return token.True
	case "false":
		return token.False
	case "null":
		return token.Null
	default:
		return token.Illegal
	}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}