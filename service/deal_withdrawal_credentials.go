package service

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/reth/bindings/StakingPool"
	"github.com/stafiprotocol/reth/types"
)

var poolPrelaunchBatchSize = int64(50)

func (s *Service) dealWithdrawalCredentials(ctx context.Context, sysErr chan<- error) error {
	glog.Info("dealWithdrawalCredentials...")
	var retry = BlockRetryLimit
	for {
		select {
		case <-ctx.Done():
			return errors.New("dealWithdrawalCredentials terminated")
		default:
			// No more retries
			if retry <= 0 {
				glog.Error("dealWithdrawalCredentials failed, retries exceeded")
				sysErr <- ErrFatalDealWithdrawalCredentials
				return nil
			}
			err := s.checkAndVoteCredentials()
			if err != nil {
				glog.Error("dealWithdrawalCredentials failed", "error", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			time.Sleep(BlockRetryInterval)
			retry = BlockRetryLimit
		}
	}
}

func (s *Service) checkAndVoteCredentials() error {
	poolAddresses, err := s.getPrelaunchStakingpools()
	if err != nil {
		return err
	}
	if len(poolAddresses) == 0 {
		return nil
	}

	expectCredentials, err := s.contract.St.GetWithdrawalCredentials(s.conn.callOpts)
	if err != nil {
		return err
	}

	for _, poolAddress := range poolAddresses {
		stakingPoolContract, err := StakingPool.NewStakingPool(poolAddress, s.conn.Client())
		if err != nil {
			return err
		}

		match, err := stakingPoolContract.GetWithdrawalCredentialsMatch(s.conn.callOpts)
		if err != nil {
			return err
		}
		if match {
			continue
		}
		pubkey, err := s.contract.StafiStakingPoolContract.GetStakingPoolPubkey(s.conn.callOpts, poolAddress)
		if err != nil {
			return err
		}

		validatorStatus, err := s.conn.eth2Conn.GetValidatorStatus(types.BytesToValidatorPubkey(pubkey), nil)
		if err != nil {
			return err
		}

		if validatorStatus.Exists {
			if validatorStatus.WithdrawalCredentials == common.BytesToHash(expectCredentials) {
				tx, err := s.voteWithdrawCredentials(stakingPoolContract)
				if err != nil {
					return err
				}
				if tx != "" {
					glog.Info("VoteWithdrawCredentials", "tx", tx, "poolAddress", poolAddress.Hex())
				}
			}
		}
	}

	return nil
}

func (s *Service) voteWithdrawCredentials(stakingPoolContract *StakingPool.StakingPool) (string, error) {
	err := s.contract.Conn.LockAndUpdateOpts()
	if err != nil {
		return "", err
	}
	defer s.contract.Conn.UnlockOpts()

	tx, err := stakingPoolContract.VoteWithdrawCredentials(s.conn.opts)
	if err != nil {
		if strings.Contains(err.Error(), "Member has already voted to withdrawCredentials") {
			return "", nil
		} else {
			return "", err
		}
	}
	return tx.Hash().String(), nil
}

func (s *Service) getPrelaunchStakingpools() ([]common.Address, error) {
	poolCount, err := s.contract.StafiStakingPoolContract.GetStakingPoolCount(s.conn.callOpts)
	if err != nil {
		glog.Error("dealWithdrawalCredentials", "GetStakingPoolCountError", err)
		return nil, err
	}

	totalPools := poolCount.Int64()
	addresses := []common.Address{}
	limit := big.NewInt(poolPrelaunchBatchSize)
	for i := int64(0); i < totalPools; i += poolPrelaunchBatchSize {
		// Get a batch of addresses
		offset := big.NewInt(i)
		newAddresses, err := s.contract.StafiStakingPoolContract.GetPrelaunchStakingpools(s.conn.callOpts, offset, limit)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, newAddresses...)
	}
	return addresses, nil
}
