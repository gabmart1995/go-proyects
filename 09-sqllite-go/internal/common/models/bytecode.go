package models

import (
	"encoding/json"
	"log"
)

type ByteCodeOperationType int

const (
	ByteCodeOperationTypeINSERT = iota
	ByteCodeOperationTypeSELECT
	ByteCodeOperationTypeTableName
	ByteCodeOperationTypeIdentifier
	ByteCodeOperationTypeValue
	ByteCodeOperationTypeCount
)

type ByteCodeValue struct {
	Type        ByteCodeOperationType
	Identifier  *string
	IntValue    *int
	StringValue *string
	Count       int
}

type ByteCode struct {
	Instructions []ByteCodeValue
}

// transforma la instancia de ByteCodeValue a un json string
func (b ByteCodeValue) String() string {
	s, err := json.Marshal(b)

	if err != nil {
		log.Panic(err)
	}

	return string(s)
}
