package filesys

import (
	"fmt"
)

func CreateProyectWeb( name string ) {
	fmt.Println("desde web")
	fmt.Println( name )
}

func CreateProyectConsole( name string ) {
	fmt.Println("desde consola")
	fmt.Println( name )
} 