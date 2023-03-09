package task_voter

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) setMerkleTree() error {
	rootHash, err := dao.GetLatestRootHash(task.db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}

	dealedEpochOnchain, err := task.MerkleTreeDealedEpoch(task.storageContract)
	if err != nil {
		return err
	}
	// ensure not set
	if dealedEpochOnchain.Uint64() >= uint64(rootHash.DealedEpoch) {
		return nil
	}

	// --- start set merkle root
	logrus.WithFields(logrus.Fields{
		"epoch":    rootHash.DealedEpoch,
		"roothash": rootHash.RootHash,
	}).Info("will set merkle root")

	// -----3 send vote tx
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	treeHash, err := hex.DecodeString(rootHash.RootHash)
	if err != nil {
		return errors.Wrap(err, "rootHash decode failed")
	}
	var merkleTreeHash [32]byte
	copy(merkleTreeHash[:], treeHash)
	tx, err := task.distributorContract.SetMerkleRoot(task.connection.TxOpts(), big.NewInt(int64(rootHash.DealedEpoch)), merkleTreeHash)
	if err != nil {
		return err
	}

	logrus.Info("send SetMerkleRoot tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("SetMerkleRoot tx reach retry limit")
		}
		_, pending, err := task.connection.Eth1Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx status")
			} else {
				logrus.WithFields(logrus.Fields{
					"hash":   tx.Hash(),
					"status": "pending",
				}).Warn("tx status")
			}
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("SetMerkleRoot tx send ok")

	return nil
}
