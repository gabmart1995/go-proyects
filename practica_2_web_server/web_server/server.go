package main

import (
  "fmt"
  "log"
  "net/http"
  "example.com/read_file"
  "html/template"
)

func handlerExample( writer http.ResponseWriter, request *http.Request )  {
  /*
    [1:] es un sub segmento de ruta desde el primer caracter hasta el
    final
    fmt.Println( writer )
    fmt.Println( request.URL.Path[1:] )
  */

  fmt.Fprintf( writer, "Hola mundo desde el server, me gustan los %s!", request.URL.Path[1:] )
}

func viewHandler( writer http.ResponseWriter, request *http.Request )  {

  // extrae la subcadena desde el view en adelante
  title := request.URL.Path[ len("/view/"): ]
  page, _ := read_file.LoadPage( title )

  fmt.Println( page )


  renderTemplate( writer, "view", page )

  /*
    templateHTML := (`
      <h1>%s</h1>
      <div>%s</div>
    `)

    fmt.Fprintf( writer, templateHTML, page.Title, page.Body )
  */
}

func editHandler( writer http.ResponseWriter, request *http.Request )  {

  title := request.URL.Path[ len("/edit/"): ]
  page, error := read_file.LoadPage( title )

  if error != nil {
    page = &read_file.Page{ Title: title }
  }

  renderTemplate( writer, "edit", page )
}


func renderTemplate( writer http.ResponseWriter, templ string, page *read_file.Page )  {

  // retorna un *template.Template
  templateHTML, _ := template.ParseFiles( "./templates/" + templ + ".html" )
  templateHTML.Execute( writer, page )
}

// ====================================================================

func main()  {

  http.HandleFunc( "/view/", viewHandler )
  http.HandleFunc( "/edit/", editHandler )

  fmt.Println("Servidor escuchando en el puerto 8080")

  log.Fatal( http.ListenAndServe( ":8080", nil ) )
}
