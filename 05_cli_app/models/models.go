package models

import (
	"encoding/json"
	"errors"
	"os"
)

type Todo struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	CompletedIn string `json:"completed_in"`
}

type Todos struct {
	Listado map[string]Todo `json:"listado"`
}

var ListTodos = Todos{
	Listado: make(map[string]Todo),
}

// constructor del modulo
func init() {

	// verifica si existe el JSON crea la estructura
	if _, err := os.Stat("todo.json"); errors.Is(err, os.ErrNotExist) {
		ListTodos.SaveJSON(ListTodos)
		return
	}

	ListTodos.GetJSON()
}

// añade la funcion a la struct Todos
func (todos *Todos) SaveJSON(newList Todos) {

	// abrimos el archivo
	// file, err := os.OpenFile("todo.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// os: CREATE: crea el archivo
	// os: WRONLY: sobreescribe todo el archivo
	// os: APPEND: coloca en la ultima posicion el contenido

	file, err := os.OpenFile("todo.json", (os.O_WRONLY | os.O_CREATE | os.O_TRUNC), 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	// transforma la data en JSON
	data, err := json.Marshal(&todos)

	if err != nil {
		panic(err)
	}

	// fmt.Println(string(data))

	if _, err := file.Write(data); err != nil {
		panic(err)
	}
}

func (todos *Todos) GetJSON() {

	// fmt.Println(todos)

	file, err := os.OpenFile("todo.json", os.O_RDONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	jsonParser := json.NewDecoder(file)

	if err := jsonParser.Decode(&todos); err != nil {
		panic(err)
	}
}
