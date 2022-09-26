package task_v1_syncer

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/types"
	"gorm.io/gorm"
)

const fetchValidatorStatusLimit = 50

func (task *Task) syncV1ValidatorEpochBalances() error {
	finalEpoch := utils.V1EndEpoch

	metaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BalanceSyncer)
	if err != nil {
		return err
	}
	logrus.Debugf("final epoch: %d, dealedEpoch: %s", finalEpoch, metaData.DealedEpoch)
	// no need fetch new balance
	if finalEpoch <= metaData.DealedEpoch {
		return nil
	}

	for epoch := metaData.DealedEpoch + 1; epoch <= finalEpoch; epoch++ {
		logrus.Debugf("check epoch: %d", epoch)
		if epoch%task.rewardEpochInterval != 0 {
			continue
		}

		validatorList, err := dao.GetValidatorListActiveEpochBefore(task.db, epoch)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"epoch":         epoch,
			"validatorList": validatorList,
		}).Debug("activeValidators")
		if len(validatorList) == 0 {
			return nil
		}

		pubkeys := make([]types.ValidatorPubkey, 0)
		pubkeyNodeMap := make(map[string]string)
		nodeAddressMap := make(map[string]struct{})
		pubkeyToIndex := make(map[string]uint64)
		for _, validator := range validatorList {
			pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
			if err != nil {
				return err
			}
			pubkeys = append(pubkeys, pubkey)
			pubkeyNodeMap[validator.Pubkey] = validator.NodeAddress
			nodeAddressMap[validator.NodeAddress] = struct{}{}
			pubkeyToIndex[validator.Pubkey] = validator.ValidatorIndex
		}

		for i := 0; i < len(pubkeys); {
			start := i
			end := i + fetchValidatorStatusLimit
			if end > len(pubkeys) {
				end = len(pubkeys)
			}
			i = end

			willUsePubkeys := pubkeys[start:end]

			validatorStatusMap := make(map[types.ValidatorPubkey]beacon.ValidatorStatus)
			if task.fakeBeaconNode {

				for _, pubkey := range willUsePubkeys {
					index := 21100 + int(pubkey.Bytes()[5]) + int(pubkey.Bytes()[25])
					fakeStatus, err := task.connection.GetValidatorStatusByIndex(fmt.Sprint(index), &beacon.ValidatorStatusOptions{
						Epoch: &epoch,
					})
					if err != nil {
						return fmt.Errorf("GetValidatorStatus err: %s", err)
					}
					validatorStatusMap[pubkey] = fakeStatus
				}
			} else {
				validatorStatusMap, err = task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
					Epoch: &epoch,
				})
				if err != nil {
					return err
				}

			}

			logrus.WithFields(logrus.Fields{
				"validatorStatuses": validatorStatusMap,
			}).Debug("validator statuses")

			for pubkey, status := range validatorStatusMap {
				pubkeyStr := hexutil.Encode(pubkey.Bytes())
				if !status.Exists {
					return fmt.Errorf("should exist status on beacon chain, pubkey: %s, epoch: %d", pubkeyStr, epoch)
				}
				validatorIndex, exist := pubkeyToIndex[pubkeyStr]
				if !exist {
					return fmt.Errorf("validator index not exit in pubkeyToIndex")
				}
				validatorBalance, err := dao.GetValidatorBalance(task.db, validatorIndex, epoch)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}

				validatorBalance.NodeAddress = pubkeyNodeMap[pubkeyStr]
				validatorBalance.Balance = status.Balance
				validatorBalance.EffectiveBalance = status.EffectiveBalance
				validatorBalance.Epoch = epoch
				validatorBalance.ValidatorIndex = validatorIndex
				validatorBalance.Timestamp = utils.EpochTime(task.eth2Config, epoch)

				err = dao.UpOrInValidatorBalance(task.db, validatorBalance)
				if err != nil {
					return err
				}

			}
		}

		// collect node address
		for node := range nodeAddressMap {

			list, err := dao.GetValidatorBalanceList(task.db, node, epoch)
			if err != nil {
				return err
			}

			nodeBalance, err := dao.GetNodeBalance(task.db, node, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			nodeBalance.NodeAddress = node
			nodeBalance.Epoch = epoch
			nodeBalance.Timestamp = utils.EpochTime(task.eth2Config, epoch)

			for _, l := range list {
				valInfo, err := dao.GetValidatorByIndex(task.db, l.ValidatorIndex)
				if err != nil {
					return err
				}

				nodeBalance.TotalNodeDepositAmount += valInfo.NodeDepositAmount

				nodeBalance.TotalBalance += l.Balance
				nodeBalance.TotalEffectiveBalance += l.EffectiveBalance
			}

			err = dao.UpOrInNodeBalance(task.db, nodeBalance)
			if err != nil {
				return err
			}

		}

		metaData.DealedEpoch = epoch
		err = dao.UpOrInMetaData(task.db, metaData)
		if err != nil {
			return err
		}

	}
	return nil
}