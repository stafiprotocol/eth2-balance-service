package task_ssv

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/pkg/credential"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

var minAmountNeedStake = decimal.NewFromBigInt(big.NewInt(31), 18)
var minAmountNeedDeposit = decimal.NewFromBigInt(big.NewInt(32), 18)

var superNodeDepositAmount = decimal.NewFromBigInt(big.NewInt(1), 18)
var superNodeStakeAmount = decimal.NewFromBigInt(big.NewInt(31), 18)

func (task *Task) checkAndStake() error {
	poolBalance, err := task.userDepositContract.GetBalance(nil)
	if err != nil {
		return err
	}
	poolBalanceDeci := decimal.NewFromBigInt(poolBalance, 0)
	if poolBalanceDeci.LessThan(minAmountNeedStake) {
		return nil
	}

	// select validators need stake
	validatorsNeedStake := make([]*Validator, 0)
	for i := 0; i < len(task.validators); i++ {
		val := task.validators[i]
		if val.status != utils.ValidatorStatusWithdrawMatch {
			continue
		}

		if poolBalanceDeci.LessThan(minAmountNeedStake) {
			continue
		}

		validatorsNeedStake = append(validatorsNeedStake, val)
		poolBalanceDeci = poolBalanceDeci.Sub(superNodeStakeAmount)
	}

	// build payload
	validatorPubkeys := make([][]byte, len(validatorsNeedStake))
	sigs := make([][]byte, len(validatorsNeedStake))
	dataRoots := make([][32]byte, len(validatorsNeedStake))
	for i, val := range validatorsNeedStake {
		credential, err := credential.NewCredential(task.seed, val.keyIndex, superNodeStakeAmount.BigInt(), task.chain, task.eth1WithdrawalAdress)
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
	}

	// send tx
	err = task.superNodeConnection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.superNodeConnection.UnlockTxOpts()

	stakeTx, err := task.superNodeContract.Stake(task.superNodeConnection.TxOpts(), validatorPubkeys, sigs, dataRoots)
	if err != nil {
		return err
	}

	err = utils.WaitTxOkCommon(task.superNodeConnection.Eth1Client(), stakeTx.Hash())
	if err != nil {
		return err
	}
	return nil
}
