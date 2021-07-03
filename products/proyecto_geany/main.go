package main

import (
	"net/http"
	"log"
	"html/template"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	// "fmt"
)

type Employee struct {
	Id int
	Name string
	Email string
}

// asigna el direcotrio de los templates
var templ = template.Must( template.ParseGlob("templates/*") )

// conexion a BD
func connectBD() ( conexion *sql.DB ) {

	Driver := "mysql"
	Usuario := "test"
	Password := "123456"
	NameDB := "prueba"

	connection, err := sql.Open( Driver, Usuario + ":" + Password +
		"@tcp(127.0.0.1)/" + NameDB )

	if err != nil {
		panic( err.Error() )
	}

	return connection
}


func index( w http.ResponseWriter, r *http.Request ) {

	// base: es el nombre definido en la plantilla si no se
	// coloca por defecto tiene que usar la extension del archivo

	connection := connectBD()

	results, err := connection.Query("SELECT * FROM empleados")

	if err != nil {
		panic( err.Error() )
	}

	employees := []Employee{}

	for results.Next() {

		var name string
		var email string
		var id int

		/* Se valida los datos de la DB */

		err := results.Scan( &id, &name, &email )

		if err != nil {
			panic( err.Error() )
		}

		employee := Employee{ Id: id, Name: name, Email: email }

		employees = append( employees, employee )
	}

	// fmt.Println( employees )

	templ.ExecuteTemplate( w, "base", employees )
}

func create( w http.ResponseWriter, r *http.Request ) {
	templ.ExecuteTemplate( w, "create", nil )
}

func insert( w http.ResponseWriter, r *http.Request ) {

	if r.Method == "POST" {

		var name string
		var correo string

		name = r.FormValue("nombre")
		correo = r.FormValue("correo")

		connection := connectBD()

		sentence, err := connection.Prepare("INSERT INTO empleados (nombre, correo) VALUES (?, ?)")

		if err != nil {
			panic( err.Error() )
		}

		// ejecuta la sentencia sql
		sentence.Exec( name, correo )

		// redirecciona a la vista
		http.Redirect( w, r, "/", 301 )
	}
}

func delete( w http.ResponseWriter, r *http.Request )  {

	// permite obtener el id del Query
	id := r.URL.Query().Get("id")

	connection := connectBD()

	sentence, err := connection.Prepare("DELETE FROM empleados WHERE id = ?")

	if err != nil {
		panic( err.Error() )
	}

	// ejecuta la sentencia sql
	sentence.Exec( id )

	// redirecciona a la vista
	http.Redirect( w, r, "/", 301 )
}

func edit( w http.ResponseWriter, r *http.Request )  {

	// permite obtener el id del Query
	id := r.URL.Query().Get("id")

	connection := connectBD()

	result, err := connection.Query( "SELECT * FROM empleados WHERE id = ?", id )
	employee := Employee{}

	if err != nil {
		panic( err.Error() )
	}

	for result.Next() {

		var id int
		var name string
		var email string

		err = result.Scan( &id, &name, &email )

		if err != nil {
			panic( err.Error() )
		}

		employee.Id = id
		employee.Name = name
		employee.Email = email
	}

	templ.ExecuteTemplate( w, "edit", employee )
}

func update( w http.ResponseWriter, r *http.Request )  {

	if r.Method == "POST" {

		name := r.FormValue("nombre") 
		id := r.FormValue("id");
		email := r.FormValue("correo")

		connection := connectBD()

		sentence, err := connection.Prepare("UPDATE empleados SET nombre = ?, correo = ? WHERE id = ?")

		if err != nil {
			panic( err.Error() )
		}

		sentence.Exec( name, email, id )

		// redirecciona
		http.Redirect( w, r, "/", 301 )
	}
}

// ========================================================

func main() {

	http.HandleFunc( "/", index )
	http.HandleFunc( "/create", create )
	http.HandleFunc( "/insert", insert )
	http.HandleFunc( "/delete", delete )
	http.HandleFunc( "/edit", edit )
	http.HandleFunc( "/update", update )

	log.Println("Servidor corriendo en el puerto 8000")

	log.Fatal( http.ListenAndServe( ":8000", nil ) )
}
