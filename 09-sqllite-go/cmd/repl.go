package repl

import (
	"bufio"
	"fmt"
	"os"
	"sqllite-go/internal/compiler"
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

	for _, instruction := range bt.Instructions {
		fmt.Println("Output ", instruction)
	}
}
