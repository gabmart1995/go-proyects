package compiler

import (
	"sqllite-go/internal/common/models"
	"strconv"
	"strings"
)

type Generator interface {
	GenerateCode(tokens []Token) *models.ByteCode
}

type generator struct {
}

func NewGenerator() Generator {
	return &generator{}
}

func (g *generator) GenerateCode(tokens []Token) *models.ByteCode {
	var bt *models.ByteCode

	// evaluamos los token keyword para pasarlos al bytecode
	if tokens[0].Type == TokenKeyword && strings.ToUpper(tokens[0].value) == "INSERT" {
		bt = g.generateCodeInsert(tokens)
	}

	if tokens[0].Type == TokenKeyword && strings.ToUpper(tokens[0].value) == "SELECT" {
		bt = g.generateCodeSelect(tokens)
	}

	return bt
}

func (g *generator) generateCodeSelect(tokens []Token) *models.ByteCode {
	bt := &models.ByteCode{}

	/** Se preparan las instrucciones al bytecode */
	instructions := []models.ByteCodeValue{
		{
			Type:       models.ByteCodeOperationTypeINSERT,
			Identifier: stringToStringPtr("SELECT"),
		},
	}

	varNames := []models.ByteCodeValue{}
	i := 1

	// recorremos los tokens para generar los bytecode
	for i < len(tokens) {
		if tokens[i].Type == TokenKeyword && strings.ToUpper(tokens[i].value) == "FROM" {
			break
		}

		// si es un token identificador se lo pasamos a varNames
		if tokens[i].Type == TokenIdentifier {
			v := models.ByteCodeValue{
				Type:       models.ByteCodeOperationTypeIdentifier,
				Identifier: stringToStringPtr(tokens[i].value),
			}

			varNames = append(varNames, v)
		}

		i++
	}

	// se aumenta el valor para pasar al siguiente token de la lista
	i++

	tblName := models.ByteCodeValue{
		Type:       models.ByteCodeOperationTypeTableName,
		Identifier: stringToStringPtr(tokens[i].value),
	}

	countVarNames := models.ByteCodeValue{
		Type:  models.ByteCodeOperationTypeCount,
		Count: len(varNames),
	}

	// insertamos las instrucciones
	instructions = append(instructions, tblName)
	instructions = append(instructions, countVarNames)
	instructions = append(instructions, varNames...)

	bt.Instructions = instructions

	return bt
}

func (g *generator) generateCodeInsert(tokens []Token) *models.ByteCode {
	bt := &models.ByteCode{}

	instructions := []models.ByteCodeValue{
		{
			Type:       models.ByteCodeOperationTypeINSERT,
			Identifier: stringToStringPtr("INSERT"),
		},
	}

	tblName := models.ByteCodeValue{
		Type:       models.ByteCodeOperationTypeTableName,
		Identifier: stringToStringPtr(tokens[2].value),
	}

	varNames := []models.ByteCodeValue{}
	i := 3

	for i < len(tokens) {
		if tokens[i].Type == TokenKeyword && strings.ToUpper(tokens[i].value) == "VALUES" {
			break
		}

		if tokens[i].Type == TokenIdentifier {
			varName := models.ByteCodeValue{
				Type:       models.ByteCodeOperationTypeIdentifier,
				Identifier: stringToStringPtr(tokens[i].value),
			}

			varNames = append(varNames, varName)
		}

		i++
	}

	// se aumenta el valor para pasar al siguiente token de la lista
	i++

	countVarNames := models.ByteCodeValue{
		Type:  models.ByteCodeOperationTypeCount,
		Count: len(varNames),
	}

	// procedemos a recorrer los valores
	varValues := []models.ByteCodeValue{}

	for i < len(tokens) {
		if tokens[i].Type == TokenSymbol && strings.ToUpper(tokens[i].value) == ";" {
			break
		}

		if tokens[i].Type == TokenIdentifier {
			varValue := models.ByteCodeValue{
				Type: models.ByteCodeOperationTypeIdentifier,
			}

			if isInteger(tokens[i].value) {
				varValue.IntValue = stringToIntPtr(tokens[i].value)

			} else {
				varValue.StringValue = stringToStringPtr(tokens[i].value)
			}

			varValues = append(varValues, varValue)
		}

		i++
	}

	countVarValues := models.ByteCodeValue{
		Type:  models.ByteCodeOperationTypeCount,
		Count: len(varValues),
	}

	// insertamos las instrucciones
	instructions = append(instructions, tblName)
	instructions = append(instructions, countVarNames)
	instructions = append(instructions, varNames...)
	instructions = append(instructions, countVarValues)
	instructions = append(instructions, varValues...)

	bt.Instructions = instructions

	return bt
}

// helper functions
func stringToStringPtr(input string) *string {
	return &input
}

func stringToIntPtr(input string) *int {
	num, _ := strconv.Atoi(input)

	return &num
}

func isInteger(input string) bool {
	_, err := strconv.Atoi(input)

	return err == nil
}
