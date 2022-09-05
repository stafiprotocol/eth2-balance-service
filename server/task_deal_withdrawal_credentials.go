package server

// import (
// 	"math/big"
// 	"strings"
// 	"time"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/stafiprotocol/reth/bindings/StakingPool"
// 	"github.com/stafiprotocol/reth/types"
// 	"github.com/stafiprotocol/reth/pkg/utils"
// )

// var poolPrelaunchBatchSize = int64(50)

// func (s *Service) dealWithdrawalCredentials() error {
// 	var retry = 0
// 	for {
// 		select {
// 		case <-s.stop:
// 			return nil
// 		default:
// 			// No more retries
// 			if retry > BlockRetryLimit {
// 				utils.ShutdownRequestChannel <- struct{}{}
// 				return nil
// 			}
// 			err := s.checkAndVoteCredentials()
// 			if err != nil {
// 				retry--
// 				time.Sleep(BlockRetryInterval)
// 				continue
// 			}

// 			time.Sleep(BlockRetryInterval)
// 			retry = BlockRetryLimit
// 		}
// 	}
// }

// func (s *Service) checkAndVoteCredentials() error {
// 	poolAddresses, err := s.getPrelaunchStakingpools()
// 	if err != nil {
// 		return err
// 	}
// 	if len(poolAddresses) == 0 {
// 		return nil
// 	}

// 	expectCredentials, err := s.contracts.Settings.GetWithdrawalCredentials(s.conn.callOpts)
// 	if err != nil {
// 		return err
// 	}

// 	for _, poolAddress := range poolAddresses {
// 		stakingPoolContract, err := StakingPool.NewStakingPool(poolAddress, s.conn.Client())
// 		if err != nil {
// 			return err
// 		}

// 		match, err := stakingPoolContract.GetWithdrawalCredentialsMatch(s.conn.callOpts)
// 		if err != nil {
// 			return err
// 		}
// 		if match {
// 			continue
// 		}
// 		pubkey, err := s.contracts.StafiStakingPoolManager.GetStakingPoolPubkey(s.conn.callOpts, poolAddress)
// 		if err != nil {
// 			return err
// 		}

// 		validatorStatus, err := s.conn.eth2Conn.GetValidatorStatus(types.BytesToValidatorPubkey(pubkey), nil)
// 		if err != nil {
// 			return err
// 		}

// 		if validatorStatus.Exists {
// 			if validatorStatus.WithdrawalCredentials == common.BytesToHash(expectCredentials) {
// 				tx, err := s.voteWithdrawCredentials(stakingPoolContract)
// 				if err != nil {
// 					return err
// 				}
// 				if tx != "" {
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// func (s *Service) voteWithdrawCredentials(stakingPoolContract *StakingPool.StakingPool) (string, error) {
// 	err := s.contracts.Conn.LockAndUpdateOpts()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer s.contracts.Conn.UnlockOpts()

// 	tx, err := stakingPoolContract.VoteWithdrawCredentials(s.conn.opts)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "Member has already voted to withdrawCredentials") {
// 			return "", nil
// 		} else {
// 			return "", err
// 		}
// 	}
// 	return tx.Hash().String(), nil
// }

// func (s *Service) getPrelaunchStakingpools() ([]common.Address, error) {
// 	poolCount, err := s.contracts.StafiStakingPoolManager.GetStakingPoolCount(s.conn.callOpts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPools := poolCount.Int64()
// 	addresses := []common.Address{}
// 	limit := big.NewInt(poolPrelaunchBatchSize)
// 	for i := int64(0); i < totalPools; i += poolPrelaunchBatchSize {
// 		// Get a batch of addresses
// 		offset := big.NewInt(i)
// 		newAddresses, err := s.contracts.StafiStakingPoolManager.GetPrelaunchStakingpools(s.conn.callOpts, offset, limit)
// 		if err != nil {
// 			return nil, err
// 		}
// 		addresses = append(addresses, newAddresses...)
// 	}
// 	return addresses, nil
// }
