package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	grpc "github.com/F24-CSE535/2pc/client/internal/grpc/dialer"
	"github.com/F24-CSE535/2pc/client/internal/manager"
	"github.com/F24-CSE535/2pc/client/internal/storage"
	"github.com/F24-CSE535/2pc/client/internal/utils"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		panic("at least fours arguments are needed (./main <path-to-iptable> <mongodb-uri> <database>)")
	}

	// load iptable
	nodes, err := utils.IPTableParseFile(args[1])
	if err != nil {
		panic(err)
	}

	// open database connection
	db, err := storage.NewDatabase(args[2], args[3])
	if err != nil {
		panic(err)
	}

	// create a new commands instance
	mg := manager.NewManager(&grpc.Dialer{Nodes: nodes}, db)

	// in a for loop, read user commands
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		input, _ := reader.ReadString('\n') // read input until newline
		input = strings.TrimSpace(input)

		// no input
		if len(input) == 0 {
			continue
		}

		// split into parts
		parts := strings.Split(input, " ")

		// create args for the client handlers
		cargs := parts[1:]
		cargsc := len(cargs)

		// switch on the first input as the command
		switch parts[0] {
		case "transaction":
			fmt.Println(mg.Transaction(cargsc, cargs))
		case "rt":
			fmt.Println(mg.RoundTrip(cargsc, cargs))
		case "printbalance":
			fmt.Println(mg.PrintBalance(cargsc, cargs))
		case "performance":
			fmt.Println(mg.Performance())
		case "exit":
			return
		default:
			fmt.Printf("command `%s` not found.\n", parts[0])
		}
	}
}
