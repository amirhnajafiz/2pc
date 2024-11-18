package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/F24-CSE535/2pc/client/internal/grpc"
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

	// create a new dialer
	dialer := grpc.Dialer{Nodes: nodes}

	// open database connection
	db, err := storage.NewDatabase(args[2], args[3])
	if err != nil {
		panic(err)
	}

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

		switch parts[0] {
		case "printbalance":
			// check the number of arguments
			if cargsc < 1 {
				fmt.Println("not enough arguments")
				continue
			}

			// get the shard
			cluster, err := db.GetClientShard(cargs[0])
			if err != nil {
				fmt.Printf("db failed: %v\n", err)
				continue
			}

			// make RPC call
			if balance, err := dialer.PrintBalance(cluster, cargs[0]); err != nil {
				fmt.Printf("failed: %v\n", err)
			} else {
				fmt.Printf("%s : %d\n", cargs[0], balance)
			}
		}
	}
}
