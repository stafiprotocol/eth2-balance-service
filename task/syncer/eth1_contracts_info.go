package task_syncer

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) syncContractsInfo() error {
	poolInfo, err := dao_chaos.GetPoolInfo(task.db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// --- deposit pool
	poolBalance, err := task.userDepositContract.GetBalance(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	poolInfo.DepositPoolBalance = poolBalance.String()

	// --- reth
	rethSupply, err := task.rethContract.TotalSupply(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	if task.version == utils.Dev {
		rethSupply = new(big.Int).Sub(rethSupply, utils.OldRethSupply)
	}
	poolInfo.REthSupply = rethSupply.String()

	// --- withdraw pool
	latestDistributeWithdrawalHeight, err := task.withdrawContract.LatestDistributeHeight(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	nextWithdrawIndex, err := task.withdrawContract.NextWithdrawIndex(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	maxClaimableWithdrawIndex, err := task.withdrawContract.MaxClaimableWithdrawIndex(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	totalMissingAmountForWithdraw, err := task.withdrawContract.TotalMissingAmountForWithdraw(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	currentWithdrawCycle, err := task.withdrawContract.CurrentWithdrawCycle(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	totalWithdrawAmountCurrentCycle, err := task.withdrawContract.TotalWithdrawAmountAtCycle(task.connection.CallOpts(nil), currentWithdrawCycle)
	if err != nil {
		return err
	}
	withdrawLimitPerCycle, err := task.withdrawContract.WithdrawLimitPerCycle(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	poolInfo.LatestDistributeWithdrawalHeight = latestDistributeWithdrawalHeight.Uint64()
	poolInfo.NextWithdrawIndex = nextWithdrawIndex.Uint64()
	poolInfo.MaxClaimableWithdrawIndex = maxClaimableWithdrawIndex.Uint64()
	poolInfo.TotalMissingAmountForWithdraw = totalMissingAmountForWithdraw.String()
	poolInfo.TotalWithdrawAmountCurrentCycle = totalWithdrawAmountCurrentCycle.String()
	poolInfo.WithdrawLimitPerCycle = withdrawLimitPerCycle.String()

	// ---- distributor
	merkleTreeDealedEpochOnchain, err := task.MerkleTreeDealedEpoch(task.storageContract)
	if err != nil {
		return err
	}
	poolInfo.LatestMerkleTreeEpoch = merkleTreeDealedEpochOnchain.Uint64()

	// --- eth price
	price, err := task.getEthPrice()
	if err != nil {
		logrus.Warnf("get eth price err: %s", err)
	} else {
		poolInfo.EthPrice = price
	}

	// --- gas fee
	base, priority, err := task.getGasPrice()
	if err != nil {
		logrus.Warnf("get gas price err: %s", err)
	} else {
		if base != 0 {
			poolInfo.BaseFee = base
		}

		if priority != 0 {
			poolInfo.PriorityFee = priority
		}
	}

	logrus.WithFields(logrus.Fields{
		"depositPoolBalance": poolInfo.DepositPoolBalance,
		"rethsupply":         poolInfo.REthSupply,
	}).Debug("poolInfo")

	return dao_chaos.UpOrInPoolInfo(task.db, poolInfo)
}

func (task *Task) getEthPrice() (string, error) {
	priceMap, err := utils.GetPriceFromCoinGecko("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return "", err
	}

	if price, exist := priceMap[utils.SymbolEth]; exist {
		return decimal.NewFromFloat(price).Mul(decimal.NewFromInt(1e6)).StringFixed(0), nil
	}
	return "", fmt.Errorf("no eth price")
}

func (task *Task) getGasPrice() (base uint64, priority uint64, err error) {
	return utils.GetGaspriceFromEthgasstation()
}
