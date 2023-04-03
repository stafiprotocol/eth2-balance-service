package task_syncer

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/eth/v1"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"github.com/stafiprotocol/eth2-balance-service/shared/types"
)

// get validator latest info of from beacon chain on finalized epoch
func (task *Task) syncValidatorLatestInfo() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	eth2ValidatorInfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	// no need fetch, if already dealed
	if finalEpoch <= eth2ValidatorInfoMetaData.DealedEpoch {
		return nil
	}

	targetEth1BlockHeight, err := task.getEpochStartBlocknumber(finalEpoch)
	if err != nil {
		return nil
	}

	eth1BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth1BlockSyncer)
	if err != nil {
		return err
	}

	// ensure all eth1 event synced before targetEth1BlockHeight
	if eth1BlockSyncerMetaData.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	validatorList, err := dao_node.GetValidatorListNeedFetchInfoFromBeacon(task.db)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"finalEpoch":       finalEpoch,
		"willDealEpoch":    finalEpoch,
		"validatorListLen": len(validatorList),
	}).Debug("syncValidatorLatestInfo")

	if len(validatorList) == 0 {
		eth2ValidatorInfoMetaData.DealedEpoch = finalEpoch
		return dao.UpOrInMetaData(task.db, eth2ValidatorInfoMetaData)
	}

	pubkeys := make([]types.ValidatorPubkey, 0)
	for _, validator := range validatorList {
		pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
		if err != nil {
			return err
		}
		pubkeys = append(pubkeys, pubkey)
	}

	willUsePubkeys := pubkeys
	var validatorStatusMap map[types.ValidatorPubkey]beacon.ValidatorStatus
	switch task.version {
	case utils.V1, utils.V2, utils.Dev:
		validatorStatusMap, err = task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
			Epoch: &finalEpoch,
		})
		if err != nil {
			return errors.Wrap(err, "syncValidatorLatestInfo GetValidatorStatuses failed")
		}
	default:
		return fmt.Errorf("unsupported version %s", task.version)
	}

	logrus.WithFields(logrus.Fields{
		"validatorStatuses len": len(validatorStatusMap),
	}).Debug("validator statuses")

	// ---- update validator info
	for pubkey, status := range validatorStatusMap {
		pubkeyStr := hexutil.Encode(pubkey.Bytes())
		if status.Exists {
			// must exist here
			validator, err := dao_node.GetValidator(task.db, pubkeyStr)
			if err != nil {
				return err
			}

			updateBaseInfo := func() {
				// validator's info may be inited at any status
				validator.ActiveEpoch = status.ActivationEpoch
				validator.EligibleEpoch = status.ActivationEligibilityEpoch
				validator.ValidatorIndex = status.Index

				exitEpoch := status.ExitEpoch
				if exitEpoch == math.MaxUint64 {
					exitEpoch = 0
				}
				withdrawableEpoch := status.WithdrawableEpoch
				if withdrawableEpoch == math.MaxUint64 {
					withdrawableEpoch = 0
				}

				validator.ExitEpoch = exitEpoch
				validator.WithdrawableEpoch = withdrawableEpoch
			}

			updateBalance := func() {
				validator.Balance = status.Balance
				validator.EffectiveBalance = status.EffectiveBalance
			}

			switch status.Status {

			case ethpb.ValidatorStatus_PENDING_INITIALIZED, ethpb.ValidatorStatus_PENDING_QUEUED: // pending
				validator.Status = utils.ValidatorStatusWaiting
				validator.ValidatorIndex = status.Index

			case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED: // active
				validator.Status = utils.ValidatorStatusActive
				if status.Slashed {
					validator.Status = utils.ValidatorStatusActiveSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED: // exited
				validator.Status = utils.ValidatorStatusExited
				if status.Slashed {
					validator.Status = utils.ValidatorStatusExitedSlash
				}
				updateBaseInfo()
				updateBalance()
			case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE: // withdrawable
				validator.Status = utils.ValidatorStatusWithdrawable
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawableSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_WITHDRAWAL_DONE: // withdrawdone
				validator.Status = utils.ValidatorStatusWithdrawDone
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawDoneSlash
				}
				updateBaseInfo()
				updateBalance()
			default:
				return fmt.Errorf("unsupported validator status %d", status.Status)
			}

			// cal total withdrawal
			totalWithdrawal, err := dao_node.GetValidatorTotalWithdrawalBeforeSlot(task.db, validator.ValidatorIndex, utils.StartSlotOfEpoch(task.eth2Config, finalEpoch))
			if err != nil {
				return err
			}
			validator.TotalWithdrawal = totalWithdrawal

			// cal total fee
			totalFee, err := task.calTotalFeeOfValidator(validator.ValidatorIndex, validator.NodeType, finalEpoch)
			if err != nil {
				return err
			}
			validator.TotalFee = totalFee

			err = dao_node.UpOrInValidator(task.db, validator)
			if err != nil {
				return err
			}
		}
	}

	// --- check distributed status
	needCheckDistributedValidatorList, err := dao_node.GetValidatorListNeedCheckDistributed(task.db)
	if err != nil {
		return err
	}

	poolInfo, err := dao_chaos.GetPoolInfo(task.db)
	if err != nil {
		return err
	}
	// tag distributed: withdrawDone && latest Distribute withdraw Height >  latest ValidatorWithdrawal.BlockNumber
	for _, val := range needCheckDistributedValidatorList {
		latestWithdrawal, err := dao_node.GetValidatorLatestWithdrawal(task.db, val.ValidatorIndex)
		if err != nil {
			logrus.Warnf("GetValidatorLatestWithdrawal failed, val: %d", val.ValidatorIndex)
			continue
		}
		// ensure withdrawDone
		if latestWithdrawal.Slot < utils.StartSlotOfEpoch(task.eth2Config, val.WithdrawableEpoch) {
			continue
		}

		// latest withdraw is partial withdrawal should skip
		if latestWithdrawal.Amount < utils.MaxPartialWithdrawalAmount {
			continue
		}

		if poolInfo.LatestDistributeWithdrawalHeight >= latestWithdrawal.BlockNumber {
			switch val.Status {
			case utils.ValidatorStatusWithdrawDone:
				val.Status = utils.ValidatorStatusDistributed
			case utils.ValidatorStatusWithdrawDoneSlash:
				val.Status = utils.ValidatorStatusDistributedSlash
			default:
				return fmt.Errorf("validator status: %d not match", val.Status)
			}

			err = dao_node.UpOrInValidator(task.db, val)
			if err != nil {
				return err
			}
		}
	}

	//----- tag validator ever slashed: ever slashed by protocol
	allValidatorList, err := dao_node.GetAllValidatorList(task.db)
	if err != nil {
		return err
	}
	for _, val := range allValidatorList {
		slashAmount, err := dao_node.GetTotalSlashAmountOfValidator(task.db, val.ValidatorIndex)
		if err != nil {
			return err
		}
		if slashAmount > 0 {
			val.EverSlashed = utils.ValidatorEverSlashedTrue
			err = dao_node.UpOrInValidator(task.db, val)
			if err != nil {
				return err
			}
		}
	}

	eth2ValidatorInfoMetaData.DealedEpoch = finalEpoch
	return dao.UpOrInMetaData(task.db, eth2ValidatorInfoMetaData)
}
