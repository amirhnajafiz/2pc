package main

import (
	"os"

	"github.com/F24-CSE535/2pc/cli/internal/handlers"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("at least one argument is needed")
	}

	// remove the first two arguments
	handlerArgs := args[2:]
	handlerArgsSize := len(handlerArgs)

	// load handlers into a list
	list := []handlers.Handler{
		&handlers.ShardHandler{},
		&handlers.ScaleHandler{},
	}

	// loop over the handlers list and execute the user command
	for _, hd := range list {
		if hd.GetName() == args[1] {
			if err := hd.Execute(handlerArgsSize, handlerArgs); err != nil {
				panic(err)
			}
		}
	}
}
