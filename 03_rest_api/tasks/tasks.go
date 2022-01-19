package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv" // parsea de string a num

	"github.com/gorilla/mux"
)

// types
type Task struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type allTasks []Task

type filterCallback func(task Task, index int, slice []Task) bool
type filterTask func(tasks []Task, callback filterCallback) []Task

type mapCallback func(task Task, index int, slice []Task) Task
type mapTask func(tasks []Task, callback mapCallback) []Task

type findIndexTask func(tasks []Task, callback filterCallback) int

var tasks = allTasks{
	{
		Id:      1,
		Name:    "Tarea 1",
		Content: "Some content",
	},
}
var filter filterTask
var maps mapTask
var findIndex findIndexTask

// inicializamos la libreria
func init() {

	// pollyfills slice fillter.js
	filter = func(tasks []Task, callback filterCallback) []Task {

		var result []Task

		for index, task := range tasks {

			// actualmente no se recomienda
			// se elimina el valor del indice
			// tasks = append( tasks[ :index ], tasks[ index + 1: ]... )

			if callback(task, index, tasks) {
				result = append(result, task)
			}
		}

		return result
	}

	/** pollyfills slice map.js */
	maps = func(tasks []Task, callback mapCallback) []Task {

		var result []Task

		for index, task := range tasks {
			result = append(result, callback(task, index, tasks))
		}

		return result

		/*
			// metodo no recomendado
			if task.Id == tasksId {

				// se elimina el valor del indice
				tasks = append( tasks[ :index ], tasks[ index + 1: ]... )

				return
			}
		*/
	}

	/** pollyfills slice .js */
	findIndex = func(tasks []Task, callback filterCallback) int {

		for index, task := range tasks {

			if callback(task, index, tasks) {
				return index
			}
		}

		// devolvemos -1 sino existe ninguna coicidencia
		return -1
	}
}

func GetTasks(writer http.ResponseWriter, request *http.Request) {

	// se colocan los datos de la cabecera
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200) // status

	json.NewEncoder(writer).Encode(tasks)
}

func GetTask(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	tasksId, error := strconv.Atoi(vars["id"])

	if error != nil {

		http.Error(writer, "Ingrese un id valido", 400)

		return
	}

	for _, task := range tasks {

		if task.Id == tasksId {

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(200) // status

			json.NewEncoder(writer).Encode(task)

			return
		}
	}

	// si no encuentra el id dentro del array mostrara bad request
	http.Error(writer, "Id no encontrado", 400)
}

func DeleteTask(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	tasksId, error := strconv.Atoi(vars["id"])

	if error != nil {
		http.Error(writer, "Ingrese un id valido", http.StatusNotFound)
		return
	}

	find := findIndex(tasks, func(task Task, index int, slice []Task) bool {
		return tasksId == task.Id
	})

	// si halla el elemento comienza a filtrarlo
	if find != -1 {

		tasks = filter(tasks, func(task Task, index int, slice []Task) bool {
			return task.Id != tasksId
		})

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200) // status

		fmt.Fprintf(writer, "La tarea con el id %v ha sido removido con éxito", tasksId)

	} else {
		http.Error(writer, "Id no encontrado", http.StatusNotFound)
	}
}

func CreateTasks(writer http.ResponseWriter, request *http.Request) {

	var newTask Task

	// tomamos los datos del body
	requestBody, error := ioutil.ReadAll(request.Body)

	if error != nil {
		http.Error(writer, "Inserte una tarea valida", 400)
		return
	}

	json.Unmarshal(requestBody, &newTask)

	newTask.Id = len(tasks) + 1

	// inserta el nuevo elemento en el array
	tasks = append(tasks, newTask)

	fmt.Println(tasks)
	fmt.Println(newTask)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)

	// devolvemos el Json
	json.NewEncoder(writer).Encode(newTask)
}

func UpdateTask(writer http.ResponseWriter, request *http.Request) {

	var updatedTask Task

	vars := mux.Vars(request)

	// id
	tasksId, error := strconv.Atoi(vars["id"])

	if error != nil {

		http.Error(writer, "Ingrese un id valido", 400)

		return
	}

	// body
	requestBody, error := ioutil.ReadAll(request.Body)

	if error != nil {

		http.Error(writer, "Inserte datos validos", 400)

		return
	}

	json.Unmarshal(requestBody, &updatedTask)

	// method map js
	taskMap := maps(tasks, func(task Task, index int, slice []Task) Task {

		if task.Id == tasksId {
			return updatedTask
		}

		return task
	})

	fmt.Println(taskMap)

	tasks = taskMap

	writer.WriteHeader(200) // status
	fmt.Fprintf(writer, "La tarea con el id %v ha sido actualizado con éxito", tasksId)
}
