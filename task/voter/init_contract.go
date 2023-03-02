package task_voter

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	distributor "github.com/stafiprotocol/reth/bindings/Distributor"
	light_node "github.com/stafiprotocol/reth/bindings/LightNode"
	network_balances "github.com/stafiprotocol/reth/bindings/NetworkBalances"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	super_node "github.com/stafiprotocol/reth/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/reth/bindings/UserDeposit"
	withdraw "github.com/stafiprotocol/reth/bindings/Withdraw"
	"github.com/stafiprotocol/reth/pkg/utils"
)

func (task *Task) initContract() error {

	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	task.storageContract = storageContract

	lightNodeAddress, err := task.getContractAddress(storageContract, "stafiLightNode")
	if err != nil {
		return err
	}
	task.lightNodeContract, err = light_node.NewLightNode(lightNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	superNodeAddress, err := task.getContractAddress(storageContract, "stafiSuperNode")
	if err != nil {
		return err
	}
	task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	if task.enableDistribute {
		stafiDistributorAddress, err := task.getContractAddress(storageContract, "stafiDistributor")
		if err != nil {
			return err
		}
		task.distributorContract, err = distributor.NewDistributor(stafiDistributorAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}

		stafiFeePoolAddress, err := task.getContractAddress(storageContract, "stafiFeePool")
		if err != nil {
			return err
		}
		task.feePoolAddress = stafiFeePoolAddress

		stafiSuperNodeFeePoolAddress, err := task.getContractAddress(storageContract, "stafiSuperNodeFeePool")
		if err != nil {
			return err
		}
		task.superNodeFeePoolAddress = stafiSuperNodeFeePoolAddress
	}

	networkBalancesAddress, err := task.getContractAddress(storageContract, "stafiNetworkBalances")
	if err != nil {
		return err
	}
	task.networkBalancesContract, err = network_balances.NewNetworkBalances(networkBalancesAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	rethAddress, err := task.getContractAddress(storageContract, "rETHToken")
	if err != nil {
		return err
	}
	task.rethContract, err = reth.NewReth(rethAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	userDepositAddress, err := task.getContractAddress(storageContract, "stafiUserDeposit")
	if err != nil {
		return err
	}
	task.userDepositContract, err = user_deposit.NewUserDeposit(userDepositAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	networkSettingsAddress, err := task.getContractAddress(storageContract, "stafiNetworkSettings")
	if err != nil {
		return err
	}
	task.networkSettingsContract, err = network_settings.NewNetworkSettings(networkSettingsAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	withdrawAddress, err := task.getContractAddress(storageContract, "stafiWithdraw")
	if err != nil {
		return err
	}
	task.withdrawContract, err = withdraw.NewWithdraw(withdrawAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) getContractAddress(storage *storage.Storage, name string) (common.Address, error) {
	address, err := storage.GetAddress(task.connection.CallOpts(nil), utils.ContractStorageKey(name))
	if err != nil {
		return common.Address{}, err
	}
	if bytes.Equal(address.Bytes(), common.Address{}.Bytes()) {
		return common.Address{}, fmt.Errorf("address empty")
	}
	return address, nil
}

func (task *Task) NodeVoted(storage *storage.Storage, sender common.Address, _block *big.Int, _totalEth *big.Int, _stakingEth *big.Int, _rethSupply *big.Int) (bool, error) {
	return storage.GetBool(task.connection.CallOpts(nil), utils.NodeSubmissionKey(sender, _block, _totalEth, _stakingEth, _rethSupply))
}
