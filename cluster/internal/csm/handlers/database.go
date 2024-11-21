package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/lock"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/enums"
	"github.com/F24-CSE535/2pc/cluster/pkg/models"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"

	"go.uber.org/zap"
)

// DatabaseHandler contains methods to perform database logic.
type DatabaseHandler struct {
	client  *client.Client
	logger  *zap.Logger
	manager *lock.Manager
	memory  *memory.SharedMemory
	storage *storage.Database
}

// Request accepts a transaction message and performs the needed logic to execute it (intr-shard).
func (d DatabaseHandler) Request(ra string, trx *database.TransactionMsg) {
	// get sessionId
	sessionId := int(trx.GetSessionId())

	// release the locks
	defer func() {
		d.manager.Unlock(trx.GetSender(), sessionId)
		d.manager.Unlock(trx.GetReceiver(), sessionId)
	}()

	d.logger.Debug("input request", zap.String("req", trx.String()))

	// check the lock before request
	if !d.manager.Lock(trx.GetSender(), sessionId) || !d.manager.Lock(trx.GetReceiver(), sessionId) {
		d.logger.Warn("failed to capture locks", zap.Int("session id", sessionId))
		return
	}

	// get the sender balance
	balance, err := d.storage.GetClientBalance(trx.GetSender())
	if err != nil {
		d.logger.Warn("failed to get client balance", zap.Error(err))
		return
	}

	response := ""

	// check the balance and transaction amount
	if trx.GetAmount() <= int64(balance) {
		// update both sender and receiver balance
		if err := d.storage.UpdateClientBalance(trx.GetSender(), -1*int(trx.GetAmount()), false); err != nil {
			d.logger.Warn("failed to update sender balance", zap.Error(err))
			return
		}
		if err := d.storage.UpdateClientBalance(trx.GetReceiver(), int(trx.GetAmount()), false); err != nil {
			d.logger.Warn("failed to update receiver balance", zap.Error(err))
			return
		}

		response = enums.RespOK

		d.logger.Debug(
			"transaction submitted",
			zap.Int64("session id", trx.GetSessionId()),
		)
	} else {
		response = enums.RespFailed

		d.logger.Debug(
			"client balance is not enough to process the transaction",
			zap.Int64("session id", trx.GetSessionId()),
		)
	}

	// call the reply RPC on client, if the node is the leader
	if d.memory.GetLeader() == d.memory.GetNodeName() {
		if err := d.client.Reply(ra, response, int(trx.GetSessionId())); err != nil {
			d.logger.Warn("failed to call reply", zap.String("client address", ra))
		}
	}
}

// Prepare accepts a prepare message and returns ack to the sender.
func (d DatabaseHandler) Prepare(msg *database.PrepareMsg) {
	// get sessionId
	sessionId := int(msg.Transaction.GetSessionId())

	// check the lock before request
	if !d.manager.Lock(msg.GetClient(), sessionId) {
		d.logger.Warn("failed to capture locks", zap.Int("session id", sessionId))

		// release the locks if one captured
		d.manager.Unlock(msg.GetClient(), sessionId)

		return
	}

	// create a list of WALs
	wals := make([]*models.Log, 0)
	wals = append(wals, &models.Log{SessionId: sessionId, Message: enums.WALStart})

	// abort flag
	abort := false

	// for (S, R, amount) we want to check if the client is S to check its balance
	if msg.GetTransaction().GetSender() == msg.GetClient() {
		// get our client balance
		balance, err := d.storage.GetClientBalance(msg.GetClient())
		if err != nil {
			d.logger.Warn("failed to get client balance", zap.Error(err))
			return
		}

		// add a log
		wals = append(wals, &models.Log{
			SessionId: sessionId,
			Message:   enums.WALUpdate,
			Record:    msg.GetTransaction().GetSender(),
			NewValue:  -1 * int(msg.GetTransaction().GetAmount())},
		)

		// check if the balance is enough
		if msg.GetTransaction().GetAmount() > int64(balance) {
			// the balance is not enough
			abort = true
		}
	} else {
		// add a log
		wals = append(wals, &models.Log{
			SessionId: sessionId,
			Message:   enums.WALUpdate,
			Record:    msg.GetTransaction().GetReceiver(),
			NewValue:  int(msg.GetTransaction().GetAmount())},
		)
	}

	// store the logs
	if err := d.storage.InsertBatchWAL(wals); err != nil {
		d.logger.Warn("failed to store logs", zap.Error(err))
		return
	}

	// send the ack message, if the node is leader
	if d.memory.GetLeader() == d.memory.GetNodeName() {
		if err := d.client.Ack(msg.GetReturnAddress(), sessionId, abort); err != nil {
			d.logger.Warn("failed to send ack message", zap.Error(err))
			return
		}
	}
}

// Commit accepts a commit message to get WALs and update the records.
func (d DatabaseHandler) Commit(msg *database.CommitMsg) {
	// get sessionId
	sessionId := int(msg.GetSessionId())

	// get all update logs
	wals, err := d.storage.RetrieveWALs(sessionId)
	if err != nil {
		d.logger.Warn("failed to get logs", zap.Error(err))
		return
	}

	// update clients
	for _, wal := range wals {
		// release the locks if one captured
		d.manager.Unlock(wal.Record, sessionId)

		// update the client balance
		if err := d.storage.UpdateClientBalance(wal.Record, wal.NewValue, false); err != nil {
			d.logger.Warn("failed to update balance", zap.Error(err), zap.String("client", wal.Record))
			return
		}
	}

	// store a commit log
	if err := d.storage.InsertWAL(&models.Log{
		SessionId: sessionId,
		Message:   enums.WALCommit,
	}); err != nil {
		d.logger.Warn("failed to store log", zap.Error(err))
		return
	}

	d.logger.Debug(
		"transaction submitted",
		zap.Int("session id", sessionId),
	)

	// call the reply RPC on client, if the node is leader
	if d.memory.GetLeader() == d.memory.GetNodeName() {
		if err := d.client.Reply(msg.GetReturnAddress(), enums.RespOK, sessionId); err != nil {
			d.logger.Warn("failed to call reply", zap.String("client address", msg.GetReturnAddress()))
		}
	}
}

// Abort will log an abort log into the logs.
func (d DatabaseHandler) Abort(sessionId int) {
	// get all update logs
	wals, err := d.storage.RetrieveWALs(sessionId)
	if err != nil {
		d.logger.Warn("failed to get logs", zap.Error(err))
		return
	}

	// release the locks if one captured
	for _, wal := range wals {
		d.manager.Unlock(wal.Record, sessionId)
	}

	// insert a abort WAL
	if err := d.storage.InsertWAL(&models.Log{
		SessionId: sessionId,
		Message:   enums.WALAbort,
	}); err != nil {
		d.logger.Warn("failed to store log", zap.Error(err))
	}
}
