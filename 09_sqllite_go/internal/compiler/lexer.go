package compiler

import "strings"

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenIdentifier
	TokenKeyword
	TokenSymbol
	TokenWhiteSpace
)

type Token struct {
	Type  TokenType
	value string
}

type Lexer interface {
	NextToken() Token
}

type LexerSimple struct {
	input    string
	position int // corresponde a la posicion de los caracteres
}

func NewLexer(input string) Lexer {
	return &LexerSimple{input: input, position: 0}
}

// examina cada letra de la instruccion
func (l *LexerSimple) NextToken() Token {
	// en caso en la entrada llega vacio
	if l.position >= len(l.input) {
		return Token{value: "", Type: TokenEOF}
	}

	// extraemos el byte caracter marcado en la posicion
	char := l.input[l.position]

	// se compara con su caracter unicode
	if char == ' ' || char == '\t' || char == '\n' {
		l.consumeWhiteSpace()
		return l.NextToken()
	}

	if char == ',' || char == '*' || char == '(' || char == ')' || char == ';' || char == '\'' {
		l.position++
		return Token{Type: TokenSymbol, value: string(char)}
	}

	if isLetter(char) || isNumber(char) {
		return l.consumeIdentifier()
	}

	return Token{Type: TokenError, value: string(char)}
}

// pasa la siguiente posicion del token si es un espacio en blanco
func (l *LexerSimple) consumeWhiteSpace() {
	for l.position < len(l.input) && (l.input[l.position] == ' ' || l.input[l.position] == '\t' || l.input[l.position] == '\n') {
		l.position++
	}
}

func (l *LexerSimple) consumeIdentifier() Token {
	start := l.position // crea una copia para conservar el inicio de la palabra

	// recorremos hasta finalizar la palabra
	for l.position < len(l.input) && (isLetter(l.input[l.position]) || isNumber(l.input[l.position])) {
		l.position++
	}

	token := TokenIdentifier
	value := l.input[start:l.position]

	// verificamos si tiene una palabra clave SQL
	if strings.ToUpper(value) == "SELECT" ||
		strings.ToUpper(value) == "INSERT" ||
		strings.ToUpper(value) == "INTO" ||
		strings.ToUpper(value) == "VALUES" ||
		strings.ToUpper(value) == "FROM" ||
		strings.ToUpper(value) == "WHERE" {
		token = TokenKeyword
	}

	return Token{value: value, Type: token}
}

// identifica si el byte es un caracter alfabetico
func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

// identifica si el byte es un caracter numerico
func isNumber(char byte) bool {
	return (char >= '0' && char <= '9')
}
