package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
)

// 1 deposited { 2 withdrawl match 3 staked 4 withdrawl unmatch } { 5 offboard 6 OffBoard can withdraw 7 OffBoard withdrawed } 8 waiting 9 active 10 exited 11 withdrawable 12 withdrawdone { 13 distributed }
// 51 active+slash 52 exit+slash 53 withdrawable+slash 54 withdrawdone+slash 55 distributed+slash
const (
	ValidatorStatusDeposited = uint8(1)

	// lightnode + super node related
	ValidatorStatusWithdrawMatch   = uint8(2)
	ValidatorStatusStaked          = uint8(3)
	ValidatorStatusWithdrawUnmatch = uint8(4)

	// lightnode related
	ValidatorStatusOffBoard            = uint8(5)
	ValidatorStatusOffBoardCanWithdraw = uint8(6)
	ValidatorStatusOffBoardWithdrawed  = uint8(7)

	// status on beacon chain
	ValidatorStatusWaiting      = uint8(8)
	ValidatorStatusActive       = uint8(9)
	ValidatorStatusExited       = uint8(10)
	ValidatorStatusWithdrawable = uint8(11)
	ValidatorStatusWithdrawDone = uint8(12)

	// after distribute reward
	ValidatorStatusDistributed = uint8(13) // distribute full withdrawal

	// after slash
	ValidatorStatusActiveSlash       = uint8(51)
	ValidatorStatusExitedSlash       = uint8(52)
	ValidatorStatusWithdrawableSlash = uint8(53)
	ValidatorStatusWithdrawDoneSlash = uint8(54)

	ValidatorStatusDistributedSlash = uint8(55) // distribute full withdrawal
)

// 1 common node 2 trust node 3 light node 4 super node
const (
	NodeTypeCommon = uint8(1)
	NodeTypeTrust  = uint8(2)
	NodeTypeLight  = uint8(3)
	NodeTypeSuper  = uint8(4)
)

const (
	ValidatorEverSlashedFalse = uint8(0)
	ValidatorEverSlashedTrue  = uint8(1)
)

const (
	FeePool          = uint8(1)
	SuperNodeFeePool = uint8(2)
)
const (
	V1  = "v1"
	V2  = "v2"
	Dev = "dev"
)

const (
	MetaTypeEth1BlockSyncer            = uint8(1) // dealed block height
	MetaTypeEth2ValidatorInfoSyncer    = uint8(2) // dealed epoch
	MetaTypeEth2ValidatorBalanceSyncer = uint8(3) // dealed epoch
	MetaTypeV1ValidatorSyncer          = uint8(4) // dealed block height
	MetaTypeEth2BlockSyncer            = uint8(5) // dealed epoch
	MetaTypeEth2NodeBalanceCollector   = uint8(6) // dealed epoch
)
const (
	SlashTypeFeeRecipient  = uint8(1)
	SlashTypeProposerSlash = uint8(2)
	SlashTypeAttesterSlash = uint8(3)
	SlashTypeSyncMiss      = uint8(4)
	SlashTypeAttesterMiss  = uint8(5)
	SlashTypeProposerMiss  = uint8(6)
)

const (
	V1EndEpoch          = uint64(148000)
	Eth1StartHeight     = uint64(15572967)
	TheMergeEpoch       = uint64(148896)
	RewardV1EndEpoch    = uint64(1) //todo mainnet
	RewardEpochInterval = uint64(75)

	StandardEffectiveBalance            = uint64(32e9) //gwei
	StandardSuperNodeFakeDepositBalance = uint64(1e9)  //gwei
	OfficialSlashAmount                 = uint64(1e9)  //gwei
	StandardLightNodeDepositAmount      = uint64(12e9) //gwei
	MaxPartialWithdrawalAmount          = uint64(8e9)  //gwei
)

//	enum ClaimType {
//	    None,
//	    CLAIMREWARD,
//	    CLAIMDEPOSIT,
//	    CLAIMTOTAL
//	}
const (
	NodeClaimTypeNone         = uint8(0)
	NodeClaimTypeClaimReward  = uint8(1)
	NodeClaimTypeClaimDeposit = uint8(2)
	NodeClaimTypeClaimTotal   = uint8(3)
)

var (
	// dev use
	OldRethSupply, _ = new(big.Int).SetString("25642334000000000000", 10)

	GweiDeci = decimal.NewFromInt(1e9)

	PlatformFeeV1Deci = decimal.NewFromFloat(0.1)
	NodeFeeV1Deci     = decimal.NewFromFloat(0.1)

	Percent5Deci  = decimal.NewFromFloat(0.05)
	Percent90Deci = decimal.NewFromFloat(0.9)
)

const (
	StakerWithdrawalClaimableTimestamp = uint64(1)
	MinValidatorWithdrawabilityDelay   = uint64(256 + 5)
	MaxDistributeSecondsInterval       = uint64(8 * 60 * 60)
	MaxDistributeEpochInterval         = uint64(75)

	EjectorUptimeInterval = uint64(10 * 60)
)

// Get an eth2 epoch number by time
func EpochAtTimestamp(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
}

func StartTimestampOfEpoch(config beacon.Eth2Config, epoch uint64) uint64 {
	return (epoch-config.GenesisEpoch)*config.SecondsPerEpoch + config.GenesisTime
}

func TimestampOfSlot(config beacon.Eth2Config, slot uint64) uint64 {
	return slot*config.SecondsPerSlot + config.GenesisTime
}

// Get an eth2 first slot number by epoch
func StartSlotOfEpoch(config beacon.Eth2Config, epoch uint64) uint64 {
	return config.SlotsPerEpoch * epoch
}

// func GetNodeReward(balance, effectiveBalance, nodeDepositAmount uint64) uint64 {
// 	if balance == 0 || effectiveBalance == 0 {
// 		return 0
// 	}
// 	reward := uint64(0)
// 	if balance > effectiveBalance {
// 		reward = balance - effectiveBalance
// 	}

// 	rewardDeci := decimal.NewFromInt(int64(reward)).Mul(decimal.NewFromFloat(0.9))
// 	nodeRewardDeci := decimal.NewFromInt(int64(nodeDepositAmount)).Mul(rewardDeci).Div(decimal.NewFromInt(int64(effectiveBalance)))
// 	stakerRawReard := rewardDeci.Sub(nodeRewardDeci)

// 	nodeRewardDeci = nodeRewardDeci.Add(stakerRawReard.Mul(decimal.NewFromFloat(0.1)))

// 	return nodeRewardDeci.BigInt().Uint64()
// }

func GetNodeManagedEth(nodeDeposit, balance uint64, status uint8) uint64 {
	switch status {
	case ValidatorStatusDeposited, ValidatorStatusWithdrawMatch, ValidatorStatusWithdrawUnmatch:
		return nodeDeposit

	case ValidatorStatusStaked, ValidatorStatusWaiting:
		return StandardEffectiveBalance

	default:
		return balance
	}
}

func GetGaspriceFromEthgasstation() (base, priority uint64, err error) {
	rsp, err := http.Get("https://api.ethgasstation.info/api/fee-estimate")
	if err != nil {
		return 0, 0, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("status err %d", rsp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, 0, err
	}
	if len(bodyBytes) == 0 {
		return 0, 0, fmt.Errorf("bodyBytes zero err")
	}
	resGasPrice := ResGasPrice{}
	err = json.Unmarshal(bodyBytes, &resGasPrice)
	if err != nil {
		return 0, 0, err
	}
	return uint64(resGasPrice.BaseFee), uint64(resGasPrice.PriorityFee.Fast), nil

}

type ResGasPrice struct {
	BaseFee     int     `json:"baseFee"`
	BlockNumber int     `json:"blockNumber"`
	BlockTime   float64 `json:"blockTime"`
	GasPrice    struct {
		Fast     int `json:"fast"`
		Instant  int `json:"instant"`
		Standard int `json:"standard"`
	} `json:"gasPrice"`
	NextBaseFee int `json:"nextBaseFee"`
	PriorityFee struct {
		Fast     int `json:"fast"`
		Instant  int `json:"instant"`
		Standard int `json:"standard"`
	} `json:"priorityFee"`
}

// statistic use
// return (user deposit, user reward, val deposit, val reward, paltform fee)
func GetUserValPlatformDepositAndRewardV1(validatorBalance, nodeDepositAmount uint64, platformFee, nodeFee decimal.Decimal) (uint64, uint64, uint64, uint64, uint64) {
	userDepositBalance := StandardEffectiveBalance - nodeDepositAmount

	switch {
	case validatorBalance == StandardEffectiveBalance:
		return userDepositBalance, 0, nodeDepositAmount, 0, 0
	case validatorBalance < StandardEffectiveBalance:
		loss := StandardEffectiveBalance - validatorBalance
		if loss < nodeDepositAmount {
			return userDepositBalance, 0, nodeDepositAmount - loss, 0, 0
		} else {
			return validatorBalance, 0, 0, 0, 0
		}
	case validatorBalance > StandardEffectiveBalance:
		// total staking reward
		reward := validatorBalance - StandardEffectiveBalance
		rewardDeci := decimal.NewFromInt(int64(reward))

		// platform Fee
		platformFeeDeci := rewardDeci.Mul(platformFee)

		// node+user stake reward
		nodeAndUserStakeRewardDeci := rewardDeci.Sub(platformFeeDeci)

		// user stake reward
		userStakeRewardDeci := nodeAndUserStakeRewardDeci.Mul(decimal.NewFromInt(int64(userDepositBalance))).Div(decimal.NewFromInt(int64(StandardEffectiveBalance)))
		// node stake reward
		nodeStakeRewardDeci := nodeAndUserStakeRewardDeci.Sub(userStakeRewardDeci)

		// node commisson reward from user
		nodeCommissionRewardFromUserDeci := userStakeRewardDeci.Mul(nodeFee)

		// user reward
		userRewardDeci := userStakeRewardDeci.Sub(nodeCommissionRewardFromUserDeci)
		// node reward
		nodeRewardDeci := nodeStakeRewardDeci.Add(nodeCommissionRewardFromUserDeci)

		return userDepositBalance, userRewardDeci.BigInt().Uint64(), nodeDepositAmount, nodeRewardDeci.BigInt().Uint64(), platformFeeDeci.BigInt().Uint64()
	default:
		// should not happen here
		panic("GetUserValPlatformDepositAndReward ")
	}
}

// v1: platform = 10% node = 90%*(nodedeposit/32)+90%*(1- nodedeposit/32)*10%  user = 90%*(1- nodedeposit/32)*90%
// return (user reward, node reward, paltform fee)
func GetUserNodePlatformRewardV1(nodeDepositBalance uint64, rewardDeci decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {

	if !rewardDeci.IsPositive() || nodeDepositBalance > StandardEffectiveBalance {
		return decimal.Zero, decimal.Zero, decimal.Zero
	}
	userDepositBalance := StandardEffectiveBalance - nodeDepositBalance

	// platform Fee
	platformFeeDeci := rewardDeci.Mul(PlatformFeeV1Deci) // 10%

	// node+user stake reward
	nodeAndUserStakeRewardDeci := rewardDeci.Sub(platformFeeDeci) // 90%

	// user stake reward
	userStakeRewardDeci := nodeAndUserStakeRewardDeci.Mul(decimal.NewFromInt(int64(userDepositBalance))).Div(decimal.NewFromInt(int64(StandardEffectiveBalance))) // 90%*(1 - nodedeposit/32)
	// node stake reward
	nodeStakeRewardDeci := nodeAndUserStakeRewardDeci.Sub(userStakeRewardDeci) // 90%*(nodedeposit/32)

	// node commisson reward from user
	nodeCommissionRewardFromUserDeci := userStakeRewardDeci.Mul(NodeFeeV1Deci) // 90%*(1 - nodedeposit/32)*10%

	// user reward
	userRewardDeci := userStakeRewardDeci.Sub(nodeCommissionRewardFromUserDeci) // 90%*(1 - nodedeposit/32)*90%
	// node reward
	nodeRewardDeci := nodeStakeRewardDeci.Add(nodeCommissionRewardFromUserDeci) // // 90%*(nodedeposit/32) + 90%*(1 - nodedeposit/32)*10%

	return userRewardDeci, nodeRewardDeci, platformFeeDeci

}

// v2: platform = 5%  node = 5% + (90% * nodedeposit/32) user = 90%*(1-nodedeposit/32)
// platform = 5%  node = 5% + (90% * nodedeposit/32)
// rewardDeci decimals maybe 9/18, also the returns
// return (user reward, node reward, paltform fee)
func GetUserNodePlatformRewardV2(nodeDepositAmount uint64, rewardDeci decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	if !rewardDeci.IsPositive() || nodeDepositAmount > StandardEffectiveBalance {
		return decimal.Zero, decimal.Zero, decimal.Zero
	}
	nodeDepositAmountDeci := decimal.NewFromInt(int64(nodeDepositAmount))
	standEffectiveBalanceDeci := decimal.NewFromInt(int64(StandardEffectiveBalance))

	// platform Fee
	platformFeeDeci := rewardDeci.Mul(Percent5Deci)
	nodeRewardDeci := platformFeeDeci.Add(rewardDeci.Mul(Percent90Deci).Mul(nodeDepositAmountDeci).Div(standEffectiveBalanceDeci))

	userRewardDeci := rewardDeci.Sub(platformFeeDeci).Sub(nodeRewardDeci)
	if userRewardDeci.IsNegative() {
		userRewardDeci = decimal.Zero
	}

	return userRewardDeci, nodeRewardDeci, platformFeeDeci

}

func GetValidatorTotalReward(balance, totalWithdrawal, totalFee uint64) uint64 {
	totalBalance := balance + totalWithdrawal + totalFee
	if totalBalance > StandardEffectiveBalance {
		return totalBalance - StandardEffectiveBalance
	}
	return 0
}

func ContractStorageKey(name string) [32]byte {
	// keccak256(abi.encodePacked("contract.address", _contractName))
	return crypto.Keccak256Hash([]byte("contract.address"), []byte(name))
}

func MerkleTreeDealedEpochStorageKey() [32]byte {
	return crypto.Keccak256Hash([]byte("stafiDistributor.merkleRoot.dealedEpoch"))
}

func NodeSubmissionKey(sender common.Address, _block *big.Int, _totalEth *big.Int, _stakingEth *big.Int, _rethSupply *big.Int) [32]byte {
	// keccak256(abi.encodePacked("network.balances.submitted.node", sender, _block, _totalEth, _stakingEth, _rethSupply))
	return crypto.Keccak256Hash([]byte("network.balances.submitted.node"), sender.Bytes(), common.LeftPadBytes(_block.Bytes(), 32),
		common.LeftPadBytes(_totalEth.Bytes(), 32), common.LeftPadBytes(_stakingEth.Bytes(), 32), common.LeftPadBytes(_rethSupply.Bytes(), 32))
}

func StafiWithdrawProposalNodeKey(sender common.Address, proposalId [32]byte) [32]byte {
	return crypto.Keccak256Hash([]byte("stafiWithdraw.proposal.node.key"), proposalId[:], sender.Bytes())
}

func DistributeWithdrawalsProposalNodeKey(sender common.Address, _dealedHeight, _userAmount, _nodeAmount, _platformAmount, _maxClaimableWithdrawIndex *big.Int) [32]byte {
	proposalId := crypto.Keccak256Hash([]byte("distributeWithdrawals"), common.LeftPadBytes(_dealedHeight.Bytes(), 32), common.LeftPadBytes(_userAmount.Bytes(), 32),
		common.LeftPadBytes(_nodeAmount.Bytes(), 32), common.LeftPadBytes(_platformAmount.Bytes(), 32), common.LeftPadBytes(_maxClaimableWithdrawIndex.Bytes(), 32))
	return StafiWithdrawProposalNodeKey(sender, proposalId)
}

func StafiDistributorProposalNodeKey(sender common.Address, proposalId [32]byte) [32]byte {
	return crypto.Keccak256Hash([]byte("stafiDistributor.proposal.node.key"), proposalId[:], sender.Bytes())
}

func ReserveEthForWithdrawProposalId(cycle *big.Int) [32]byte {
	return crypto.Keccak256Hash([]byte("reserveEthForWithdraw"), common.LeftPadBytes(cycle.Bytes(), 32))
}

func DistributeFeeProposalNodeKey(sender common.Address, _dealedHeight, _userAmount, _nodeAmount, _platformAmount *big.Int) [32]byte {
	proposalId := crypto.Keccak256Hash([]byte("distributeFee"), common.LeftPadBytes(_dealedHeight.Bytes(), 32), common.LeftPadBytes(_userAmount.Bytes(), 32),
		common.LeftPadBytes(_nodeAmount.Bytes(), 32), common.LeftPadBytes(_platformAmount.Bytes(), 32))
	return StafiDistributorProposalNodeKey(sender, proposalId)
}

func DistributeSuperNodeFeeProposalNodeKey(sender common.Address, _dealedHeight, _userAmount, _nodeAmount, _platformAmount *big.Int) [32]byte {
	proposalId := crypto.Keccak256Hash([]byte("distributeSuperNodeFee"), common.LeftPadBytes(_dealedHeight.Bytes(), 32), common.LeftPadBytes(_userAmount.Bytes(), 32),
		common.LeftPadBytes(_nodeAmount.Bytes(), 32), common.LeftPadBytes(_platformAmount.Bytes(), 32))
	return StafiDistributorProposalNodeKey(sender, proposalId)
}
