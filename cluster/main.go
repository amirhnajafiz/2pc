package main

import (
	"os"
	"strconv"

	"github.com/F24-CSE535/2pc/cluster/cmd"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		panic("at least four arguments are needed (./main <cluster-name> <config-path> <index>)")
	}

	// create a new cluster manager
	cm := cmd.Cluster{
		ClusterName: args[1],
		ConfigPath:  args[2],
	}
	cm.Index, _ = strconv.Atoi(args[3])

	// start the cluster manager
	if err := cm.Main(); err != nil {
		panic(err)
	}
}
