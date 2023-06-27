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

		validatorPubkeys[i] = credential.SigningSk.PublicKey().Marshal()
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

	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	stakeTx, err := task.superNodeContract.Deposit(task.connection.TxOpts(), validatorPubkeys, sigs, dataRoots)
	if err != nil {
		return err
	}

	err = utils.WaitTxOkCommon(task.connection.Eth1Client(), stakeTx.Hash())
	if err != nil {
		return err
	}

	return nil
}
