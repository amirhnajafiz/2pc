package cmd

import (
	"github.com/F24-CSE535/2pc/cluster/internal/config/paxos"
	"github.com/F24-CSE535/2pc/cluster/internal/csm"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"

	"go.uber.org/zap"
)

// node is a wrapper for a single cluster entity.
type node struct {
	cfg      *paxos.Config
	logger   *zap.Logger
	database *storage.Database

	iptable map[string]string
	cluster string
	leader  string
}

func (n node) main(port int, name string) {
	// create a new CSM manager
	manager := csm.Manager{
		Cfg:     n.cfg,
		Storage: n.database,
		Memory:  memory.NewSharedMemory(n.leader, name, n.cluster, n.iptable),
	}

	// initialize CSMs with desired replica
	manager.Initialize(n.logger)

	// create a bootstrap
	b := grpc.Bootstrap{
		Logger: n.logger,
	}

	// run the grpc server
	if err := b.ListenAnsServer(port, manager.Channel, manager.DispatcherChannel, n.database); err != nil {
		n.logger.Panic("grpc server failed", zap.Error(err))
	}
}
