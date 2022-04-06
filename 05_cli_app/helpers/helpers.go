package helpers

import (
	"05_cli_app/models"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateTodo(description string) {

	// currentTime := time.Now()
	// currentTime.Format("2006-01-02 15:04:05")

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

	fmt.Println("Tarea creada con éxito")
}

func UpdateTodo(todo models.Todo) {

	currentTime := time.Now()
	todos := models.ListTodos

	if len(todo.CompletedIn) == 0 {
		todo.CompletedIn = currentTime.Format("2006-01-02 15:04:05")

	} else {
		todo.CompletedIn = ""

	}

	todos.Listado[todo.Id] = todo

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

	// lee nuevamente los valores en el JSON
	models.ListTodos.GetJSON()

	// tranforma el objeto a slice
	var listTodos []models.Todo

	for _, todo := range models.ListTodos.Listado {
		listTodos = append(listTodos, todo)
	}

	return listTodos
}

func GetTodosCompleted() []models.Todo {

	// lee nuevamente los valores en el JSON
	models.ListTodos.GetJSON()

	// tranforma el objeto a slice
	var listTodos []models.Todo

	for _, todo := range models.ListTodos.Listado {

		if len(todo.CompletedIn) > 0 {
			listTodos = append(listTodos, todo)
		}
	}

	return listTodos
}

func GetTodosPendient() []models.Todo {

	// lee nuevamente los valores en el JSON
	models.ListTodos.GetJSON()

	// tranforma el objeto a slice
	var listTodos []models.Todo

	for _, todo := range models.ListTodos.Listado {

		if len(todo.CompletedIn) == 0 {
			listTodos = append(listTodos, todo)
		}
	}

	return listTodos
}

func DeleteTodo(id string) {

	// se extrae la propiedad del diccionario
	// el ok indica si la propiedad existe en el diccionario
	_, ok := models.ListTodos.Listado[id]

	if ok {

		delete(models.ListTodos.Listado, id)

		models.ListTodos.SaveJSON(models.ListTodos)

		return
	}

	fmt.Println("No se encontro id de la tarea")
}
