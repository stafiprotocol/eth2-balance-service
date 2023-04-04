package task_voter

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) distributeSuperNodeFeePool() error {
	latestDistributeHeight, targetEth1BlockHeight, shouldGoNext, err := task.checkStateForDistriSuperNodeFeePool()
	if err != nil {
		return errors.Wrap(err, "distributeSuperNodeFeePool checkSyncState failed")
	}

	if !shouldGoNext {
		logrus.Debug("distributeSuperNodeFeePool should not go next")
		return nil
	}

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, err := task.getUserNodePlatformFromSuperNodeFeePool(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromFeePool failed")
	}

	// return if smaller than minDistributeAmount
	if totalAmountDeci.IsZero() {
		logrus.Debugf("distributeSuperNodeFeePool totalAmountDeci: %s ", totalAmountDeci.String())
		return nil
	}
	// check voted
	voted, err := task.NodeVotedDistributeSuperNodeFeePool(task.storageContract, task.connection.Keypair().CommonAddress(),
		big.NewInt(int64(targetEth1BlockHeight)), totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt())
	if err != nil {
		return fmt.Errorf("NodeVotedDistributeFeePoolerr: %s", err)
	}
	if voted {
		logrus.Debug("NodeVotedDistributeFeePool voted")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1BlockHeight": targetEth1BlockHeight,
		"totalUserEthDeci":      totalUserEthDeci.String(),
		"totalNodeEthDeci":      totalNodeEthDeci.String(),
		"totalPlatformEthDeci":  totalPlatformEthDeci.String(),
	}).Info("Will distributeSuperNodeFeePool")

	return task.sendDistributeSuperNodeFeeTx(big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt())
}

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext, err)
func (task *Task) checkStateForDistriSuperNodeFeePool() (uint64, uint64, bool, error) {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	targetEpoch := (finalEpoch / task.rewardEpochInterval) * task.rewardEpochInterval
	targetEth1BlockHeight, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	logrus.Debugf("targetEth1Block %d", targetEth1BlockHeight)

	latestDistributeHeight, err := task.distributorContract.GetDistributeSuperNodeFeeDealedHeight(task.connection.CallOpts(nil))
	if err != nil {
		return 0, 0, false, err
	}
	// init case
	if latestDistributeHeight.Uint64() == 0 {
		latestDistributeHeight = big.NewInt(task.distributeSuperNodeFeeInitDealedHeight)
	}

	if latestDistributeHeight.Uint64() >= targetEth1BlockHeight {
		logrus.Debug("latestDistributeHeight >= targetEth1BlockHeight")
		return 0, 0, false, nil
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	eth2BlockSyncerBlockHeight, err := task.getEpochStartBlocknumber(eth2BlockSyncerMetaData.DealedEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	// ensure eth2 block have synced
	if eth2BlockSyncerBlockHeight < targetEth1BlockHeight {
		logrus.Debugf("eth2BlockSyncerBlockHeight %d < targetEth1BlockHeight %d", eth2BlockSyncerBlockHeight, targetEth1BlockHeight)
		return 0, 0, false, nil
	}

	return latestDistributeHeight.Uint64(), targetEth1BlockHeight, true, nil
}

func (task *Task) sendDistributeSuperNodeFeeTx(targetEth1BlockHeight, totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci *big.Int) error {
	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.DistributeSuperNodeFee(task.connection.TxOpts(),
		targetEth1BlockHeight, totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci)
	if err != nil {
		return err
	}
	logrus.Info("send DistributeFee tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

// return (user reward, node reward, platform fee, totalFee) decimals 18
func (task Task) getUserNodePlatformFromSuperNodeFeePool(latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	proposedBlockList, err := dao_node.GetProposedBlockListBetween(task.db, latestDistributeHeight, targetEth1BlockHeight, task.superNodeFeePoolAddress.String())
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
	}

	totalAmountDeci := decimal.Zero
	totalUserEthDeci := decimal.Zero
	totalNodeEthDeci := decimal.Zero
	totalPlatformEthDeci := decimal.Zero
	for _, w := range proposedBlockList {
		validator, err := dao_node.GetValidatorByIndex(task.db, w.ValidatorIndex)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
		}
		feeAmountDeci, err := decimal.NewFromString(w.FeeAmount)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
		}

		totalAmountDeci = totalAmountDeci.Add(feeAmountDeci)

		// cal rewards
		var userRewardDeci, nodeRewardDeci, platformFeeDeci = decimal.Zero, decimal.Zero, decimal.Zero
		if w.Slot <= utils.StartSlotOfEpoch(task.eth2Config, task.rewardV1EndEpoch) {
			userRewardDeci, nodeRewardDeci, platformFeeDeci = utils.GetUserNodePlatformRewardV1(validator.NodeDepositAmount, feeAmountDeci)
		} else {

			userRewardDeci, nodeRewardDeci, platformFeeDeci = utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, feeAmountDeci)
		}
		// cal reward + deposit
		totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci)
		totalNodeEthDeci = totalNodeEthDeci.Add(nodeRewardDeci)
		totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)

	}

	return totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, nil
}
