package grpc

import (
	"fmt"
	"net"

	"github.com/F24-CSE535/2pc/cluster/internal/grpc/services"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Bootstrap is a wrapper that holds every required thing for the gRPC server starting.
type Bootstrap struct {
	Logger *zap.Logger
}

// ListenAnsServer creates a new gRPC instance with all required services.
func (b *Bootstrap) ListenAnsServer(
	port int,
	channel chan *packets.Packet,
	dchannel chan *packets.Packet,
	st *storage.Database,
) error {
	// on the local network, listen to a port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start the listener server: %v", err)
	}

	// create a new grpc instance
	server := grpc.NewServer(
		grpc.UnaryInterceptor(b.allUnaryInterceptor),   // set an unary interceptor
		grpc.StreamInterceptor(b.allStreamInterceptor), // set a stream interceptor
	)

	// register all gRPC services
	database.RegisterDatabaseServer(server, &services.DatabaseService{
		Storage:           st,
		Channel:           channel,
		DispatcherChannel: dchannel,
	})
	paxos.RegisterPaxosServer(server, &services.PaxosService{
		Channel: channel,
	})

	// starting the server
	b.Logger.Info("grpc started", zap.Int("port", port))
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}

	return nil
}
