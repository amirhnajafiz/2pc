package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"

	"go.uber.org/zap"
)

// DatabaseHandler contains methods to perform database logic.
type DatabaseHandler struct {
	logger  *zap.Logger
	storage *storage.Database
}

// Request accepts a transaction message and performs the needed logic to execute it (intr-shard).
func (d DatabaseHandler) Request(trx *database.TransactionMsg) {
	// get the sender balance
	balance, err := d.storage.GetClientBalance(trx.GetSender())
	if err != nil {
		d.logger.Warn("failed to get client balance", zap.Error(err))
		return
	}

	// check the balance and transaction amount
	if trx.GetAmount() <= int64(balance) {
		// update both sender and receiver balance
		if err := d.storage.UpdateClientBalance(trx.GetSender(), -1*int(trx.GetAmount())); err != nil {
			d.logger.Warn("failed to update sender balance", zap.Error(err))
			return
		}
		if err := d.storage.UpdateClientBalance(trx.GetReceiver(), int(trx.GetAmount())); err != nil {
			d.logger.Warn("failed to update receiver balance", zap.Error(err))
			return
		}

		d.logger.Debug(
			"transaction submitted",
			zap.Int64("session id", trx.GetSessionId()),
		)
	} else {
		d.logger.Debug(
			"client balance is not enough to process the transaction",
			zap.Int64("session id", trx.GetSessionId()),
		)
	}
}

func (d DatabaseHandler) Prepare() {

}

func (d DatabaseHandler) Commit() {

}

func (d DatabaseHandler) Abort() {

}
