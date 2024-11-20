package client

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

// Accept gets a target and request to call Accept RPC on the target.
func (c *Client) Accept(address string, req *paxos.Request) error {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call accept RPC
	if _, err := paxos.NewPaxosClient(conn).Accept(context.Background(), &paxos.AcceptMsg{
		Request: req,
	}); err != nil {
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
func (c *Client) Commit() error {
	return nil
}
