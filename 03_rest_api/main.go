package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/rest_api/tasks"
	"github.com/gorilla/mux"
)

/*
	CompileDeamon es un demonio que compila de forma automatica la aplicacion
	en modo de desarrollo
	tienes que setear el env de GO111MODULE=on para que funcione
	CompileDaemon -comand="./ejecutable"
*/

func index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Bienvendo al API")
}

// =========================================================

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)
	router.HandleFunc("/tasks", tasks.GetTasks).Methods("GET")
	router.HandleFunc("/tasks", tasks.CreateTasks).Methods("POST")
	router.HandleFunc("/tasks/{id}", tasks.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", tasks.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", tasks.UpdateTask).Methods("PUT")

	// el segundo parametro es el enrutador

	fmt.Println("Servidor escuchando en el puerto 3000")

	log.Fatal(http.ListenAndServe(":3000", router))
}
