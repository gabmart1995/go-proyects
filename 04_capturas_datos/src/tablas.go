package main

import (
  "fmt"
  "os"
  // "regexp"
  "strings"
  "strconv"
  "log"
  "errors"
)

type Options struct {
  Name string
  Value int
}

type allOptions []Options

// array de opciones
var options = allOptions {}

func showHelp() {

  fmt.Println(`
    Bienvendo al generador de tablas de multiplicar realizado con go.
    Para generar una tablas puedes inclur las siguientes opciones:

    --table=value | -t=value:  Genera una tabla con el valor proporcionado por
        el parametro (value)

    --limite=value | -l=value: Ejecuta la tabla con limite colocado (value), sino se
        coloca el valor es 10
  `)
}

func main()  {

  // se preparan los logs
  log.SetFlags(0)

  // se obtienen los argumentos usado en el ejecutable
  args := os.Args

  if ( len( args ) == 1 ) {
    showHelp()
    return
  }

  // elimina el primer vector que es el nombre del programa
  args = append( args[1:] )

  for _, arg := range args {

    if ( arg == "--help" || arg == "-h" ) {
      showHelp()
      return
    }

    value, error := strconv.Atoi( strings.Split( arg, "=")[1] )

    if ( error != nil ) {
      log.SetPrefix( strings.Split( arg, "=")[0] + ": " )
      log.Fatal( errors.New( "El Valor colocado es invalido" ) )
    }

    newOption := &Options{ Name: strings.Split( arg, "=")[0], Value: value }

    // insert element
    // options = append( options, newOption )

    fmt.Println( options )
    fmt.Println( newOption )

  }
}
