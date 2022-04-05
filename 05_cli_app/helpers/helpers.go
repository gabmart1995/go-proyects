package helpers

import (
	"05_cli_app/models"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func CreateTodo(description string) {

	idTodo := uuid.NewString()
	todos := models.ListTodos

	// inserta la nueva propiedad
	todos.Listado[idTodo] = models.Todo{
		Id:          idTodo,
		Description: description,
		CompletedIn: "",
	}

	// salvamos los datos en el json
	todos.SaveJSON(todos)
}

func GetInput(message string) string {

	var value string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(message + " ")

	if !scanner.Scan() {
		err := errors.New("hubo un problema al capturar los datos")
		panic(err)
	}

	// le pasamos la data a description
	value = strings.Trim(scanner.Text(), " ")

	return value
}

func GetTodos() []models.Todo {

	// tranforma el objeto a slice
	var values []models.Todo

	for _, value := range models.ListTodos.Listado {
		values = append(values, value)
	}

	return values
}
