package task_syncer

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

// calc merkle tree on the basis of [node balance/slash event]
func (task *Task) calcMerkleTree() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	// todo mainnet config
	calMerkleTreeDu := task.rewardEpochInterval // 8 hours for test

	targetEpoch := (beaconHead.FinalizedEpoch / calMerkleTreeDu) * calMerkleTreeDu
	eth2NodeBalanceSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2NodeBalanceCollector)
	if err != nil {
		return err
	}
	// ensure node balances already caled
	if eth2NodeBalanceSyncerMetaData.DealedEpoch < targetEpoch {
		return nil
	}

	rootHash, err := dao.GetRootHash(task.db, targetEpoch)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	// just return if already cal
	if err == nil {
		return nil
	}

	// -- start cal
	nodeBalanceList, err := dao.GetNodeBalanceListByEpoch(task.db, targetEpoch)
	if err != nil {
		return err
	}

	proofList := make([]*dao.Proof, len(nodeBalanceList))
	for i, nodeBalance := range nodeBalanceList {
		proof, err := dao.GetProof(task.db, targetEpoch, nodeBalance.NodeAddress)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao.GetProof failed")
		}
		// fetch total slash amount
		valList, err := dao.GetValidatorListByNode(task.db, nodeBalance.NodeAddress, 0)
		if err != nil {
			return err
		}

		valIndexList := make([]uint64, len(valList))
		for i, val := range valList {
			valIndexList[i] = val.ValidatorIndex
		}
		totalSlashAmount, err := dao.GetTotalSlashAmountWithIndexList(task.db, valIndexList, targetEpoch)
		if err != nil {
			return errors.Wrap(err, "GetTotalSlashAmountWithIndexList failed")
		}
		totalSlashAmountDeci := decimal.NewFromInt(int64(totalSlashAmount)).
			Mul(utils.GweiDeci)

		totalExitDepositAmountDeci := decimal.NewFromInt(int64(nodeBalance.TotalExitNodeDepositAmount)).
			Mul(utils.GweiDeci)

		totalRewardAmountDeci := decimal.NewFromInt(int64(nodeBalance.TotalReward)).
			Mul(utils.GweiDeci).
			Sub(totalSlashAmountDeci)

		if totalRewardAmountDeci.IsNegative() {
			totalRewardAmountDeci = decimal.Zero
		}

		proof.Address = nodeBalance.NodeAddress
		proof.TotalRewardAmount = totalRewardAmountDeci.StringFixed(0)
		proof.TotalExitDepositAmount = totalExitDepositAmountDeci.StringFixed(0)
		proof.Index = uint32(i)
		proof.DealedEpoch = uint32(targetEpoch)

		proofList[i] = proof
	}

	tree, err := BuildMerkleTree(proofList)
	if err != nil {
		return err
	}
	treeHash, err := tree.GetRootHash()
	if err != nil {
		return err
	}

	// cal and save  proof
	for _, proof := range proofList {
		totalRewardAmountDeci, err := decimal.NewFromString(proof.TotalRewardAmount)
		if err != nil {
			return err
		}
		totalExitDepositAmountDeci, err := decimal.NewFromString(proof.TotalExitDepositAmount)
		if err != nil {
			return err
		}

		nodeHash := utils.GetNodeHash(big.NewInt(int64(proof.Index)), common.HexToAddress(proof.Address), totalRewardAmountDeci.BigInt(), totalExitDepositAmountDeci.BigInt())
		proofList, err := tree.GetProof(nodeHash)
		if err != nil {
			return err
		}

		proofStrList := make([]string, len(proofList))
		for i, p := range proofList {
			proofStrList[i] = p.String()
		}
		// set proof
		proof.Proof = strings.Join(proofStrList, ":")

		err = dao.UpOrInProof(task.db, proof)
		if err != nil {
			return err
		}
	}

	rootHash.DealedEpoch = uint32(targetEpoch)
	rootHash.RootHash = treeHash.String()
	err = dao.UpOrInRootHash(task.db, rootHash)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"epoch":    targetEpoch,
		"roothash": treeHash.String(),
	}).Info("cal merkleTree")

	return nil
}

func BuildMerkleTree(datas []*dao.Proof) (*utils.MerkleTree, error) {
	if len(datas) == 0 {
		return nil, fmt.Errorf("proof list empty")
	}
	list := make(utils.NodeHashList, len(datas))
	for i, data := range datas {
		totalRewardAmountDeci, err := decimal.NewFromString(data.TotalRewardAmount)
		if err != nil {
			return nil, err
		}
		totalExitDepositAmountDeci, err := decimal.NewFromString(data.TotalExitDepositAmount)
		if err != nil {
			return nil, err
		}
		list[i] = utils.GetNodeHash(big.NewInt(int64(data.Index)), common.HexToAddress(data.Address), totalRewardAmountDeci.BigInt(), totalExitDepositAmountDeci.BigInt())
	}
	mt := utils.NewMerkleTree(list)
	return mt, nil
}
