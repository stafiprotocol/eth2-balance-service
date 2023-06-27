package task_ssv

import (
	"encoding/hex"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/pkg/credential"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) checkAndDeposit() error {
	poolBalance, err := task.userDepositContract.GetBalance(nil)
	if err != nil {
		return err
	}
	poolBalanceDeci := decimal.NewFromBigInt(poolBalance, 0)
	if poolBalanceDeci.LessThan(minAmountNeedDeposit) {
		return nil
	}

	// validators need stake
	depositLen := poolBalanceDeci.Div(minAmountNeedDeposit).IntPart()
	validatorPubkeys := make([][]byte, depositLen)
	sigs := make([][]byte, depositLen)
	dataRoots := make([][32]byte, depositLen)
	preKeyIndex := task.nextKeyIndex

	defer func() {
		if err := recover(); err != nil {
			task.nextKeyIndex = preKeyIndex
			for i := task.nextKeyIndex; i < task.nextKeyIndex+int(depositLen); i++ {
				delete(task.validators, i)
			}
		}
	}()

	for i := 0; i < int(depositLen); i++ {
		credential, err := credential.NewCredential(task.seed, task.nextKeyIndex, superNodeDepositAmount.BigInt(), task.chain, task.eth1WithdrawalAdress)
		if err != nil {
			return err
		}

		pubkeyBts := credential.SigningSk.PublicKey().Marshal()
		pubkeyStatus, err := task.mustGetSuperNodePubkeyStatus(pubkeyBts)
		if err != nil {
			return fmt.Errorf("mustGetSuperNodePubkeyStatus err: %s", err.Error())
		}
		if pubkeyStatus != utils.ValidatorStatusUnInitial {
			return fmt.Errorf("pubkey %s at index %d already on chain", hex.EncodeToString(pubkeyBts), task.nextKeyIndex)
		}

		validatorPubkeys[i] = pubkeyBts
		depositData, err := credential.SigningDepositData()
		if err != nil {
			return err
		}
		sigBts, err := hex.DecodeString(depositData.Signature)
		if err != nil {
			return err
		}
		sigs[i] = sigBts
		dataRootBts, err := hex.DecodeString(depositData.DepositDataRoot)
		if err != nil {
			return err
		}
		if len(dataRootBts) != 32 {
			return fmt.Errorf("dataRoot length %d not match", len(dataRootBts))
		}
		dataRoots[i] = [32]byte(dataRootBts)

		task.validators[task.nextKeyIndex] = &Validator{
			privateKey: credential.SigningSk,
			status:     utils.ValidatorStatusUnInitial,
			keyIndex:   task.nextKeyIndex,
		}

		task.nextKeyIndex++
	}

	err = task.superNodeConnection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.superNodeConnection.UnlockTxOpts()

	stakeTx, err := task.superNodeContract.Deposit(task.superNodeConnection.TxOpts(), validatorPubkeys, sigs, dataRoots)
	if err != nil {
		return err
	}

	err = utils.WaitTxOkCommon(task.superNodeConnection.Eth1Client(), stakeTx.Hash())
	if err != nil {
		return err
	}

	for _, pubkey := range validatorPubkeys {
		status, err := task.mustGetSuperNodePubkeyStatus(pubkey)
		if err != nil {
			return fmt.Errorf("mustGetSuperNodePubkeyStatus err: %s", err.Error())
		}
		if status == utils.ValidatorStatusUnInitial {
			return fmt.Errorf("validator %s not exist on chain", hex.EncodeToString(pubkey))
		}
	}

	return nil
}
