package task_ssv

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

const (
	// todo: automatically detect event fetch limit from target rpc
	fetchEventBlockLimit      = uint64(4900)
	fetchEth1WaitBlockNumbers = uint64(5)
)

func (task *Task) updateOffchainState() error {
	latestBlockNumber, err := task.connectionOfSuperNodeAccount.Eth1LatestBlock()
	if err != nil {
		return err
	}

	if latestBlockNumber > fetchEth1WaitBlockNumbers {
		latestBlockNumber -= fetchEth1WaitBlockNumbers
	}

	logrus.Debugf("latestBlockNumber: %d, dealedBlockNumber: %d", latestBlockNumber, task.dealedEth1Block)
	if latestBlockNumber <= uint64(task.dealedEth1Block) {
		return nil
	}

	start := uint64(task.dealedEth1Block + 1)
	end := latestBlockNumber
	maxBlock := uint64(0)

	increaseNonce := func() {
		task.latestRegistrationNonce = new(big.Int).Add(task.latestRegistrationNonce, big.NewInt(1))
	}

	// 'ClusterDeposited',
	// 'ClusterWithdrawn',
	// 'ValidatorRemoved',
	// 'ValidatorAdded',
	// 'ClusterLiquidated',
	// 'ClusterReactivated',
	for i := start; i <= end; i += fetchEventBlockLimit {
		subStart := i
		subEnd := i + fetchEventBlockLimit - 1
		if end < i+fetchEventBlockLimit {
			subEnd = end
		}

		clusterDepositedIter, err := task.ssvClustersContract.FilterClusterDeposited(&bind.FilterOpts{
			Start:   subStart,
			End:     &subEnd,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterDepositedIter.Next() {
			increaseNonce()
			if clusterDepositedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterDepositedIter.Event.Cluster
				maxBlock = clusterDepositedIter.Event.Raw.BlockNumber
			}
		}

		clusterWithdrawnIter, err := task.ssvClustersContract.FilterClusterWithdrawn(&bind.FilterOpts{
			Start:   subStart,
			End:     &subEnd,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterWithdrawnIter.Next() {
			increaseNonce()
			if clusterWithdrawnIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterWithdrawnIter.Event.Cluster
				maxBlock = clusterWithdrawnIter.Event.Raw.BlockNumber
			}
		}

		validatorRemovedIter, err := task.ssvClustersContract.FilterValidatorRemoved(&bind.FilterOpts{
			Start:   subStart,
			End:     &subEnd,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for validatorRemovedIter.Next() {
			increaseNonce()
			if validatorRemovedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &validatorRemovedIter.Event.Cluster
				maxBlock = validatorRemovedIter.Event.Raw.BlockNumber
			}
		}

		validatorAdddedIter, err := task.ssvClustersContract.FilterValidatorAdded(&bind.FilterOpts{
			Start:   subStart,
			End:     &subEnd,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for validatorAdddedIter.Next() {
			increaseNonce()
			if validatorAdddedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &validatorAdddedIter.Event.Cluster
				maxBlock = validatorAdddedIter.Event.Raw.BlockNumber
			}
		}

		clusterLiquidatedIter, err := task.ssvClustersContract.FilterClusterLiquidated(&bind.FilterOpts{
			Start:   subStart,
			End:     &subEnd,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterLiquidatedIter.Next() {
			increaseNonce()
			if clusterLiquidatedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterLiquidatedIter.Event.Cluster
				maxBlock = clusterLiquidatedIter.Event.Raw.BlockNumber
			}
		}

		clusterReactivatedIter, err := task.ssvClustersContract.FilterClusterReactivated(&bind.FilterOpts{
			Start:   subStart,
			End:     &subEnd,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterReactivatedIter.Next() {
			increaseNonce()
			if clusterReactivatedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterReactivatedIter.Event.Cluster
				maxBlock = clusterReactivatedIter.Event.Raw.BlockNumber
			}
		}

		task.dealedEth1Block = subEnd

		logrus.WithFields(logrus.Fields{
			"start": subStart,
			"end":   subEnd,
		}).Debug("already dealed blocks")

	}
	return nil
}
