package client

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

// Accept gets a target and request to call Accept RPC on the target.
func (c *Client) Accept(address string, msg *paxos.AcceptMsg) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call accept RPC
	if _, err := paxos.NewPaxosClient(conn).Accept(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// Accepted gets a target to call Accepted RPC on the target.
func (c *Client) Accepted(address string, acceptedNum *paxos.BallotNumber, acceptedVal *paxos.AcceptMsg) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call accepted RPC
	if _, err := paxos.NewPaxosClient(conn).Accepted(context.Background(), &paxos.AcceptedMsg{
		AcceptedNumber: acceptedNum,
		AcceptedValue:  acceptedVal,
	}); err != nil {
		return err
	}

	return nil
}

// Commit gets a target to call Commit RPC on the target.
func (c *Client) Commit(address string, acceptedNum *paxos.BallotNumber, acceptedVal *paxos.AcceptMsg) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call commit RPC
	if _, err := paxos.NewPaxosClient(conn).Commit(context.Background(), &paxos.CommitMsg{
		AcceptedNumber: acceptedNum,
		AcceptedValue:  acceptedVal,
	}); err != nil {
		return err
	}

	return nil
}

// Ping gets a target and sends a ping message to the target.
func (c *Client) Ping(address string, lcm *paxos.BallotNumber) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call ping RPC
	if _, err := paxos.NewPaxosClient(conn).Ping(context.Background(), &paxos.PingMsg{
		LastCommitted: lcm,
		NodeId:        c.nodeId,
	}); err != nil {
		return err
	}

	return nil
}

// Pong gets a target and sends a pong message to the target.
func (c *Client) Pong(address string, lcm *paxos.BallotNumber) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call pong RPC
	if _, err := paxos.NewPaxosClient(conn).Pong(context.Background(), &paxos.PongMsg{
		LastCommitted: lcm,
		NodeId:        c.nodeId,
	}); err != nil {
		return err
	}

	return nil
}

// Sync gets a target and sends a sync message to the target.
func (c *Client) Sync(address string, lcm *paxos.BallotNumber, items []*paxos.SyncItem) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call sync RPC
	if _, err := paxos.NewPaxosClient(conn).Sync(context.Background(), &paxos.SyncMsg{
		LastCommitted: lcm,
		Items:         items,
	}); err != nil {
		return err
	}

	return nil
}
