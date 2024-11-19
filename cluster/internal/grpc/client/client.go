package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is responsible for calling RPCs for CSMs.
type Client struct {
	nodeId string
}

// connect should be called in the beginning of each method to establish a connection.
func (c *Client) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to %s: %v", address, err)
	}

	return conn, nil
}

// NewClient returns an instance of Client.
func NewClient(nodeId string) *Client {
	return &Client{
		nodeId: nodeId,
	}
}
