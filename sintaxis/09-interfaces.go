package main

import (
	"fmt"
)

// interfaz
type User interface {
	Permissions() int // 1-5 levels
	Name() string
}

// =================
// admin
// =================

type Admin struct {
	name string
}


func ( this Admin ) Permissions() int {
	return 5
}

func ( this Admin ) Name() string {
	return this.name
}


// =================
// editor
// =================


type Editor struct {
	name string
}


func ( this Editor ) Permissions() int {
	return 3
}

func ( this Editor ) Name() string {
	return this.name
}

func auth( user User ) string {

	permission := user.Permissions()

	if permission > 4 {
		return user.Name() + " tiene permisos de administrador"

	} else if permission == 3 {
		return user.Name() + " tiene permisos de editor"

	} else {
		return "no posee acceso"

	}
}

// =================
// main
// =================

func main() {

	admin := Admin{ name: "Gabriel" }
	editor := Editor{ name: "Alfonso" }

	// slice de usuarios
	users := []User{ admin, editor }

	for _, user := range users {
		fmt.Println( auth( user ) )
	}

	// fmt.Println( auth( editor ) )
}
