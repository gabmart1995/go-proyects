package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"strconv" // parsea de string a num
)

/*
	CompileDeamon es un demonio que compila de forma automatica la aplicacion
	en modo de desarrollo

	tienes que setear el env de GO111MODULE=on para que funcione

	CompileDaemon -comand="./ejecutable"

*/

type Task struct {
	Id int `json:Id`
	Name string `json:Name`
	Content string `json:Content`
}

// vector principal
type allTasks []Task


var tasks = allTasks {
	{
		Id: 1,
		Name: "Tarea 1",
		Content: "Some content",
	},
}

func indexRoute( writer http.ResponseWriter, request *http.Request ) {
	
	fmt.Fprintf( writer, "Bienvendo al API" )
}

func getTasks( writer http.ResponseWriter, request *http.Request ) {

	// se colocan los datos de la cabecera
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader( 200 ) // status
	
	json.NewEncoder( writer ).Encode( tasks )
}

func getTask( writer http.ResponseWriter, request *http.Request ) {

	vars := mux.Vars( request )

	tasksId, error := strconv.Atoi( vars["id"] )

	if error != nil {
		
		http.Error( writer, "Ingrese un id valido", 400 )
		
		return
	}

	for _, task := range tasks {

		if task.Id == tasksId {

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader( 200 ) // status
			
			json.NewEncoder( writer ).Encode( task )

			return
		}
	}

	// si no encuentra el id dentro del array mostrara bad request
	http.Error( writer, "Id no encontrado", 400 )
}

func deleteTask( writer http.ResponseWriter, request *http.Request ) {

	vars := mux.Vars( request )

	tasksId, error := strconv.Atoi( vars["id"] )

	if error != nil {
		
		http.Error( writer, "Ingrese un id valido", 400 )
		
		return
	}

	for index, task := range tasks {

		if task.Id == tasksId {

			// se elimina el valor del indice
			tasks = append( tasks[ :index ], tasks[ index + 1: ]... )

			// writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader( 200 ) // status
			
			fmt.Fprintf( writer, "La tarea con el id %v ha sido removido con éxito", tasksId )

			return
		}
	}

	// si no encuentra el id dentro del array mostrara bad request
	http.Error( writer, "Id no encontrado", 400 )
}

func createTasks( writer http.ResponseWriter, request *http.Request ) {

	var newTask Task

	// tomamos los datos del body
	requestBody, error := ioutil.ReadAll( request.Body )

	if error != nil {
		http.Error( writer, "Inserte una tarea valida", 400 )
		return
	}

	json.Unmarshal( requestBody, &newTask )

	newTask.Id = len( tasks ) + 1 

	// inserta el nuevo elemento en el array
	tasks = append( tasks, newTask )

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader( 201 )

	// devolvemos el Json
	json.NewEncoder( writer ).Encode( newTask )
}


func updateTask( writer http.ResponseWriter, request *http.Request ) {

	var updatedTask Task

	vars := mux.Vars( request )

	// id
	tasksId, error := strconv.Atoi( vars["id"] )

	if error != nil {
		
		http.Error( writer, "Ingrese un id valido", 400 )
		
		return
	}

	// body
	requestBody, error := ioutil.ReadAll( request.Body )

	if error != nil {
		
		http.Error( writer, "Inserte datos validos", 400 )

		return 
	}

	json.Unmarshal( requestBody, &updatedTask )


	for index, task := range tasks {

		if task.Id == tasksId {

			// se elimina y luego se agrega el nuevo elemento
			tasks = append( tasks[:index], tasks[index + 1:]... )

			updatedTask.Id = tasksId

			tasks = append( tasks, updatedTask )

			writer.WriteHeader( 200 ) // status
			
			fmt.Fprintf( writer, "La tarea con el id %v ha sido actualizado con éxito", tasksId )

			return
		}
	}
}

// =========================================================

func main() {

	router := mux.NewRouter().StrictSlash( true )

	router.HandleFunc( "/", indexRoute )
	router.HandleFunc( "/tasks", getTasks ).Methods("GET")
	router.HandleFunc( "/tasks", createTasks ).Methods("POST")
	router.HandleFunc( "/tasks/{id}", getTask ).Methods("GET")
	router.HandleFunc( "/tasks/{id}", deleteTask ).Methods("DELETE")
	router.HandleFunc( "/tasks/{id}", updateTask ).Methods("PUT")

	// el segundo parametro es el enrutador

	fmt.Println("Servidor escuchando en el puerto 3000")

	log.Fatal( http.ListenAndServe( ":3000", router ) )
}