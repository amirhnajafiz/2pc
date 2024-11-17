package main

import (
	"os"

	"github.com/F24-CSE535/2pc/cluster/cmd"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		panic("at least one argument is needed")
	}

	// create a new cluster manager
	cm := cmd.Cluster{
		ClusterName: args[1],
		ConfigPath:  args[2],
	}

	// start the cluster manager
	if err := cm.Main(); err != nil {
		panic(err)
	}
}
