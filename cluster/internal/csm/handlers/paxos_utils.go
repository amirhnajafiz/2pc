package handlers

import (
	"fmt"
	"strings"

	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

func getClusterIPs(nodeName string, clusterName string, iptable map[string]string) []string {
	// split the cluster endpoints by ':'
	parts := strings.Split(iptable[fmt.Sprintf("E%s", clusterName)], ":")

	// ip list
	list := make([]string, 0)
	for _, key := range parts {
		if key != nodeName {
			list = append(list, iptable[key])
		}
	}

	return list
}

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
		Receiver:  req.GetSender(),
		Amount:    req.GetAmount(),
		SessionId: req.GetSessionId(),
	}
}
