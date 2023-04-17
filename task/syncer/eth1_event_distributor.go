package task_syncer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) fetchDistributorContractEvents(start, end uint64) error {
	// unstake
	iterClaimed, err := task.distributorContract.FilterClaimed(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}
	for iterClaimed.Next() {
		txHashStr := iterClaimed.Event.Raw.TxHash.String()
		logIndex := uint32(iterClaimed.Event.Raw.Index)

		nodeClaim, err := dao_node.GetNodeClaim(task.db, txHashStr, logIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}
		nodeClaim.TxHash = txHashStr
		nodeClaim.LogIndex = logIndex

		nodeClaim.Address = iterClaimed.Event.Account.String()
		nodeClaim.BlockNumber = iterClaimed.Event.Raw.BlockNumber
		nodeClaim.ClaimableDeposit = decimal.NewFromBigInt(iterClaimed.Event.ClaimableDeposit, 0).StringFixed(0)
		nodeClaim.ClaimableReward = decimal.NewFromBigInt(iterClaimed.Event.ClaimableReward, 0).StringFixed(0)
		nodeClaim.ClaimedType = iterClaimed.Event.ClaimType

		block, err := task.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(nodeClaim.BlockNumber)))
		if err != nil {
			return err
		}
		nodeClaim.Timestamp = block.Header().Time

		err = dao_node.UpOrInNodeClaim(task.db, nodeClaim)
		if err != nil {
			return err
		}
	}

	// distributeSlash
	iterDistributerSlash, err := task.distributorContract.FilterDistributeSlash(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}
	for iterDistributerSlash.Next() {

		txHashStr := iterDistributerSlash.Event.Raw.TxHash.String()
		logIndex := uint32(iterDistributerSlash.Event.Raw.Index)

		distributerSlash, err := dao_node.GetDistributeSlash(task.db, txHashStr, logIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}
		distributerSlash.TxHash = txHashStr
		distributerSlash.LogIndex = logIndex

		distributerSlash.BlockNumber = iterDistributerSlash.Event.Raw.BlockNumber
		distributerSlash.SlashAmount = decimal.NewFromBigInt(iterDistributerSlash.Event.SlashAmount, 0).Div(utils.GweiDeci).BigInt().Uint64() // gwei
		distributerSlash.DealedHeight = iterDistributerSlash.Event.DealedHeight.Uint64()

		block, err := task.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(distributerSlash.BlockNumber)))
		if err != nil {
			return err
		}
		distributerSlash.Timestamp = block.Header().Time

		err = dao_node.UpOrInDistributeSlash(task.db, distributerSlash)
		if err != nil {
			return err
		}
	}

	return nil
}
