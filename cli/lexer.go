package cli

const a = int('a')
const z = int('z')
const A = int('A')
const Z = int('Z')
const _0 = int('0')
const _9 = int('9')

type Token struct {
	Type_   int
	Literal string
}

func CreateToken(tType int, lit string) Token {
	return Token{
		Type_:   tType,
		Literal: lit,
	}
}

type Lexer struct {
	position     int
	readPosition int
	ch           rune
	input        string
}

func NewLexer(inputText string) Lexer {
	l := Lexer{
		position:     0,
		readPosition: 0,
		ch:           0,
		input:        inputText,
	}
	l.readChar()
	return l
}

func (l *Lexer) GetNextToken() Token {
	l.skipWhiteSpace()

	var token Token
	if isLetter(l.ch) {
		token = CreateToken(0, l.readIdent())
	}

	return token
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = '\x00'

	} else {
		l.ch = rune(l.input[l.readPosition])
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdent() string {
	var position = l.position

	for isLetter(l.ch) || isNumber(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func isLetter(char rune) bool {
	charInt := int(char)
	return a <= charInt && charInt <= z ||
		A <= charInt && charInt <= Z
}
func isNumber(char rune) bool {
	charInt := int(char)
	return _0 <= charInt && charInt <= _9
}

func isWhiteSpace(char rune) bool {
	return char == ' ' || char == '\n' || char == '\r' || char == '\t'
}
