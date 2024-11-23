package client

import (
	"context"
	"fmt"

	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
)

// Request accepts a transaction parameters for an inter-shard transaction.
func (c *Client) Request(address string, msg *database.RequestMsg) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Request RPC
	if _, err = database.NewDatabaseClient(conn).Request(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// Request accepts a transaction parameters for a cross-shard transaction.
func (c *Client) Prepare(address string, msg *database.PrepareMsg) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Prepare RPC
	if _, err = database.NewDatabaseClient(conn).Prepare(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// Reply calls the reply RPC on the client.
func (c *Client) Reply(address, text string, sessionId int) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Reply RPC
	if _, err := database.NewDatabaseClient(conn).Reply(context.Background(), &database.ReplyMsg{
		SessionId: int64(sessionId),
		Text:      text,
	}); err != nil {
		return fmt.Errorf("failed to call reply RPC: %v", err)
	}

	return nil
}

// Ack calls the ack RPC on the client.
func (c *Client) Ack(address string, sessionId int, isAborted bool) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Ack RPC
	if _, err := database.NewDatabaseClient(conn).Ack(context.Background(), &database.AckMsg{
		SessionId: int64(sessionId),
		IsAborted: isAborted,
		NodeId:    c.nodeId,
	}); err != nil {
		return fmt.Errorf("failed to call ack RPC: %v", err)
	}

	return nil
}

// DatabaseCommit accepts a target and sessionId to send a commit message.
func (c *Client) DatabaseCommit(address string, sessionId int) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call Commit RPC
	if _, err = database.NewDatabaseClient(conn).Commit(context.Background(), &database.CommitMsg{
		SessionId:     int64(sessionId), // set the session id
		ReturnAddress: c.nodeId,         // set the return address
	}); err != nil {
		return err
	}

	return nil
}

// DatabaseAbort accepts a target and sessionId to send an abort message.
func (c *Client) DatabaseAbort(address string, sessionId int) error {
	// base connection
	conn, err := c.connect(address)
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
