package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

func databaseRequestToPaxosRequest(req *database.RequestMsg) *paxos.Request {
	return &paxos.Request{
		Sender:        req.GetTransaction().GetSender(),
		Receiver:      req.GetTransaction().GetReceiver(),
		Amount:        req.GetTransaction().GetAmount(),
		SessionId:     req.GetTransaction().GetSessionId(),
		ReturnAddress: req.GetReturnAddress(),
	}
}

func databasePrepareToPaxosRequest(req *database.PrepareMsg) *paxos.Request {
	return &paxos.Request{
		Sender:        req.GetTransaction().GetSender(),
		Receiver:      req.GetTransaction().GetReceiver(),
		Amount:        req.GetTransaction().GetAmount(),
		SessionId:     req.GetTransaction().GetSessionId(),
		Client:        req.GetClient(),
		ReturnAddress: req.GetReturnAddress(),
	}
}

func paxosRequestToDatabaseTransaction(req *paxos.Request) *database.TransactionMsg {
	return &database.TransactionMsg{
		Sender:    req.GetSender(),
		Receiver:  req.GetReceiver(),
		Amount:    req.GetAmount(),
		SessionId: req.GetSessionId(),
	}
}
