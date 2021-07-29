package main

import (
	"os"
	"fmt"
	"log"
	"errors"
	"bit-lab.com/filesys"
)

type App struct {
	Option string
	Name string
	Basic_web bool
	Basic_console bool
}

func showHelp() {

	fmt.Println(`
    Bienvendo al generador de proyectos web realizado con go.
    Para generar un estructura debes incluir las siguientes opciones:

    new "name" Genera un nuevo proyecto con el valor pasado por el nombre ( requerido )

    flags:
    	--basic-web | --web: crea una estructura para proyecto basico web
    	--basic-console | --console: crea una estructura basica para crear proyectos de consola 
  `)
}

func main() {
	args := os.Args
  app := App{}

	if ( len( args ) == 1 ) { 
    
    showHelp()
   
    return
  }

  // elimina el primer vector que es el nombre del programa
  args = append( args[1:] )

  if len( args ) > 3 {
  	log.Fatal( errors.New( "Demasiados arguementos" ) )

  } else if len( args ) < 3 {
  	log.Fatal( errors.New( "Faltan arguementos" ) )
  
  }

  // option
  if args[0] == "new" {
		app.Option = args[0]
  
  } else {
		log.Fatal( errors.New( "Opcion incorrecta" ) )

  }

  // name
  if args[1] != "" {
  	app.Name = args[1]

  } else {
  	log.Fatal( errors.New( "Debes incluir el nombre del proyecto" ) )
  
  }

  // type proyect
  if args[2] == "--basic-web" || args[2] == "--web" {
  	app.Basic_web = true
  	app.Basic_console = false

  
  } else if args[2] == "--basic-console" || args[2] == "--console" {
  	app.Basic_web = false
  	app.Basic_console = true

  } else {
  	log.Fatal( errors.New( "Tipo de proyecto incorrecto" ) )

  }


  // ejecuta la instruccion que crea el proyect

  if app.Basic_console {
  	filesys.CreateProyectConsole( app.Name )

  } else {
  	filesys.CreateProyectWeb( app.Name )
  
  }
}