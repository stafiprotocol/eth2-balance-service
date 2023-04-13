package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	distributor "github.com/stafiprotocol/eth2-balance-service/bindings/Distributor"
	light_node "github.com/stafiprotocol/eth2-balance-service/bindings/LightNode"
	network_balances "github.com/stafiprotocol/eth2-balance-service/bindings/NetworkBalances"
	reth "github.com/stafiprotocol/eth2-balance-service/bindings/Reth"
	"github.com/stafiprotocol/eth2-balance-service/bindings/Settings"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	super_node "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	withdraw "github.com/stafiprotocol/eth2-balance-service/bindings/Withdraw"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
)

type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1Endpoint           string
	eth2Endpoint           string
	keyPair                *secp256k1.Keypair
	gasLimit               *big.Int
	maxGasPrice            *big.Int
	storageContractAddress common.Address
	rewardEpochInterval    uint64
	rewardV1EndEpoch       uint64
	slashStartEpoch        uint64

	distributeFeeInitDealedHeight          int64 // dealedHeight is zero after upgrade new contract
	distributeSuperNodeFeeInitDealedHeight int64
	distributeWithdrawalInitDealedHeight   int64
	dev                                    bool

	// need init on start()
	connection           *shared.Connection
	db                   *db.WrapDb
	withdrawCredientials string

	storageContract         *storage.Storage
	userDepositContract     *user_deposit.UserDeposit
	lightNodeContract       *light_node.LightNode
	superNodeContract       *super_node.SuperNode
	networkSettingsContract *network_settings.NetworkSettings
	networkBalancesContract *network_balances.NetworkBalances
	rethContract            *reth.Reth
	distributorContract     *distributor.Distributor
	withdrawContract        *withdraw.Withdraw

	feePoolAddress          common.Address
	superNodeFeePoolAddress common.Address

	eth2Config beacon.Eth2Config

	domain []byte // for eth2 sigs
}

func NewTask(cfg *config.Config, dao *db.WrapDb, keyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address err")
	}

	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}

	if gasLimitDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("gas limit is zero")
	}
	maxGasPriceDeci, err := decimal.NewFromString(cfg.MaxGasPrice)
	if err != nil {
		return nil, err
	}
	if maxGasPriceDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("max gas price is zero")
	}

	s := &Task{
		taskTicker:             15,
		stop:                   make(chan struct{}),
		db:                     dao,
		keyPair:                keyPair,
		eth1Endpoint:           cfg.Eth1Endpoint,
		eth2Endpoint:           cfg.Eth2Endpoint,
		gasLimit:               gasLimitDeci.BigInt(),
		maxGasPrice:            maxGasPriceDeci.BigInt(),
		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),

		rewardEpochInterval: utils.RewardEpochInterval,
		rewardV1EndEpoch:    utils.RewardV1EndEpoch,
		slashStartEpoch:     utils.SlashStartEpoch,

		distributeFeeInitDealedHeight:          16638921,
		distributeSuperNodeFeeInitDealedHeight: 17024852,
		distributeWithdrawalInitDealedHeight:   1,
	}

	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair, task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}

	// set domain by chainId
	chainId, err := task.connection.Eth1Client().ChainID(context.Background())
	if err != nil {
		return err
	}
	switch chainId.Uint64() {
	case 1:
		domain, err := signing.ComputeDomain(
			params.MainnetConfig().DomainDeposit,
			params.MainnetConfig().GenesisForkVersion,
			params.MainnetConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		task.domain = domain
		task.dev = false
	case 1337803: //zhejiang
		domain, err := signing.ComputeDomain(
			params.PraterConfig().DomainDeposit,
			[]byte{0x00, 0x00, 0x00, 0x69},
			params.PraterConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		task.domain = domain
		task.dev = true
	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}

	if len(task.domain) == 0 {
		return fmt.Errorf("domain not ok")
	}

	if task.dev {
		task.rewardV1EndEpoch = utils.DevRewardV1EndEpoch
		task.slashStartEpoch = utils.DevTheMergeEpoch

		task.distributeFeeInitDealedHeight = 1
		task.distributeSuperNodeFeeInitDealedHeight = 1
		task.distributeWithdrawalInitDealedHeight = 1
	}

	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	err = task.initContract()
	if err != nil {
		return err
	}

	credentials, err := task.networkSettingsContract.GetWithdrawalCredentials(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	task.withdrawCredientials = hexutil.Encode(credentials)

	utils.SafeGoWithRestart(task.voteHandler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) voteHandler() {
	logrus.Info("start voteHandler")
	ticker := time.NewTicker(time.Duration(task.taskTicker) * time.Second)
	defer ticker.Stop()
	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return
		}

		select {
		case <-task.stop:
			logrus.Info("task has stopped")
			return
		case <-ticker.C:
			logrus.Debug("voteWithdrawalCredential start -----------")

			err := task.voteWithdrawalCredential()
			if err != nil {
				logrus.Warnf("vote voteWithdrawalCredential err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("voteWithdrawalCredential end -----------\n")

			logrus.Debug("distributeFee start -----------")
			err = task.distributeFee()
			if err != nil {
				logrus.Warnf("distributeFee err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("distributeFee end -----------\n")

			logrus.Debug("voteRate start -----------")
			err = task.voteRate()
			if err != nil {
				logrus.Warnf("vote rate err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}

			logrus.Debug("voteRate end -----------\n")

			logrus.Debug("setMerkleTree start -----------")
			err = task.setMerkleTree()
			if err != nil {
				logrus.Warnf("setMerkleTree err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("setMerkleTree end -----------\n")

			logrus.Debug("notifyValidatorExit start -----------")
			err = task.notifyValidatorExit()
			if err != nil {
				logrus.Warnf("notifyValidatorExit err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("notifyValidatorExit end -----------\n")
			retry = 0
		}
	}
}

func (task Task) getEpochStartBlocknumber(epoch uint64) (uint64, error) {
	eth2ValidatorBalanceSyncerStartSlot := utils.StartSlotOfEpoch(task.eth2Config, epoch)
	blocknumber := uint64(0)
	retry := 0
	for {
		if retry > 5 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := task.connection.GetBeaconBlock(eth2ValidatorBalanceSyncerStartSlot)
		if err != nil {
			return 0, err
		}
		// we will use next slot if not exist
		if !exist {
			eth2ValidatorBalanceSyncerStartSlot++
			retry++
			continue
		}
		if targetBeaconBlock.ExecutionBlockNumber == 0 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}
		blocknumber = targetBeaconBlock.ExecutionBlockNumber
		break
	}
	return blocknumber, nil
}

func (task *Task) waitTxOk(txHash common.Hash) error {
	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("waitTxOk reach retry limit")
		}
		tx, pending, err := task.connection.Eth1Client().TransactionByHash(context.Background(), txHash)
		if err == nil && !pending {
			receipt, err := task.connection.Eth1Client().TransactionReceipt(context.Background(), txHash)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx TransactionReceipt")
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			if receipt.Status == 1 { //success
				break
			} else { //failed
				utils.ShutdownRequestChannel <- struct{}{}
				return fmt.Errorf("tx %s failed", txHash)
			}
		} else {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx status")
			} else {
				logrus.WithFields(logrus.Fields{
					"hash":   tx.Hash(),
					"status": "pending",
				}).Warn("tx status")
			}
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}

	}
	logrus.WithFields(logrus.Fields{
		"tx": txHash.String(),
	}).Info("tx send ok")
	return nil
}
