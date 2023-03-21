package task_voter

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"gorm.io/gorm"
)

func (task *Task) setMerkleTree() error {
	rootHash, err := dao_node.GetLatestRootHash(task.db)
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

	logrus.WithFields(logrus.Fields{
		"epoch":    rootHash.DealedEpoch,
		"roothash": rootHash.RootHash,
	}).Info("will set merkle root")
	return task.sendSetMerkleRootTx(rootHash)
}

func (task *Task) sendSetMerkleRootTx(rootHash *dao_node.RootHash) error {
	treeHash, err := hex.DecodeString(rootHash.RootHash)
	if err != nil {
		return errors.Wrap(err, "rootHash decode failed")
	}
	var merkleTreeRootHash [32]byte
	copy(merkleTreeRootHash[:], treeHash)

	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.SetMerkleRoot(task.connection.TxOpts(), big.NewInt(int64(rootHash.DealedEpoch)), merkleTreeRootHash)
	if err != nil {
		return err
	}

	logrus.Info("send SetMerkleRoot tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}
