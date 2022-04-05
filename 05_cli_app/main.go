package main

import (
	"05_cli_app/helpers"
	"fmt"
	"os"
	"os/exec"
)

// colores de consola
const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)

func menu() {

	title := (`
=======================================
	Seleccione una opción
=======================================
`)

	menu := (`
1.- Crear una tarea
2.- Listar tareas
3.- Listar tareas completadas
4.- Listar tareas pendientes
5.- Completar tarea(s)
6.- Borrar tarea
0.- Salir

`)

	fmt.Printf(ColorGreen + title)
	fmt.Printf(ColorReset + menu)
}

func main() {

	var option int

	menu()

	fmt.Printf("Seleccione una opción: ")
	fmt.Scan(&option)

	for option != 0 {

		switch option {
		case 1:

			input := helpers.GetInput("Ingrese la descripción de la tarea:")
			helpers.CreateTodo(input)

			break

		case 2:

			todos := helpers.GetTodos()
			fmt.Println(todos)

			break

		case 3:
			fmt.Printf("option 3")
			break

		case 4:
			fmt.Printf("option 4")
			break

		case 5:
			fmt.Printf("option 5")
			break

		case 6:
			fmt.Printf("option 6")

		default:
			fmt.Printf("La opción seleccionada no es válida")
			break
		}

		fmt.Printf("\n\nPresione " + ColorGreen + "ENTER" + ColorReset + " para continuar...\n")
		fmt.Scanln()

		// limpiamos la consola
		command := exec.Command("clear")
		command.Stdout = os.Stdout
		command.Run()

		// corremos nuevamente el menu y cambiamos el valor
		menu()
		fmt.Printf("Seleccione una opción: ")
		fmt.Scan(&option)
	}
}
