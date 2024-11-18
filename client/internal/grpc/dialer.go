package grpc

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/F24-CSE535/2pc/client/internal/rpc/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dialer is a module for making RPC calls from client to clusters.
type Dialer struct {
	Nodes map[string][]string
}

// connect should be called in the beginning of each method to establish a connection.
func (d *Dialer) connect(target string) (*grpc.ClientConn, error) {
	randomIndex := rand.Intn(len(d.Nodes[target]))
	address := d.Nodes[target][randomIndex]

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to %s: %v", address, err)
	}

	return conn, nil
}

// PrintBalance accepts a target and client to return the client balance.
func (d *Dialer) PrintBalance(target string, client string) (int, error) {
	// base connection
	conn, err := d.connect(target)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// call PrintBalance RPC
	resp, err := database.NewDatabaseClient(conn).PrintBalance(context.Background(), &database.PrintBalanceMsg{
		Client: client,
	})
	if err != nil {
		return 0, err
	}

	return int(resp.GetBalance()), nil
}
