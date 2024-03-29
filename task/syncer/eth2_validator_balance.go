package task_syncer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"github.com/stafiprotocol/eth2-balance-service/shared/types"
	"gorm.io/gorm"
)

// calc validator epoch balance info(balance/effective balance/totalwithdrawal/totalfee) at target epoch(every 75 epoch) on the basis of [beaconchain/proposed blocks/withdrawals]
func (task *Task) syncValidatorEpochBalances() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch
	eth2ValidatorInfoSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		return err
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return err
	}
	// ensure validators latest info already synced
	if finalEpoch > eth2ValidatorInfoSyncerMetaData.DealedEpoch {
		finalEpoch = eth2ValidatorInfoSyncerMetaData.DealedEpoch
	}
	// ensure validators block info(withdrawals) already synced
	if finalEpoch > eth2BlockSyncerMetaData.DealedEpoch {
		finalEpoch = eth2BlockSyncerMetaData.DealedEpoch
	}

	eth2ValidatorBalanceMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return err
	}

	// no need fetch new balance
	if finalEpoch <= eth2ValidatorBalanceMetaData.DealedEpoch {
		return nil
	}

	for epoch := eth2ValidatorBalanceMetaData.DealedEpoch + 1; epoch <= finalEpoch; epoch++ {
		// we fetch epoch info every 75 epoch
		if epoch%task.rewardEpochInterval != 0 {
			continue
		}

		validatorList, err := dao_node.GetValidatorListActiveEpochBefore(task.db, epoch)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"dealedEpoch":              eth2ValidatorBalanceMetaData.DealedEpoch,
			"willDealEpoch":            epoch,
			"willDealValidatorListLen": len(validatorList),
		}).Debug("syncValidatorEpochBalances")

		// should skip if no validator
		if len(validatorList) == 0 {
			eth2ValidatorBalanceMetaData.DealedEpoch = epoch
			err = dao.UpOrInMetaData(task.db, eth2ValidatorBalanceMetaData)
			if err != nil {
				return err
			}
			continue
		}

		pubkeys := make([]types.ValidatorPubkey, 0)
		pubkeyToNodeAddress := make(map[string]string)
		pubkeyToValidator := make(map[string]*dao_node.Validator)
		for _, validator := range validatorList {
			pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
			if err != nil {
				return err
			}
			pubkeys = append(pubkeys, pubkey)
			pubkeyToNodeAddress[validator.Pubkey] = validator.NodeAddress
			pubkeyToValidator[validator.Pubkey] = validator
		}
		willUsePubkeys := pubkeys

		// fetch validators info at target epoch
		validatorStatusMap, err := task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
			Epoch: &epoch,
		})
		if err != nil {
			return errors.Wrap(err, "syncValidatorEpochBalances GetValidatorStatuses failed")
		}

		logrus.WithFields(logrus.Fields{
			"validatorStatuses len": len(validatorStatusMap),
		}).Debug("validator statuses")

		if len(validatorStatusMap) != len(willUsePubkeys) {
			return fmt.Errorf("validatorStatusMap len: %d not equal pubkeys len: %d", len(validatorStatusMap), len(willUsePubkeys))
		}

		for pubkey, status := range validatorStatusMap {
			pubkeyStr := hexutil.Encode(pubkey.Bytes())
			if !status.Exists {
				return fmt.Errorf("status not exist on beacon chain, pubkey: %s, epoch: %d", pubkeyStr, epoch)
			}
			validatorIndex := status.Index
			nodeAddress, exist := pubkeyToNodeAddress[pubkeyStr]
			if !exist {
				return fmt.Errorf("node address not exist in pubkeyToNodeAddress")
			}
			validatorInfo, exist := pubkeyToValidator[pubkeyStr]
			if !exist {
				return fmt.Errorf("validator info not exist in pubkeyToValidator")
			}

			// calc total withdrawal
			totalWithdrawal, err := dao_node.GetValidatorTotalWithdrawalBeforeSlot(task.db, validatorIndex, utils.StartSlotOfEpoch(task.eth2Config, epoch))
			if err != nil {
				return errors.Wrap(err, "GetValidatorTotalWithdrawalBeforeSlot failed")
			}

			// calc total fee to fee pool/super fee pool
			totalFee, err := task.calTotalFeeOfValidator(validatorIndex, validatorInfo.NodeType, epoch)
			if err != nil {
				return errors.Wrap(err, "calTotalFeeOfValidator failed")
			}

			// insert valdiator balance
			validatorBalance, err := dao_node.GetValidatorBalance(task.db, validatorIndex, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			validatorBalance.NodeAddress = nodeAddress
			validatorBalance.Balance = status.Balance
			validatorBalance.TotalWithdrawal = totalWithdrawal
			validatorBalance.TotalFee = totalFee
			validatorBalance.EffectiveBalance = status.EffectiveBalance
			validatorBalance.Epoch = epoch
			validatorBalance.ValidatorIndex = validatorIndex
			validatorBalance.Timestamp = utils.StartTimestampOfEpoch(task.eth2Config, epoch)

			// !!!!!safe check
			if validatorBalance.Balance+validatorBalance.TotalWithdrawal > utils.StandardEffectiveBalance+utils.MaxPartialWithdrawalAmount {
				return fmt.Errorf("validator: %d staking reward abnormal, balance: %d, withdrawals: %d", validatorBalance.ValidatorIndex, validatorBalance.Balance, validatorBalance.TotalWithdrawal)
			}
			if validatorBalance.TotalFee > utils.StandardEffectiveBalance {
				return fmt.Errorf("validator: %d priority fee abnormal, totalFee: %d", validatorBalance.ValidatorIndex, validatorBalance.TotalFee)
			}

			err = dao_node.UpOrInValidatorBalance(task.db, validatorBalance)
			if err != nil {
				return err
			}
		}

		eth2ValidatorBalanceMetaData.DealedEpoch = epoch
		err = dao.UpOrInMetaData(task.db, eth2ValidatorBalanceMetaData)
		if err != nil {
			return err
		}

	}
	return nil
}

// return gwei
func (task *Task) calTotalFeeOfValidator(validatorIndex uint64, nodeType uint8, epoch uint64) (uint64, error) {
	feePoolAddressList := []string{task.lightNodeFeePoolAddress.String(), task.superNodeFeePoolAddress.String()}

	proposedBlockList, err := dao_node.GetProposedBlockListBefore(task.db, validatorIndex, utils.StartSlotOfEpoch(task.eth2Config, epoch), feePoolAddressList)
	if err != nil {
		return 0, errors.Wrap(err, "GetProposedBlockListBefore failed")
	}
	// we use gwei here
	totalFee := uint64(0)
	for _, block := range proposedBlockList {
		feeAmountDeci, err := decimal.NewFromString(block.FeeAmount)
		if err != nil {
			return 0, errors.Wrap(err, "fee amount cast decimal failed")
		}
		totalFee += feeAmountDeci.Div(utils.GweiDeci).BigInt().Uint64()
	}
	return totalFee, nil
}
