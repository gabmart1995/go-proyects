package repl

import (
	"bufio"
	"fmt"
	"os"
	"sqllite-go/internal/compiler"
	"sqllite-go/internal/vm"
)

type REPL interface {
	StartREPL()
}

type replStruct struct {
}

func New() REPL {
	return &replStruct{}
}

// inicia el proceso REPL
func (r *replStruct) StartREPL() {
	// inicia el reader para leer los datos
	// de entrada del teclado
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("DB >> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("Error reading input %v \n", err)
			continue
		}

		// toma todo lo demas menos el salto de linea
		input = input[:len(input)-1]

		// si contiene la palabra exit
		if input == "exit" {
			fmt.Println("Exiting ...")
			break
		}

		r.evaluate(input)
	}
}

func (r *replStruct) evaluate(input string) {
	comp := compiler.NewCompiler(input)
	bt := comp.Call()

	/*for _, instruction := range bt.Instructions {
		fmt.Println("Output ", instruction)
	}*/

	// generamos la maquina virtual
	machine := vm.NewVM("users", bt)
	response := machine.ExecuteByteCode()

	if response.Err != nil {
		fmt.Println(response.Err.Error())
		return
	}

	if response.MSG != nil {
		fmt.Println(*response.MSG)
	}

	if response.Cursor != nil {
		for i := 0; i < len(response.Cursor.Columns); i++ {
			fmt.Print("| ", response.Cursor.Columns[i], " ")
		}

		fmt.Println("\n-------------------------------------")

		// extraemos la informacion de las filas
		for _, record := range response.Cursor.Pages {
			for _, currentPage := range record {

				// validamos cada campo
				if currentPage.IsPrimaryKey != nil && *currentPage.IsPrimaryKey {
					fmt.Print("PK: ", *currentPage.IntValue, " ")
					continue
				}

				if currentPage.IntValue != nil {
					fmt.Print("| ", *currentPage.IntValue, " ")
					continue
				}

				if currentPage.StringValue != nil {
					fmt.Print("| ", *currentPage.StringValue, " ")
					continue
				}
			}

			fmt.Println("\n-------------------------------------")
		}
	}
}
