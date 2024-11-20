package services

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"google.golang.org/protobuf/types/known/emptypb"
)

// PaxosService is used for handling paxos RPCs.
type PaxosService struct {
	paxos.UnimplementedPaxosServer

	Channel chan *packets.Packet
}

// Accept RPC creates a new accept packet and sends it to the CSMs.
func (p *PaxosService) Accept(_ context.Context, msg *paxos.AcceptMsg) (*emptypb.Empty, error) {
	p.Channel <- &packets.Packet{Label: packets.PktPaxosAccept, Payload: msg}

	return &emptypb.Empty{}, nil
}

// Accepted RPC creates a new accepted packet and sends it to the CSMs.
func (p *PaxosService) Accepted(_ context.Context, msg *paxos.AcceptedMsg) (*emptypb.Empty, error) {
	p.Channel <- &packets.Packet{Label: packets.PktPaxosAccepted, Payload: msg}

	return &emptypb.Empty{}, nil
}

// Commit RPC creates a new commit packet and sends it to the CSMs.
func (p *PaxosService) Commit(_ context.Context, msg *paxos.CommitMsg) (*emptypb.Empty, error) {
	p.Channel <- &packets.Packet{Label: packets.PktPaxosCommit, Payload: msg}

	return &emptypb.Empty{}, nil
}
