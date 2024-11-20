package grpc

import (
	"context"
	"fmt"

	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dialer is a module for making RPC calls from client to clusters.
type Dialer struct {
	Nodes map[string]string
}

// connect should be called in the beginning of each method to establish a connection.
func (d *Dialer) connect(target string) (*grpc.ClientConn, error) {
	address := d.Nodes[d.Nodes[target]]
	fmt.Println(address)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to %s: %v", address, err)
	}

	return conn, nil
}

// Request accepts a transaction parameters for an inter-shard transaction.
func (d *Dialer) Request(target, sender, receiver string, amount, sessionId int) error {
	// base connection
	conn, err := d.connect(target)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Request RPC
	if _, err = database.NewDatabaseClient(conn).Request(context.Background(), &database.RequestMsg{
		Transaction: &database.TransactionMsg{ // initialize a new transaction
			Sender:    sender,
			Receiver:  receiver,
			Amount:    int64(amount),
			SessionId: int64(sessionId),
		},
		ReturnAddress: d.Nodes["client"], // set the return address
	}); err != nil {
		return err
	}

	return nil
}

// Request accepts a transaction parameters for a cross-shard transaction.
func (d *Dialer) Prepare(target, client, sender, receiver string, amount, sessionId int) error {
	// base connection
	conn, err := d.connect(target)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Prepare RPC
	if _, err = database.NewDatabaseClient(conn).Prepare(context.Background(), &database.PrepareMsg{
		Transaction: &database.TransactionMsg{ // initialize a new transaction
			Sender:    sender,
			Receiver:  receiver,
			Amount:    int64(amount),
			SessionId: int64(sessionId),
		},
		Client:        client,            // set the client for cluster usage
		ReturnAddress: d.Nodes["client"], // set the return address
	}); err != nil {
		return err
	}

	return nil
}

// Commit accepts a target and sessionId to send a commit message.
func (d *Dialer) Commit(target string, sessionId int) error {
	// base connection
	conn, err := d.connect(target)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Commit RPC
	if _, err = database.NewDatabaseClient(conn).Commit(context.Background(), &database.CommitMsg{
		SessionId:     int64(sessionId),  // set the session id
		ReturnAddress: d.Nodes["client"], // set the return address
	}); err != nil {
		return err
	}

	return nil
}

// Abort accepts a target and sessionId to send an abort message.
func (d *Dialer) Abort(target string, sessionId int) error {
	// base connection
	conn, err := d.connect(target)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Abort RPC
	if _, err = database.NewDatabaseClient(conn).Abort(context.Background(), &database.AbortMsg{
		SessionId: int64(sessionId), // set the session id
	}); err != nil {
		return err
	}

	return nil
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
