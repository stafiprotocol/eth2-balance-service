package task_ssv

import (
	// "bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/eth2-balance-service/bindings/SsvClusters"
	"github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetwork"
	"github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetworkViews"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	super_node "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection/beacon"
	"github.com/stafiprotocol/eth2-balance-service/pkg/constants"
	"github.com/stafiprotocol/eth2-balance-service/pkg/crypto/bls"
	"github.com/stafiprotocol/eth2-balance-service/pkg/keyshare"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

var (
	minAmountNeedStake   = decimal.NewFromBigInt(big.NewInt(31), 18)
	minAmountNeedDeposit = decimal.NewFromBigInt(big.NewInt(32), 18)

	superNodeDepositAmount = decimal.NewFromBigInt(big.NewInt(1), 18)
	superNodeStakeAmount   = decimal.NewFromBigInt(big.NewInt(31), 18)
)

// only support super node account now !!!
// 0. find next key index and cache validator status on start
// 1. update validator status periodically
// 2. check stakepool balance periodically, call stake/deposit if match
// 3. register validator on ssv, if status is staked on eth
// 3. register validator on ssv, if status is exited on eth
type Task struct {
	taskTicker                     int64
	stop                           chan struct{}
	eth1StartHeight                uint64
	eth1Endpoint                   string
	eth2Endpoint                   string
	superNodeKeyPair               *secp256k1.Keypair
	ssvKeyPair                     *secp256k1.Keypair
	gasLimit                       *big.Int
	maxGasPrice                    *big.Int
	clusterInitSsvAmount           *big.Int
	storageContractAddress         common.Address
	ssvNetworkContractAddress      common.Address
	ssvNetworkViewsContractAddress common.Address
	seed                           []byte

	// --- need init on start
	dev                          bool
	ssvApiNetwork                string
	chain                        constants.Chain
	connectionOfSuperNodeAccount *connection.Connection
	connectionOfSsvAccount       *connection.Connection
	eth1WithdrawalAdress         common.Address
	superNodeContract            *super_node.SuperNode
	userDepositContract          *user_deposit.UserDeposit
	ssvNetworkContract           *ssv_network.SsvNetwork
	ssvNetworkViewsContract      *ssv_network_views.SsvNetworkViews
	ssvClustersContract          *ssv_clusters.SsvClusters
	nextKeyIndex                 int
	dealedEth1Block              uint64
	validators                   map[int]*Validator

	latestCluster *ssv_clusters.ISSVNetworkCoreCluster

	operators []*keyshare.Operator

	eth2Config beacon.Eth2Config
}

type Validator struct {
	privateKey *bls.PrivateKey
	keyIndex   int
	status     uint8 // status: 0 uninitiated 1 deposited 2 staked 3 registed on ssv 4 exited on eth 5 removed on ssv
}

const (
	valStatusUnInitiated = uint8(0)
	valStatusDeposited   = uint8(1)
	valStatusMatch       = uint8(2)
	valStatusUnmatch     = uint8(3)
	valStatusStaked      = uint8(4)

	valStatusRegistedOnSsv = uint8(5)

	valStatusActiveOnBeacon = uint8(6)
	valStatusExitedOnBeacon = uint8(7)

	valStatusRemovedOnSsv = uint8(8)
)

func NewTask(cfg *config.Config, seed []byte, superNodeKeyPair, ssvKeyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("storage contract address fmt err")
	}
	if !common.IsHexAddress(cfg.Contracts.SsvNetworkAddress) {
		return nil, fmt.Errorf("ssvnetwork contract address fmt err")
	}
	if !common.IsHexAddress(cfg.Contracts.SsvNetworkViewsAddress) {
		return nil, fmt.Errorf("ssvnetworkviews contract address fmt err")
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

	clusterInitSsvAmount, err := decimal.NewFromString(cfg.ClusterInitSsvAmount)
	if err != nil {
		return nil, err
	}
	if clusterInitSsvAmount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("clusterInitSsvAmount is zero")
	}

	if len(cfg.Operators) == 1 || len(cfg.Operators)%3 != 1 {
		return nil, fmt.Errorf("operators length not match")
	}

	sort.Slice(cfg.Operators, func(i, j int) bool {
		return cfg.Operators[i] < cfg.Operators[j]
	})

	operaters := make([]*keyshare.Operator, len(cfg.Operators))
	for i := 0; i < len(cfg.Operators); i++ {
		operaters[i] = &keyshare.Operator{Id: int(cfg.Operators[i])}
	}

	s := &Task{
		taskTicker:           15,
		stop:                 make(chan struct{}),
		eth1Endpoint:         cfg.Eth1Endpoint,
		eth2Endpoint:         cfg.Eth2Endpoint,
		superNodeKeyPair:     superNodeKeyPair,
		ssvKeyPair:           ssvKeyPair,
		seed:                 seed,
		gasLimit:             gasLimitDeci.BigInt(),
		maxGasPrice:          maxGasPriceDeci.BigInt(),
		clusterInitSsvAmount: clusterInitSsvAmount.BigInt(),

		eth1StartHeight:                utils.TheMergeBlockNumber,
		storageContractAddress:         common.HexToAddress(cfg.Contracts.StorageContractAddress),
		ssvNetworkContractAddress:      common.HexToAddress(cfg.Contracts.SsvNetworkAddress),
		ssvNetworkViewsContractAddress: common.HexToAddress(cfg.Contracts.SsvNetworkViewsAddress),

		operators: operaters,

		validators: map[int]*Validator{},

		latestCluster: &ssv_clusters.ISSVNetworkCoreCluster{
			ValidatorCount:  0,
			NetworkFeeIndex: 0,
			Index:           0,
			Active:          true,
			Balance:         big.NewInt(0),
		},
	}

	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connectionOfSuperNodeAccount, err = connection.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.superNodeKeyPair, task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}
	task.connectionOfSsvAccount, err = connection.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.ssvKeyPair, task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}
	chainId, err := task.connectionOfSuperNodeAccount.Eth1Client().ChainID(context.Background())
	if err != nil {
		return err
	}

	// task.eth2Config, err = task.connectionOfSuperNodeAccount.Eth2Client().GetEth2Config()
	// if err != nil {
	// 	return err
	// }

	switch chainId.Uint64() {
	case 1: //mainnet
		task.dev = false
		task.chain = constants.GetChain(constants.ChainMAINNET)
		// if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.MainnetConfig().GenesisForkVersion) {
		// 	return fmt.Errorf("endpoint network not match")
		// }
		task.dealedEth1Block = 17705353
		task.ssvApiNetwork = "mainnet"

	case 11155111: // sepolia
		task.dev = true
		task.chain = constants.GetChain(constants.ChainSEPOLIA)
		// if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.SepoliaConfig().GenesisForkVersion) {
		// 	return fmt.Errorf("endpoint network not match")
		// }
		task.dealedEth1Block = 9354882
		task.ssvApiNetwork = "prater"
	case 5: // goerli
		task.dev = true
		task.chain = constants.GetChain(constants.ChainGOERLI)
		// if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.PraterConfig().GenesisForkVersion) {
		// 	return fmt.Errorf("endpoint network not match")
		// }
		task.dealedEth1Block = 9354882
		task.ssvApiNetwork = "prater"

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

	err = task.initNextKeyIndex()
	if err != nil {
		return err
	}
	logrus.Infof("nextKeyIndex: %d", task.nextKeyIndex)

	// init operators
	for _, op := range task.operators {
		operatorDetail, err := utils.GetOperatorDetail(task.ssvApiNetwork, op.Id)
		if err != nil {
			return err
		}
		op.PublicKey = operatorDetail.PublicKey
		feeDeci, err := decimal.NewFromString(operatorDetail.Fee)
		if err != nil {
			return err
		}
		op.Fee = feeDeci.BigInt().Uint64()
	}

	utils.SafeGo(task.handler)
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
	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connectionOfSuperNodeAccount.Eth1Client())
	if err != nil {
		return err
	}
	stafiWithdrawAddress, err := utils.GetContractAddress(storageContract, "stafiWithdraw")
	if err != nil {
		return err
	}

	task.eth1WithdrawalAdress = stafiWithdrawAddress

	logrus.Debugf("stafiWithdraw address: %s", task.eth1WithdrawalAdress.String())

	superNodeAddress, err := utils.GetContractAddress(storageContract, "stafiSuperNode")
	if err != nil {
		return err
	}
	task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.connectionOfSuperNodeAccount.Eth1Client())
	if err != nil {
		return err
	}

	userDepositAddress, err := utils.GetContractAddress(storageContract, "stafiUserDeposit")
	if err != nil {
		return err
	}
	task.userDepositContract, err = user_deposit.NewUserDeposit(userDepositAddress, task.connectionOfSuperNodeAccount.Eth1Client())
	if err != nil {
		return err
	}
	task.ssvNetworkContract, err = ssv_network.NewSsvNetwork(task.ssvNetworkContractAddress, task.connectionOfSuperNodeAccount.Eth1Client())
	if err != nil {
		return err
	}
	task.ssvNetworkViewsContract, err = ssv_network_views.NewSsvNetworkViews(task.ssvNetworkViewsContractAddress, task.connectionOfSuperNodeAccount.Eth1Client())
	if err != nil {
		return err
	}
	task.ssvClustersContract, err = ssv_clusters.NewSsvClusters(task.ssvNetworkContractAddress, task.connectionOfSuperNodeAccount.Eth1Client())
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) mustGetSuperNodePubkeyStatus(pubkey []byte) (uint8, error) {
	retry := 0
	var pubkeyStatus *big.Int
	var err error
	for {
		if retry > utils.RetryLimit {
			return 0, fmt.Errorf("updateValStatus reach retry limit")
		}
		pubkeyStatus, err = task.superNodeContract.GetSuperNodePubkeyStatus(nil, pubkey)
		if err != nil {
			logrus.Warnf("GetSuperNodePubkeyStatus err: %s", err.Error())
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
		break
	}

	return uint8(pubkeyStatus.Uint64()), nil
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
			logrus.Debug("checkAnddRepairNexKeyIndex start -----------")
			err := task.checkAnddRepairNexKeyIndex()
			if err != nil {
				logrus.Warnf("checkAnddRepairNexKeyIndex err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAnddRepairNexKeyIndex end -----------")

			logrus.Debug("updateValStatus start -----------")
			err = task.updateValStatus()
			if err != nil {
				logrus.Warnf("updateValStatus err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("updateValStatus end -----------")

			logrus.Debug("checkAndStake start -----------")
			err = task.checkAndStake()
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

			logrus.Debug("updateLatestClusters start -----------")
			err = task.updateLatestClusters()
			if err != nil {
				logrus.Warnf("updateLatestClusters err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("updateLatestClusters end -----------")

			logrus.Debug("checkAndRegisterOnSSV start -----------")
			err = task.checkAndRegisterOnSSV()
			if err != nil {
				logrus.Warnf("checkAndRegisterOnSSV err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndRegisterOnSSV end -----------")

			logrus.Debug("updateLatestClusters start -----------")
			err = task.updateLatestClusters()
			if err != nil {
				logrus.Warnf("updateLatestClusters err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("updateLatestClusters end -----------")

			logrus.Debug("checkAndRemoveOnSSV start -----------")
			err = task.checkAndRemoveOnSSV()
			if err != nil {
				logrus.Warnf("checkAndRemoveOnSSV err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndRemoveOnSSV end -----------")
		}
	}
}
