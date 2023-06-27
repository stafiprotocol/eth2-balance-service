package task_voter

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	distributor "github.com/stafiprotocol/eth2-balance-service/bindings/Distributor"
	light_node "github.com/stafiprotocol/eth2-balance-service/bindings/LightNode"
	network_balances "github.com/stafiprotocol/eth2-balance-service/bindings/NetworkBalances"
	reth "github.com/stafiprotocol/eth2-balance-service/bindings/Reth"
	"github.com/stafiprotocol/eth2-balance-service/bindings/Settings"
	stake_portal_rate "github.com/stafiprotocol/eth2-balance-service/bindings/StakePortalRate"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	super_node "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	withdraw "github.com/stafiprotocol/eth2-balance-service/bindings/Withdraw"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) initContract() error {

	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	task.storageContract = storageContract

	lightNodeAddress, err := utils.GetContractAddress(storageContract, "stafiLightNode")
	if err != nil {
		return err
	}
	task.lightNodeContract, err = light_node.NewLightNode(lightNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	superNodeAddress, err := utils.GetContractAddress(storageContract, "stafiSuperNode")
	if err != nil {
		return err
	}
	task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	stafiDistributorAddress, err := utils.GetContractAddress(storageContract, "stafiDistributor")
	if err != nil {
		return err
	}
	task.distributorContract, err = distributor.NewDistributor(stafiDistributorAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	stafiFeePoolAddress, err := utils.GetContractAddress(storageContract, "stafiFeePool")
	if err != nil {
		return err
	}
	task.feePoolAddress = stafiFeePoolAddress

	stafiSuperNodeFeePoolAddress, err := utils.GetContractAddress(storageContract, "stafiSuperNodeFeePool")
	if err != nil {
		return err
	}
	task.superNodeFeePoolAddress = stafiSuperNodeFeePoolAddress

	networkBalancesAddress, err := utils.GetContractAddress(storageContract, "stafiNetworkBalances")
	if err != nil {
		return err
	}
	task.networkBalancesContract, err = network_balances.NewNetworkBalances(networkBalancesAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	rethAddress, err := utils.GetContractAddress(storageContract, "rETHToken")
	if err != nil {
		return err
	}
	task.rethContract, err = reth.NewReth(rethAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	userDepositAddress, err := utils.GetContractAddress(storageContract, "stafiUserDeposit")
	if err != nil {
		return err
	}
	task.userDepositContract, err = user_deposit.NewUserDeposit(userDepositAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	networkSettingsAddress, err := utils.GetContractAddress(storageContract, "stafiNetworkSettings")
	if err != nil {
		return err
	}
	task.networkSettingsContract, err = network_settings.NewNetworkSettings(networkSettingsAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	withdrawAddress, err := utils.GetContractAddress(storageContract, "stafiWithdraw")
	if err != nil {
		return err
	}
	task.withdrawContract, err = withdraw.NewWithdraw(withdrawAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	//arbitrum
	task.arbitrumStakePortalRateContract, err = stake_portal_rate.NewStakePortalRate(task.arbitrumStakePortalRateAddress, task.arbitrumConn.Eth1Client())
	if err != nil {
		return err
	}

	return nil
}

// balance network related
func (task *Task) NodeVotedBalanceSubmission(storage *storage.Storage, sender common.Address, _block *big.Int, _totalEth *big.Int, _stakingEth *big.Int, _rethSupply *big.Int) (bool, error) {
	return storage.GetBool(task.connection.CallOpts(nil), utils.NodeSubmissionKey(sender, _block, _totalEth, _stakingEth, _rethSupply))
}

// withdraw pool related
func (task *Task) NodeVotedDistributeWithdrawals(storage *storage.Storage, sender common.Address, _dealedHeight, _userAmount, _nodeAmount, _platformAmount, _maxClaimableWithdrawIndex *big.Int) (bool, error) {
	return storage.GetBool(task.connection.CallOpts(nil), utils.DistributeWithdrawalsProposalNodeKey(sender, _dealedHeight, _userAmount, _nodeAmount, _platformAmount, _maxClaimableWithdrawIndex))
}

// distributor related
func (task *Task) NodeVotedDistributeFeePool(storage *storage.Storage, sender common.Address, _dealedHeight, _userAmount, _nodeAmount, _platformAmount *big.Int) (bool, error) {
	return storage.GetBool(task.connection.CallOpts(nil), utils.DistributeFeeProposalNodeKey(sender, _dealedHeight, _userAmount, _nodeAmount, _platformAmount))
}

func (task *Task) NodeVotedDistributeSuperNodeFeePool(storage *storage.Storage, sender common.Address, _dealedHeight, _userAmount, _nodeAmount, _platformAmount *big.Int) (bool, error) {
	return storage.GetBool(task.connection.CallOpts(nil), utils.DistributeSuperNodeFeeProposalNodeKey(sender, _dealedHeight, _userAmount, _nodeAmount, _platformAmount))
}
