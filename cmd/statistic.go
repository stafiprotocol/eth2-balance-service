// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	user_deposit "github.com/stafiprotocol/reth/bindings/UserDeposit"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
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
			logrus.Infof("statistic config info:\nlogFilePath: %s\nlogLevel: %s\neth1Endpoint: %s\neth2Endpoint: %s\nstorageAddress:%s\neraCount:%d\nrewardEpochInterval:%d",
				cfg.LogFilePath, logLevelStr, cfg.Eth1Endpoint, cfg.Eth2Endpoint, cfg.Contracts.StorageContractAddress, cfg.EraCount, cfg.RewardEpochInterval)
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

			networkSettingsAddress, err := storageContract.GetAddress(&bind.CallOpts{}, utils.ContractStorageKey("stafiNetworkSettings"))
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

			userDepositAddress, err := storageContract.GetAddress(&bind.CallOpts{}, utils.ContractStorageKey("stafiUserDeposit"))
			if err != nil {
				return err
			}
			userDepositContract, err := user_deposit.NewUserDeposit(userDepositAddress, connection.Eth1Client())
			if err != nil {
				return err
			}

			rethAddress, err := storageContract.GetAddress(&bind.CallOpts{}, utils.ContractStorageKey("rETHToken"))
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
				valBalanceList, err := dao.GetValidatorBalanceListByEpoch(db, willDealEpoch)
				if err != nil {
					return err
				}
				if len(valBalanceList) == 0 {
					break
				}
				targetBeaconBlock, _, err := connection.Eth2Client().GetBeaconBlock(fmt.Sprint(utils.SlotAt(eth2Config, willDealEpoch)))
				if err != nil {
					return err
				}
				if targetBeaconBlock.ExecutionBlockNumber == 0 {
					return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
				}
				targetEth1BlockHeight := targetBeaconBlock.ExecutionBlockNumber

				totalUserEthFromValidator := uint64(0)
				totalValidatorEth := uint64(0)
				totalValidatorDepositEth := uint64(0)
				allEth := uint64(0)
				totalPlatformFee := uint64(0)
				old1 := uint64(471105)
				old2 := uint64(471085)
				old1Exist := false
				old2Exist := false

				for _, validatorBalance := range valBalanceList {
					validator := valIndex[validatorBalance.ValidatorIndex]

					userDepositAndRewardEth, valDepositAndReward, platformFee := utils.GetUserValPlatformDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount, platformFeeDeci, nodeFeeDeci)

					totalUserEthFromValidator += userDepositAndRewardEth

					totalValidatorDepositEth += validator.NodeDepositAmount
					totalValidatorEth += valDepositAndReward

					allEth += validatorBalance.Balance

					totalPlatformFee += platformFee

					if validatorBalance.ValidatorIndex == old1 {
						old1Exist = true
					}
					if validatorBalance.ValidatorIndex == old2 {
						old2Exist = true
					}
				}

				if !old1Exist {
					userDepositAndRewardEth, valDepositAndReward, platformFee := utils.GetUserValPlatformDepositAndReward(utils.StandardEffectiveBalance, 0, platformFeeDeci, nodeFeeDeci)

					totalUserEthFromValidator += userDepositAndRewardEth

					totalValidatorDepositEth += 0
					totalValidatorEth += valDepositAndReward

					allEth += utils.StandardEffectiveBalance

					totalPlatformFee += platformFee
				}
				if !old2Exist {
					userDepositAndRewardEth, valDepositAndReward, platformFee := utils.GetUserValPlatformDepositAndReward(utils.StandardEffectiveBalance, 0, platformFeeDeci, nodeFeeDeci)

					totalUserEthFromValidator += userDepositAndRewardEth

					totalValidatorDepositEth += 0
					totalValidatorEth += valDepositAndReward

					allEth += utils.StandardEffectiveBalance

					totalPlatformFee += platformFee
				}

				userDepositPoolBalance, err := userDepositContract.GetBalance(connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight))))
				if err != nil {
					return err
				}
				userDepositPoolBalanceDeci := decimal.NewFromBigInt(userDepositPoolBalance, 0)

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
				totalStakerRewardDeci := totalStakerEthDeci.Sub(totalStakerDepositEthDeci)

				// validator
				totalValidatorEthDeci := decimal.NewFromInt(int64(totalValidatorEth)).Mul(utils.GweiDeci)
				totalValidatorDepositEthDeci := decimal.NewFromInt(int64(totalValidatorDepositEth)).Mul(utils.GweiDeci)
				totalValidatorRewardDeci := totalValidatorEthDeci.Sub(totalValidatorDepositEthDeci)

				// platform
				totalPlatformFeeDeci := decimal.NewFromInt(int64(totalPlatformFee)).Mul(utils.GweiDeci)

				// all
				allEthDeci := userDepositPoolBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(allEth)), 9))
				allDepositEthDeci := decimal.NewFromInt(int64(len(valBalanceList))).Mul(decimal.NewFromInt(int64(utils.StandardEffectiveBalance))).Mul(utils.GweiDeci)
				allRewardDeci := allEthDeci.Sub(allDepositEthDeci)

				// exchange rate
				rethTotalSupplyDeci := decimal.NewFromBigInt(rethTotalSupply, 0)
				exchangeRateDeci := totalStakerEthDeci.Mul(decimal.NewFromInt(1e18)).Div(rethTotalSupplyDeci)

				content := fmt.Sprintf("\nepoch: %d timestamp: %d\nstaker: totalEth: %s totalDepositEth: %s totalReward: %s\nvalidator: totalEth: %s totalDepositEth: %s totalReward: %s\nplatform: fee: %s\nall: totalEth: %s totalDepositEth: %s totalReward: %s\nexchangeRate: %s\n",
					willDealEpoch, utils.EpochTime(eth2Config, willDealEpoch),
					totalStakerEthDeci.StringFixed(0), totalStakerDepositEthDeci.StringFixed(0), totalStakerRewardDeci.StringFixed(0),
					totalValidatorEthDeci.StringFixed(0), totalValidatorDepositEthDeci.StringFixed(0), totalValidatorRewardDeci.StringFixed(0),
					totalPlatformFeeDeci.StringFixed(0),
					allEthDeci.StringFixed(0), allDepositEthDeci.StringFixed(0), allRewardDeci.StringFixed(0),
					exchangeRateDeci.StringFixed(0))

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
