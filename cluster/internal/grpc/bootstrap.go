package grpc

import (
	"fmt"
	"net"

	"github.com/F24-CSE535/2pc/cluster/internal/grpc/services"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"

	"google.golang.org/grpc"
)

// Bootstrap is a wrapper that holds every required thing for the gRPC server starting.
type Bootstrap struct {
	Storage storage.Database
}

// ListenAnsServer creates a new gRPC instance with all required services.
func (b *Bootstrap) ListenAnsServer(port int) error {
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
		Storage: b.Storage,
	})

	// starting the server
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}

	return nil
}
