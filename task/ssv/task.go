package task_ssv

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection/beacon"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// sync deposit/stake events and pool latest info from execute chain
// sync validator latest info and epoch balance from consensus chain
// sync beacon block info from consensus chain
// sort by head: 0 eth1 syncer -> 1 latestInfo syncer -> 2 eth2Block syncer -> 3 valBalance syncer -> 4 nodeBalance collector -> 5 merkle tree
type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1StartHeight        uint64
	eth1Endpoint           string
	eth2Endpoint           string
	storageContractAddress common.Address
	rewardEpochInterval    uint64 //75
	calMerkleTreeDu        uint64 //75
	rewardV1EndEpoch       uint64
	slashStartEpoch        uint64
	dev                    bool

	eth2BlockStartEpoch uint64 // for eth2 block syncer

	// --- need init on start
	connection *connection.Connection

	eth2Config beacon.Eth2Config
}

func NewTask(cfg *config.Config) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address fmt err")
	}

	s := &Task{
		taskTicker:   15,
		stop:         make(chan struct{}),
		eth1Endpoint: cfg.Eth1Endpoint,
		eth2Endpoint: cfg.Eth2Endpoint,

		eth1StartHeight:     utils.TheMergeBlockNumber,
		eth2BlockStartEpoch: utils.TheMergeEpoch,
		rewardV1EndEpoch:    utils.RewardV1EndEpoch,
		slashStartEpoch:     utils.SlashStartEpoch,

		rewardEpochInterval:    utils.RewardEpochInterval,
		calMerkleTreeDu:        utils.RewardEpochInterval * 3, // 24 h
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
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.MainnetConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
	case 1337803: //zhejiang
		task.dev = true
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, []byte{0x00, 0x00, 0x00, 0x69}) {
			return fmt.Errorf("endpoint network not match")
		}
	case 11155111: // sepolia
		task.dev = true
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.SepoliaConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}

	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}
