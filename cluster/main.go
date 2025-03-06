package main

import (
	"os"

	"github.com/F24-CSE535/2pc/cluster/cmd"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		panic("at least four arguments are needed (./main <config-path> <iptable-path>)")
	}

	// create a new cluster manager
	cm := cmd.Cluster{
		ConfigPath:  args[1],
		IPTablePath: args[2],
	}

	// start the cluster manager
	if err := cm.Main(); err != nil {
		panic(err)
	}
}
