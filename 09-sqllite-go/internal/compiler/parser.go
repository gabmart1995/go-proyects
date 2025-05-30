package compiler

type Parser interface {
	ParseStatement()
	GetTokens() []Token
}

type ParserSimple struct {
	lexer  Lexer
	Tokens []Token
}

func NewParser(l Lexer) Parser {
	return &ParserSimple{lexer: l, Tokens: []Token{}}
}

func (p *ParserSimple) ParseStatement() {
	for {
		token := p.lexer.NextToken()
		p.Tokens = append(p.Tokens, token)

		// si el ultimo token es el fin de la instruccion
		// sale del parser
		if token.Type == TokenEOF {
			break
		}
	}
}

func (p *ParserSimple) GetTokens() []Token {
	return p.Tokens
}
