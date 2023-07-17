package task_ssv

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

const (
	fetchEventBlockLimit      = uint64(4900)
	fetchEth1WaitBlockNumbers = uint64(5)
)

func (task *Task) updateLatestClusters() error {
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
			Start:   start,
			End:     &end,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterDepositedIter.Next() {
			if clusterDepositedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterDepositedIter.Event.Cluster
				maxBlock = clusterDepositedIter.Event.Raw.BlockNumber
			}
		}

		clusterWithdrawnIter, err := task.ssvClustersContract.FilterClusterWithdrawn(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterWithdrawnIter.Next() {
			if clusterWithdrawnIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterWithdrawnIter.Event.Cluster
				maxBlock = clusterWithdrawnIter.Event.Raw.BlockNumber
			}
		}

		validatorRemovedIter, err := task.ssvClustersContract.FilterValidatorRemoved(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for validatorRemovedIter.Next() {
			if validatorRemovedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &validatorRemovedIter.Event.Cluster
				maxBlock = validatorRemovedIter.Event.Raw.BlockNumber
			}
		}

		validatorAdddedIter, err := task.ssvClustersContract.FilterValidatorAdded(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for validatorAdddedIter.Next() {
			if validatorAdddedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &validatorAdddedIter.Event.Cluster
				maxBlock = validatorAdddedIter.Event.Raw.BlockNumber
			}
		}

		clusterLiquidatedIter, err := task.ssvClustersContract.FilterClusterLiquidated(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterLiquidatedIter.Next() {
			if clusterLiquidatedIter.Event.Raw.BlockNumber > maxBlock {
				task.latestCluster = &clusterLiquidatedIter.Event.Cluster
				maxBlock = clusterLiquidatedIter.Event.Raw.BlockNumber
			}
		}

		clusterReactivatedIter, err := task.ssvClustersContract.FilterClusterReactivated(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: context.Background(),
		}, []common.Address{task.ssvKeyPair.CommonAddress()})

		if err != nil {
			return err
		}
		for clusterReactivatedIter.Next() {
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
