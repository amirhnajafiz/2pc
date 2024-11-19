package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/enums"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"

	"go.uber.org/zap"
)

// DatabaseHandler contains methods to perform database logic.
type DatabaseHandler struct {
	client  *client.Client
	logger  *zap.Logger
	storage *storage.Database
}

// Request accepts a transaction message and performs the needed logic to execute it (intr-shard).
func (d DatabaseHandler) Request(ra string, trx *database.TransactionMsg) {
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

	// call the reply RPC on client
	if err := d.client.Reply(ra, response, int(trx.GetSessionId())); err != nil {
		d.logger.Warn("failed to call reply", zap.String("client address", ra))
	}
}

func (d DatabaseHandler) Prepare(ra string, trx *database.TransactionMsg) {

}

func (d DatabaseHandler) Commit() {

}

func (d DatabaseHandler) Abort() {

}
