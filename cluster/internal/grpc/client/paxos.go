package client

import "github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

// Accept gets a target and request to call Accept RPC on the target.
func (c *Client) Accept(target string, req *paxos.Request) error {
	return nil
}

// Accepted gets a target to call Accepted RPC on the target.
func (c *Client) Accepted() error {
	return nil
}

// Commit gets a target to call Commit RPC on the target.
func (c *Client) Commit() error {
	return nil
}
