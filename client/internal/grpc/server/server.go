package server

import (
	"context"

	"github.com/F24-CSE535/2pc/client/pkg/enums"
	"github.com/F24-CSE535/2pc/client/pkg/models"
	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Server is a partial gRPC server in client side.
type Server struct {
	database.UnimplementedDatabaseServer

	channel chan *models.Packet
}

// Reply accepts all reply messages.
func (s *Server) Reply(_ context.Context, msg *database.ReplyMsg) (*emptypb.Empty, error) {
	s.channel <- &models.Packet{Label: enums.PktReply, Payload: msg}

	return &emptypb.Empty{}, nil
}

// Ack accepts all ack messages.
func (s *Server) Ack(_ context.Context, msg *database.AckMsg) (*emptypb.Empty, error) {
	s.channel <- &models.Packet{Label: enums.PktAck, Payload: msg}

	return &emptypb.Empty{}, nil
}
