package cmd

import (
	"github.com/F24-CSE535/2pc/cluster/internal/config"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"

	"go.uber.org/zap"
)

// node is a wrapper for a single cluster entity.
type node struct {
	cfg                config.Config
	logger             *zap.Logger
	database           *storage.Database
	terminationChannel chan bool
}

func (n node) main() {
	go func() {
		// create a bootstrap
		b := grpc.Bootstrap{
			Logger:  n.logger.Named("grpc"),
			Storage: n.database,
		}

		// run the grpc server
		if err := b.ListenAnsServer(n.cfg.GRPCPort); err != nil {
			n.logger.Panic("grpc server failed", zap.Error(err))
		}
	}()

	// wait until the manager sends an event
	<-n.terminationChannel
}
