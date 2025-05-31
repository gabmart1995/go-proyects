package compiler

import "sqllite-go/internal/common/models"

type Compiler interface {
	Call() *models.ByteCode
}

type compiler struct {
	lexer     Lexer
	parser    Parser
	generator Generator
}

func NewCompiler(sqlText string) compiler {
	lexer := NewLexer(sqlText)
	parser := NewParser(lexer)
	generator := NewGenerator()

	return compiler{
		lexer:     lexer,
		parser:    parser,
		generator: generator,
	}
}

func (c *compiler) Call() *models.ByteCode {
	c.parser.ParseStatement()

	return c.generator.GenerateCode(c.parser.GetTokens())
}
