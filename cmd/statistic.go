// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/reth/bindings/FeePool"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	"github.com/stafiprotocol/reth/bindings/SuperNodeFeePool"
	user_deposit "github.com/stafiprotocol/reth/bindings/UserDeposit"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
	"gorm.io/gorm"
)

func statisticCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statistic",
		Short: "Statistic history reward info",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return err
			}
			fmt.Printf("config path: %s\n", configPath)

			logLevelStr, err := cmd.Flags().GetString(flagLogLevel)
			if err != nil {
				return err
			}
			logLevel, err := logrus.ParseLevel(logLevelStr)
			if err != nil {
				return err
			}
			logrus.SetLevel(logLevel)

			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			logrus.Infof(
				`statistic config info ->
	logFilePath: %s
	logLevel: %s
	eth1Endpoint: %s
	eth2Endpoint: %s
	storageAddress:%s
	eraCount:%d
	rewardEpochInterval:%d`,
				cfg.LogFilePath, logLevelStr,
				cfg.Eth1Endpoint, cfg.Eth2Endpoint,
				cfg.Contracts.StorageContractAddress,
				cfg.EraCount, cfg.RewardEpochInterval)

			if cfg.EraCount == 0 {
				cfg.EraCount = 22
			}
			if cfg.RewardEpochInterval == 0 {
				cfg.RewardEpochInterval = 75
			}

			statisticFilePath := cfg.LogFilePath + "/statistic_info.txt"
			logrus.WithFields(
				logrus.Fields{
					"path": statisticFilePath,
				}).Info("statistic info file")

			//init db
			db, err := db.NewDB(&db.Config{
				Host:     cfg.Db.Host,
				Port:     cfg.Db.Port,
				User:     cfg.Db.User,
				Pass:     cfg.Db.Pwd,
				DBName:   cfg.Db.Name,
				LogLevel: logLevelStr})
			if err != nil {
				logrus.Errorf("db err: %s", err)
				return err
			}
			logrus.Infof("db connect success")

			defer func() {
				sqlDb, err := db.DB.DB()
				if err != nil {
					logrus.Errorf("db.DB() err: %s", err)
					return
				}
				logrus.Infof("shutting down the db ...")
				sqlDb.Close()
			}()
			err = dao.AutoMigrate(db)
			if err != nil {
				logrus.Errorf("dao autoMigrate err: %s", err)
				return err
			}
			connection, err := shared.NewConnection(cfg.Eth1Endpoint, cfg.Eth2Endpoint, nil, nil, nil)
			if err != nil {
				return err
			}
			eth2Config, err := connection.Eth2Client().GetEth2Config()
			if err != nil {
				return err
			}

			// get contracts
			storageAddress := common.HexToAddress(cfg.Contracts.StorageContractAddress)
			storageContract, err := storage.NewStorage(storageAddress, connection.Eth1Client())
			if err != nil {
				return err
			}

			networkSettingsAddress, err := storageContract.GetAddress(connection.CallOpts(nil), utils.ContractStorageKey("stafiNetworkSettings"))
			if err != nil {
				return err
			}
			networkSettingsContract, err := network_settings.NewNetworkSettings(networkSettingsAddress, connection.Eth1Client())
			if err != nil {
				return err
			}

			platformFee, err := networkSettingsContract.GetPlatformFee(connection.CallOpts(nil))
			if err != nil {
				return err
			}
			platformFeeDeci := decimal.NewFromBigInt(platformFee, -18)

			nodeFee, err := networkSettingsContract.GetNodeFee(connection.CallOpts(nil))
			if err != nil {
				return err
			}
			nodeFeeDeci := decimal.NewFromBigInt(nodeFee, -18)

			userDepositAddress, err := storageContract.GetAddress(connection.CallOpts(nil), utils.ContractStorageKey("stafiUserDeposit"))
			if err != nil {
				return err
			}
			userDepositContract, err := user_deposit.NewUserDeposit(userDepositAddress, connection.Eth1Client())
			if err != nil {
				return err
			}

			rethAddress, err := storageContract.GetAddress(connection.CallOpts(nil), utils.ContractStorageKey("rETHToken"))
			if err != nil {
				return err
			}
			if bytes.Equal(rethAddress.Bytes(), common.Address{}.Bytes()) {
				return fmt.Errorf("adderss empty")
			}

			rethContract, err := reth.NewReth(rethAddress, connection.Eth1Client())
			if err != nil {
				return err
			}
			stafiFeePoolAddress, err := storageContract.GetAddress(connection.CallOpts(nil), utils.ContractStorageKey("stafiFeePool"))
			if err != nil {
				return err
			}

			stafiSuperNodeFeePoolAddress, err := storageContract.GetAddress(connection.CallOpts(nil), utils.ContractStorageKey("stafiSuperNodeFeePool"))
			if err != nil {
				return err
			}
			feePoolContract, err := fee_pool.NewFeePool(stafiFeePoolAddress, connection.Eth1Client())
			if err != nil {
				return err
			}

			superNodeFeePoolContract, err := super_node_fee_pool.NewSuperNodeFeePool(stafiSuperNodeFeePoolAddress, connection.Eth1Client())
			if err != nil {
				return err
			}

			// sync distribute events
			logrus.Info("sync distribute events start...")
			distributeIter, err := feePoolContract.FilterEtherWithdrawn(&bind.FilterOpts{
				Context: context.Background(),
			}, nil, nil)

			if err != nil {
				return err
			}
			for distributeIter.Next() {
				txHashStr := distributeIter.Event.Raw.TxHash.String()
				logIndex := uint32(distributeIter.Event.Raw.Index)

				distributeFee, err := dao.GetDistributeFee(db, txHashStr, logIndex)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}
				if err == nil {
					continue
				}
				distributeFee.LogIndex = logIndex
				distributeFee.TxHash = txHashStr

				distributeFee.Amount = decimal.NewFromBigInt(distributeIter.Event.Amount, 0).StringFixed(0)
				distributeFee.Timestamp = distributeIter.Event.Time.Uint64()
				distributeFee.BlockNumber = distributeIter.Event.Raw.BlockNumber

				err = dao.UpOrInDistributeFee(db, distributeFee)
				if err != nil {
					return err
				}
			}

			superNodeDistributeIter, err := superNodeFeePoolContract.FilterEtherWithdrawn(&bind.FilterOpts{
				Context: context.Background(),
			}, nil, nil)

			if err != nil {
				return err
			}
			for superNodeDistributeIter.Next() {
				txHashStr := superNodeDistributeIter.Event.Raw.TxHash.String()
				logIndex := uint32(superNodeDistributeIter.Event.Raw.Index)

				distributeFee, err := dao.GetDistributeFee(db, txHashStr, logIndex)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}
				if err == nil {
					continue
				}
				distributeFee.LogIndex = logIndex
				distributeFee.TxHash = txHashStr

				distributeFee.Amount = decimal.NewFromBigInt(superNodeDistributeIter.Event.Amount, 0).StringFixed(0)
				distributeFee.Timestamp = superNodeDistributeIter.Event.Time.Uint64()
				distributeFee.BlockNumber = superNodeDistributeIter.Event.Raw.BlockNumber

				err = dao.UpOrInDistributeFee(db, distributeFee)
				if err != nil {
					return err
				}
			}
			logrus.Info("sync distribute events end")

			// get metadata of balancy syncer
			eth2BalanceSyncerMetaData, err := dao.GetMetaData(db, utils.MetaTypeEth2BalanceSyncer)
			if err != nil {
				return err
			}

			logrus.Info("start statistic...")

			allVal, err := dao.GetAllValidatorList(db)
			if err != nil {
				return err
			}

			valIndex := make(map[uint64]*dao.Validator)
			for _, val := range allVal {
				valIndex[val.ValidatorIndex] = val
			}

			for i := uint64(0); i < cfg.EraCount; i++ {
				if eth2BalanceSyncerMetaData.DealedEpoch < i*cfg.RewardEpochInterval {
					break
				}
				willDealEpoch := eth2BalanceSyncerMetaData.DealedEpoch - i*cfg.RewardEpochInterval

				// get validator balance list
				valBalanceList, err := dao.GetValidatorBalanceListByEpoch(db, willDealEpoch)
				if err != nil {
					return err
				}

				// get eth1 block height
				targetBeaconBlock, _, err := connection.Eth2Client().GetBeaconBlock(fmt.Sprint(utils.SlotAt(eth2Config, willDealEpoch)))
				if err != nil {
					return err
				}
				if targetBeaconBlock.ExecutionBlockNumber == 0 {
					return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
				}
				targetEth1BlockHeight := targetBeaconBlock.ExecutionBlockNumber

				//get deposited validator before target height
				valDepositedList, err := dao.GetValidatorDepositedListBefore(db, targetEth1BlockHeight)
				if err != nil {
					return err
				}
				depositedIndex := make(map[uint64]*dao.Validator)
				for _, val := range valDepositedList {
					depositedIndex[val.ValidatorIndex] = val
				}

				totalUserEthFromValidator := uint64(0)
				totalValidatorEth := uint64(0)
				totalValidatorDepositEth := uint64(0)
				allEthFromNode := uint64(0)
				totalPlatformFee := uint64(0)

				totalRewardFromConsensus := uint64(0)
				totalUserRewardFromConsensus := uint64(0)
				totalValRewardFromConsensus := uint64(0)
				totalPlatformRewardFromConsensus := uint64(0)

				for _, validatorBalance := range valBalanceList {
					validator := valIndex[validatorBalance.ValidatorIndex]

					userDeposit, userReward, valDeposit, valReward, platformFee := utils.GetUserValPlatformDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount, platformFeeDeci, nodeFeeDeci)

					totalUserEthFromValidator += userDeposit
					totalUserEthFromValidator += userReward

					totalValidatorEth += valDeposit
					totalValidatorEth += valReward

					totalValidatorDepositEth += valDeposit

					allEthFromNode += validatorBalance.Balance

					totalPlatformFee += platformFee

					totalRewardFromConsensus += userReward + valReward + platformFee
					totalUserRewardFromConsensus += userReward
					totalValRewardFromConsensus += valReward
					totalPlatformRewardFromConsensus += platformFee

					// rm from depositedIndex if exist in val balance list
					delete(depositedIndex, validatorBalance.ValidatorIndex)
				}

				//we should deal val deposited before target height but not in balance list
				for _, val := range depositedIndex {
					userDeposit, userReward, valDeposit, valReward, platformFee := utils.GetUserValPlatformDepositAndReward(utils.StandardEffectiveBalance, val.NodeDepositAmount, platformFeeDeci, nodeFeeDeci)

					totalUserEthFromValidator += userDeposit
					totalUserEthFromValidator += userReward

					totalValidatorEth += valDeposit
					totalValidatorEth += valReward

					totalValidatorDepositEth += valDeposit

					allEthFromNode += utils.StandardEffectiveBalance

					totalPlatformFee += platformFee

					totalRewardFromConsensus += userReward + valReward + platformFee
				}

				//cal reward from fee pool
				distributeFeeList, err := dao.GetDistributeFeeListBefore(db, targetEth1BlockHeight)
				if err != nil {
					return err
				}

				totalRewardFromFeeDeci := decimal.Zero
				totalUserRewardFromFeeDeci := decimal.Zero
				totalValRewardFromFeeDeci := decimal.Zero
				totalPlatformFeeFromFeeDeci := decimal.Zero
				for _, distributeFee := range distributeFeeList {
					rewardDeci, err := decimal.NewFromString(distributeFee.Amount)
					if err != nil {
						return err
					}

					userReward, valReward, platformFee := utils.GetUserValPlatformReward(rewardDeci, platformFeeDeci, nodeFeeDeci, distributeFee.FeePoolType)

					totalRewardFromFeeDeci = totalRewardFromFeeDeci.Add(rewardDeci)
					totalUserRewardFromFeeDeci = totalUserRewardFromFeeDeci.Add(userReward)
					totalValRewardFromFeeDeci = totalValRewardFromFeeDeci.Add(valReward)
					totalPlatformFeeFromFeeDeci = totalPlatformFeeFromFeeDeci.Add(platformFee)
				}

				// get userDepositPool balance
				userDepositPoolBalance, err := userDepositContract.GetBalance(connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight))))
				if err != nil {
					return err
				}
				userDepositPoolBalanceDeci := decimal.NewFromBigInt(userDepositPoolBalance, 0)

				// get eth totalsupply
				rethTotalSupply, err := rethContract.TotalSupply(connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight))))
				if err != nil {
					return err
				}

				// get total staker deposit eth
				totalStakerDepositEth, err := dao.GetTotalStakerDepositEthBefore(db, targetEth1BlockHeight)
				if err != nil {
					return err
				}
				totalUserEthFromValidatorDeci := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(utils.GweiDeci)

				// staker
				totalStakerEthDeci := totalUserEthFromValidatorDeci.Add(userDepositPoolBalanceDeci)
				totalStakerDepositEthDeci, err := decimal.NewFromString(totalStakerDepositEth)
				if err != nil {
					return err
				}

				// totalStakerRewardDeci := totalStakerEthDeci.Sub(totalStakerDepositEthDeci)
				totalStakerRewardDeci := decimal.NewFromInt(int64(totalUserRewardFromConsensus)).Add(totalUserRewardFromFeeDeci)

				// validator
				totalValidatorEthDeci := decimal.NewFromInt(int64(totalValidatorEth)).Mul(utils.GweiDeci)
				totalValidatorDepositEthDeci := decimal.NewFromInt(int64(totalValidatorDepositEth)).Mul(utils.GweiDeci)

				// totalValidatorRewardDeci := totalValidatorEthDeci.Sub(totalValidatorDepositEthDeci)
				totalValidatorRewardDeci := decimal.NewFromInt(int64(totalValRewardFromConsensus)).Add(totalValRewardFromFeeDeci)

				// platform
				totalPlatformFeeDeci := decimal.NewFromInt(int64(totalPlatformFee)).Mul(utils.GweiDeci)

				// all
				allEthDeci := userDepositPoolBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(allEthFromNode)), 9))
				allDepositEthDeci := decimal.NewFromInt(int64(len(valBalanceList))).Mul(decimal.NewFromInt(int64(utils.StandardEffectiveBalance))).Mul(utils.GweiDeci).Add(userDepositPoolBalanceDeci)

				allRewardDeci := decimal.NewFromInt(int64(totalRewardFromConsensus)).Add(totalRewardFromFeeDeci)

				// exchange rate
				rethTotalSupplyDeci := decimal.NewFromBigInt(rethTotalSupply, 0)
				exchangeRateDeci := totalStakerEthDeci.Mul(decimal.NewFromInt(1e18)).Div(rethTotalSupplyDeci)

				content := fmt.Sprintf(
					`
epoch: %d timestamp: %d
	staker -> totalEth: %s totalDepositEth: %s totalReward: %s
	validator -> totalEth: %s totalDepositEth: %s totalReward: %s
	platform -> totalfee: %s
	feePool -> totalReward: %s stakerReward: %s validatorReward: %s platformReward: %s
	all -> totalEth: %s totalDepositEth: %s totalReward: %s
	exchangeRate -> rate: %s
`,
					willDealEpoch, utils.EpochTime(eth2Config, willDealEpoch),
					totalStakerEthDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalStakerDepositEthDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalStakerRewardDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6),
					totalValidatorEthDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalValidatorDepositEthDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalValidatorRewardDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6),
					totalPlatformFeeDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6),
					totalRewardFromFeeDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalUserRewardFromFeeDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalValRewardFromFeeDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), totalPlatformFeeFromFeeDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6),
					allEthDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), allDepositEthDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6), allRewardDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6),
					exchangeRateDeci.Div(decimal.NewFromInt(1e18)).StringFixed(6))

				utils.AppendToFile(statisticFilePath, content)
			}

			logrus.Info("statistic end")
			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
