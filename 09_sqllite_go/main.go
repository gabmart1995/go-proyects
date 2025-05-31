package main

import (
	"fmt"
	"os"
	repl "sqllite-go/cmd"
)

func main() {
	cmd := repl.New()
	cmd.StartREPL()

	fmt.Println("finished ...")
	os.Exit(0)
}
