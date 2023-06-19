package task_voter

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// sync rate from eth to arbitrum
func (task *Task) syncRate() error {
	rateOnEth, err := task.rethContract.GetExchangeRate(nil)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}

	rateOnArbitrum, err := task.arbitrumStakePortalRateContract.GetRate(nil)
	if err != nil {
		return fmt.Errorf("arbitrumStakePortalRateContract.GetRate err: %s", err)
	}

	if rateOnEth.Cmp(rateOnArbitrum) != 0 {
		logrus.WithFields(logrus.Fields{"rateOnEth": rateOnEth.String(), "rateOnArbitrum": rateOnArbitrum.String()}).Info("rateInfo")

		proposalId := getProposalId(0, rateOnEth, 0)
		proposal, err := task.arbitrumStakePortalRateContract.Proposals(nil, proposalId)
		if err != nil {
			return fmt.Errorf("arbitrumStakePortalRateContract Proposals error %s ", err)
		}
		if proposal.Status == 2 { // success status
			return nil
		}
		hasVoted, err := task.arbitrumStakePortalRateContract.HasVoted(nil, proposalId, task.arbitrumConn.TxOpts().From)
		if err != nil {
			return fmt.Errorf("arbitrumStakePortalRateContract HasVoted error %s", err)
		}
		if hasVoted {
			return nil
		}

		// send tx
		err = task.arbitrumConn.LockAndUpdateTxOpts()
		if err != nil {
			return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
		}
		defer task.arbitrumConn.UnlockTxOpts()

		tx, err := task.arbitrumStakePortalRateContract.VoteRate(task.arbitrumConn.TxOpts(), proposalId, rateOnEth)
		if err != nil {
			return err
		}

		logrus.Info("send VoteRate tx hash: ", tx.Hash().String())

		err = task.waitArbitrumTxOk(tx.Hash())
		if err != nil {
			return fmt.Errorf("waitArbitrumTxOk failed, err: %s", err)
		}
		return task.waitRateUpdated(proposalId)
	}

	return nil
}

func getProposalId(era uint32, rate *big.Int, factor int) common.Hash {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-%s-%s-%d", era, "voteRate", rate.String(), factor)))
}

func (task *Task) waitRateUpdated(proposalId [32]byte) error {
	retry := 0
	for {
		if retry > utils.RetryLimit*3 {
			return fmt.Errorf("waitRateUpdated tx reach retry limit")
		}

		proposal, err := task.arbitrumStakePortalRateContract.Proposals(nil, proposalId)
		if err != nil {
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
		if proposal.Status != 2 {
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
		break
	}
	return nil
}
