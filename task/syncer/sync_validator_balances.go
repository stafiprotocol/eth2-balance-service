package task_syncer

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

// get staked validator info from beacon on target slot, and update balance/effective balance
func (task *Task) syncValidatorEpochBalances() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	finalizedEpoch := beaconHead.FinalizedEpoch

	metaData, err := dao.GetMetaData(task.db, utils.MetaTypeSyncBalances)
	if err != nil {
		return err
	}

	// no need fetch new balance
	if finalizedEpoch <= metaData.DealedEpoch {
		return nil
	}

	for epoch := metaData.DealedEpoch; epoch <= finalizedEpoch; epoch++ {
		if epoch%task.RewardEpochInterval != 0 {
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
		nodeAddressMap := make(map[string]struct{}, 0)
		for _, validator := range validatorList {
			pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
			if err != nil {
				return err
			}
			pubkeys = append(pubkeys, pubkey)
			pubkeyNodeMap[validator.Pubkey] = validator.NodeAddress
			nodeAddressMap[validator.NodeAddress] = struct{}{}
		}

		for i := 0; i < len(pubkeys); {
			start := i
			end := i + getValidatorStatusLimit
			if end > len(pubkeys) {
				end = len(pubkeys)
			}
			i = end

			willUsePubkeys := pubkeys[start:end]

			validatorStatusMap := make(map[types.ValidatorPubkey]beacon.ValidatorStatus)
			if task.FakeBeaconNode {
				pubkey, err := types.HexToValidatorPubkey("91b92af1781da257d3564a03f10c1f3b572695e1b4de50709096cf960260570768c17cd69c5a4ce6be9ae7e7f8e86f1f")
				if err != nil {
					return err
				}
				fakeStatus, err := task.connection.Eth2Client().GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{
					Epoch: &epoch,
				})
				if err != nil {
					return fmt.Errorf("GetValidatorStatus err: %s", err)
				}

				for _, pubkey := range willUsePubkeys {
					validatorStatusMap[pubkey] = fakeStatus
				}
			} else {
				validatorStatusMap, err = task.connection.Eth2Client().GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
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

				validatorBalance, err := dao.GetValidatorBalance(task.db, pubkeyStr, epoch)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}

				validatorBalance.NodeAddress = pubkeyNodeMap[pubkeyStr]
				validatorBalance.Balance = status.Balance
				validatorBalance.EffectiveBalance = status.EffectiveBalance
				validatorBalance.Epoch = epoch
				validatorBalance.Pubkey = pubkeyStr
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
				valInfo, err := dao.GetValidator(task.db, l.Pubkey)
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
