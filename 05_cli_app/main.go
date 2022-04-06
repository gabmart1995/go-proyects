package main

import (
	"05_cli_app/helpers"
	"05_cli_app/models"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// colores de consola
const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
	ColorCyan   = "\u001B[36m"
	ColorWhite  = "\u001B[37m"
)

func menu() {

	var title string
	var menu string

	title += ColorGreen + "=======================================\n"
	title += ColorReset + "\tSeleccione una opción\n"
	title += ColorGreen + "=======================================\n"

	menu += ColorCyan + "1.-" + ColorReset + " Crear una tarea\n"
	menu += ColorCyan + "2.-" + ColorReset + " Listar tareas\n"
	menu += ColorCyan + "3.-" + ColorReset + " Listar tareas completadas\n"
	menu += ColorCyan + "4.-" + ColorReset + " Listar tareas pendientes\n"
	menu += ColorCyan + "5.-" + ColorReset + " Completar tarea(s)\n"
	menu += ColorCyan + "6.-" + ColorReset + " Borrar tarea\n"
	menu += ColorCyan + "0.-" + ColorReset + " Salir\n\n"

	fmt.Print(title)
	fmt.Print(menu)
}

func clearConsole() {

	// limpiamos la consola
	command := exec.Command("clear")
	command.Stdout = os.Stdout
	command.Run()
}

func formatTodo(todos []models.Todo, subTitle string) {

	var title string

	title += ColorGreen + "============================================\n"
	title += ColorReset + "\tListado de tareas " + subTitle + "\n"
	title += ColorGreen + "============================================\n"

	fmt.Print(title)

	if len(todos) == 0 {
		fmt.Println(ColorRed + "\nNo existen registros de tareas" + ColorReset)
	}

	for index, todo := range todos {

		var todoFormat string

		if len(todo.CompletedIn) > 0 {
			todoFormat = ColorGreen + todo.CompletedIn

		} else {
			todoFormat = ColorRed + "Pendiente"
		}

		fmt.Printf("%s%d.- %s%s :: %s%s\n",
			ColorCyan,
			(index + 1),
			ColorReset,
			todo.Description,
			todoFormat,
			ColorReset,
		)
	}
}

func updateOrDeleteTodo(todos []models.Todo, delete bool, message string) error {

	find := func(todos []models.Todo, index int) (models.Todo, error) {

		for idx, todo := range todos {
			if idx == index {
				return todo, nil
			}
		}

		return models.Todo{}, errors.New("no se encontro el elemento")
	}

	input := helpers.GetInput(message)
	index, err := strconv.Atoi(input)

	if err != nil {
		return errors.New("el valor ingresado debe ser un número")
	}

	if index != -1 {

		todo, errFind := find(todos, (index - 1))

		// en caso de no encontrar el valor en el array
		if errFind != nil {
			return errFind
		}

		// flag delete
		if delete {

			helpers.DeleteTodo(todo.Id)
			fmt.Println("Tarea borrada con éxito")

		} else {

			// actualizamos la tarea
			helpers.UpdateTodo(todo)
			fmt.Println("Tarea actualizada con éxito")
		}
	}

	return nil
}

func main() {

	var option int

	menu()

	fmt.Printf("Seleccione una opción númerica: ")
	fmt.Scan(&option)

	for option != 0 {

		switch option {
		case 1:
			input := helpers.GetInput("Ingrese la descripción de la tarea o -1 para cancelar:")

			if input != "-1" {
				helpers.CreateTodo(input)
			}

			break

		// listar tareas
		case 2:
			todos := helpers.GetTodos()
			clearConsole()
			formatTodo(todos, "")

			break

		// listar tareas completadas
		case 3:
			todos := helpers.GetTodosCompleted()
			clearConsole()
			formatTodo(todos, "completadas")
			break

		// listar tareas pendientes
		case 4:
			todos := helpers.GetTodosPendient()
			clearConsole()
			formatTodo(todos, "pendientes")
			break

		// completar tareas
		case 5:

			todos := helpers.GetTodos()

			clearConsole()

			formatTodo(todos, "para completar")

			err := updateOrDeleteTodo(
				todos,
				false,
				"\nIngrese el numero de la tarea para completarla o -1 para cancelar:",
			)

			if err != nil {
				fmt.Println(err.Error())
			}

			break

		// borrar tarea
		case 6:

			todos := helpers.GetTodos()

			clearConsole()

			formatTodo(todos, "para borrar")

			err := updateOrDeleteTodo(
				todos,
				true,
				"\nIngrese el numero de la tarea para borrarla o -1 para cancelar:",
			)

			if err != nil {
				fmt.Println(err.Error())
			}

			break

		default:
			fmt.Printf("La opción seleccionada no es válida")
			break
		}

		fmt.Printf("\n\nPresione " + ColorGreen + "ENTER" + ColorReset + " para continuar...\n")
		fmt.Scanln()

		clearConsole()

		// corremos nuevamente el menu y cambiamos el valor
		menu()

		fmt.Printf("Seleccione una opción númerica: ")
		fmt.Scan(&option)
	}
}
