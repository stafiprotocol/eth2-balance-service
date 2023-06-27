package task_ssv

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	super_node "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection/beacon"
	"github.com/stafiprotocol/eth2-balance-service/pkg/constants"
	"github.com/stafiprotocol/eth2-balance-service/pkg/credential"
	"github.com/stafiprotocol/eth2-balance-service/pkg/crypto/bls"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// only support super node account now !!!
// 0. find next key index and cache validator status on start
// 1. check stakepool balance, deposit if match
type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1StartHeight        uint64
	eth1Endpoint           string
	eth2Endpoint           string
	keyPair                *secp256k1.Keypair
	gasLimit               *big.Int
	maxGasPrice            *big.Int
	storageContractAddress common.Address
	seed                   []byte

	// --- need init on start
	dev                  bool
	chain                constants.Chain
	connection           *connection.Connection
	eth1WithdrawalAdress common.Address
	superNodeContract    *super_node.SuperNode
	userDepositContract  *user_deposit.UserDeposit
	nextKeyIndex         int
	validators           map[int]*Validator

	eth2Config beacon.Eth2Config
}

type Validator struct {
	privateKey *bls.PrivateKey
	status     uint8
	keyIndex   int
}

func NewTask(cfg *config.Config, seed []byte, keyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address fmt err")
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
		taskTicker:   15,
		stop:         make(chan struct{}),
		eth1Endpoint: cfg.Eth1Endpoint,
		eth2Endpoint: cfg.Eth2Endpoint,
		keyPair:      keyPair,
		seed:         seed,
		gasLimit:     gasLimitDeci.BigInt(),
		maxGasPrice:  maxGasPriceDeci.BigInt(),

		eth1StartHeight:        utils.TheMergeBlockNumber,
		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
	}

	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = connection.NewConnection(task.eth1Endpoint, task.eth2Endpoint, nil, nil, nil)
	if err != nil {
		return err
	}
	chainId, err := task.connection.Eth1Client().ChainID(context.Background())
	if err != nil {
		return err
	}
	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	switch chainId.Uint64() {
	case 1:
		task.dev = false
		task.chain = constants.GetChain(constants.ChainMAINNET)
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.MainnetConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}

	case 1337803: //zhejiang
		task.dev = true
		task.chain = constants.GetChain(constants.ChainZHEJIANG)
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, []byte{0x00, 0x00, 0x00, 0x69}) {
			return fmt.Errorf("endpoint network not match")
		}
	case 11155111: // sepolia
		task.dev = true
		task.chain = constants.GetChain(constants.ChainSEPOLIA)
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.SepoliaConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}

	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}
	if err != nil {
		return err
	}

	err = task.initContract()
	if err != nil {
		return err
	}

	err = task.findNextKeyIndex()
	if err != nil {
		return err
	}

	logrus.Infof("nextKeyIndex: %d", task.nextKeyIndex)

	utils.SafeGoWithRestart(task.handler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) copySeed() []byte {
	copyBts := make([]byte, len(task.seed))
	copy(copyBts, task.seed)
	return copyBts
}

func (task *Task) initContract() error {
	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	stafiWithdrawAddress, err := utils.GetContractAddress(storageContract, "stafiWithdraw")
	if err != nil {
		return err
	}

	task.eth1WithdrawalAdress = stafiWithdrawAddress

	superNodeAddress, err := utils.GetContractAddress(storageContract, "stafiSuperNode")
	if err != nil {
		return err
	}
	task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.connection.Eth1Client())
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
	return nil
}

func (task *Task) findNextKeyIndex() error {
	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("findNextKeyIndex reach retry limit")
		}
		credential, err := credential.NewCredential(task.copySeed(), task.nextKeyIndex, nil, constants.Chain{}, task.eth1WithdrawalAdress)
		if err != nil {
			return err
		}
		pubkey := credential.SigningPK().Marshal()
		pubkeyStatus, err := task.superNodeContract.GetSuperNodePubkeyStatus(nil, pubkey)
		if err != nil {
			logrus.Warnf("GetSuperNodePubkeyStatus err: %s", err.Error())
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}

		if pubkeyStatus.Uint64() == 0 {
			break
		}

		task.validators[task.nextKeyIndex] = &Validator{
			privateKey: credential.SigningSk,
			status:     uint8(pubkeyStatus.Uint64()),
			keyIndex:   task.nextKeyIndex,
		}

		task.nextKeyIndex++
	}
	return nil
}

func (task *Task) handler() {
	logrus.Info("start handler")
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
			logrus.Debug("checkAndStake start -----------")
			err := task.checkAndStake()
			if err != nil {
				logrus.Warnf("checkAndStake err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndStake end -----------")

			logrus.Debug("checkAndDeposit start -----------")
			err = task.checkAndDeposit()
			if err != nil {
				logrus.Warnf("checkAndDeposit err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndDeposit end -----------")
		}
	}
}
