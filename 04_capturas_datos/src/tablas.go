package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "strings"
  "strconv"
  "log"
  "errors"
)

type Options struct {
  Name string
  Value int
}

type allOptions []*Options

// array de opciones
var options = allOptions {}

func showHelp() {

  fmt.Println(`
    Bienvendo al generador de tablas de multiplicar realizado con go.
    Para generar una tabla debes inclur las siguientes opciones:

    --table=value | -t=value:  Genera una tabla con el valor proporcionado por
        el parametro (value)

    --limit=value | -l=value: (opcional) Ejecuta la tabla con limite colocado (value), sino se
        coloca el valor por defecto es 10
  `)
}

func storeOption( arg string ) ( []*Options )  {

  value, error := strconv.Atoi( strings.Split( arg, "=" )[1] )

  if ( error != nil ) {
    log.SetPrefix( strings.Split( arg, "=")[0] + ": " )
    log.Fatal( errors.New( "El Valor colocado es invalido" ) )
  }

  // insert element in array
  return append( options, &Options{ Name: strings.Split( arg, "=")[0], Value: value } )
}

func checkOption() {

  if ( len( options ) > 2 ) {
    log.SetPrefix( "Arg: " )
    log.Fatal( errors.New( "Demasiados argumentos inicializados" ) )
  }

  table := false
  limit := false

  for _, option := range options {

    switch option.Name {

      case "--table":
        table = true
        break

      case "-t":
        table = true
        break

      case "--limit":
        limit = true
        break

      case "-l":
        limit = true
        break

      default:
        break
    }
  }

  if table == true && limit == true {

    saveTableLimit()

  } else if table == true && limit == false {

    saveTable()

  } else if table == false && limit == true  {

    log.SetPrefix( "Tabla: " )
    log.Fatal( errors.New( "Ingresa la opcion tabla, campo obligatiorio" ) )

  } else if table == false && limit {

    log.SetPrefix( "Opcion: " )
    log.Fatal( errors.New( "La opcion o opciones colocadas es invalido" ) )

  }

}

func saveTable() {

  var operation string
  value := options[ len( options ) - 1 ].Value

  for index := 1;  index <= 10; index++ {

    operation += strconv.Itoa( value ) + " * " + strconv.Itoa( index ) +
      " = " + strconv.Itoa( (index * value) ) + "\n"
  }

  error := ioutil.WriteFile( "tabla_" + strconv.Itoa( value ) + ".txt", []byte( operation ), 0600 )

  if error != nil {

    log.Fatal( error )
  }

  fmt.Printf( "La tabla del %d fue generada con exito\n", value )
}

func saveTableLimit() {

  var valueTable int
  var valueLimit int
  var operation string

  for _, option := range options {

    if option.Name == "-t" || option.Name == "--table" {
      
      valueTable = option.Value

    } else {
      
      valueLimit = option.Value
    }
  }

  for index := 1;  index <= valueLimit; index++ {
    
    operation += strconv.Itoa( valueTable ) + " * " + strconv.Itoa( index ) +
      " = " + strconv.Itoa( ( index * valueTable ) ) + "\n"
  }

  error := ioutil.WriteFile( "tabla_" + strconv.Itoa( valueTable ) + ".txt", []byte( operation ), 0600 )

  if error != nil {

    log.Fatal( error )
  }

  fmt.Printf( "La tabla del %d fue generada con exito\n", valueTable )
}

// ======================================================================================

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

    options = storeOption( arg )
  }

  checkOption()
}
