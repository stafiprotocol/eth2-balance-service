// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	stafi_ether "github.com/stafiprotocol/eth2-balance-service/bindings/StafiEther"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	withdraw "github.com/stafiprotocol/eth2-balance-service/bindings/Withdraw"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func poolInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-info",
		Short: "Pool info",
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
				`syncer config info:
  logFilePath: %s
  logLevel: %s
  eth1Endpoint: %s
  eth2Endpoint: %s
  storageAddress:%s`,
				cfg.LogFilePath, logLevelStr, cfg.Eth1Endpoint, cfg.Eth2Endpoint, cfg.Contracts.StorageContractAddress)

			client, err := ethclient.Dial(cfg.Eth1Endpoint)
			if err != nil {
				return err
			}
			s, err := storage.NewStorage(common.HexToAddress(cfg.Contracts.StorageContractAddress), client)
			if err != nil {
				return err
			}

			// --------  stafi ether
			stafiEtherAddress, err := s.GetAddress(&bind.CallOpts{
				Context: context.Background(),
			}, utils.ContractStorageKey("stafiEther"))
			if err != nil {
				return err
			}
			logrus.Info("stafiEtherAddress: ", stafiEtherAddress)
			stafiEtherContract, err := stafi_ether.NewStafiEther(stafiEtherAddress, client)
			if err != nil {
				return err
			}

			// -------- stafi distributor
			stafiDistributorAddress, err := s.GetAddress(&bind.CallOpts{
				Context: context.Background(),
			}, utils.ContractStorageKey("stafiDistributor"))
			if err != nil {
				return err
			}
			logrus.Info("stafiDistributorAddress: ", stafiDistributorAddress)
			distributorBalance, err := stafiEtherContract.BalanceOf(&bind.CallOpts{}, stafiDistributorAddress)
			if err != nil {
				return err
			}

			logrus.Info("stafiDistributor balance: ", decimal.NewFromBigInt(distributorBalance, -18))

			// ------ withdrawal pool
			withdrawPoolAddress, err := s.GetAddress(&bind.CallOpts{
				Context: context.Background(),
			}, utils.ContractStorageKey("stafiWithdraw"))
			if err != nil {
				return err
			}
			logrus.Info("withdrawPoolAddress: ", withdrawPoolAddress)

			withdrawPoolBalance, err := client.BalanceAt(context.Background(), withdrawPoolAddress, nil)
			if err != nil {
				return err
			}
			logrus.Info("withdrawPoolBalance: ", decimal.NewFromBigInt(withdrawPoolBalance, -18))

			withdrawPoolContract, err := withdraw.NewWithdraw(withdrawPoolAddress, client)
			if err != nil {
				return err
			}
			totalMissingAmountForWithdraw, err := withdrawPoolContract.TotalMissingAmountForWithdraw(&bind.CallOpts{})
			if err != nil {
				return err
			}
			logrus.Info("totalMissingAmountForWithdraw: ", decimal.NewFromBigInt(totalMissingAmountForWithdraw, -18))
			latestDistributeHeight, err := withdrawPoolContract.LatestDistributeHeight(&bind.CallOpts{})
			if err != nil {
				return err
			}
			logrus.Info("latestDistributeWithdrawalHeight: ", latestDistributeHeight)

			maxClaimableWithdrawIndex, err := withdrawPoolContract.MaxClaimableWithdrawIndex(&bind.CallOpts{})
			if err != nil {
				return err
			}
			logrus.Info("maxClaimableWithdrawIndex: ", maxClaimableWithdrawIndex)

			NextWithdrawIndex, err := withdrawPoolContract.NextWithdrawIndex(&bind.CallOpts{})
			if err != nil {
				return err
			}
			logrus.Info("NextWithdrawIndex: ", NextWithdrawIndex)

			//---------user deposit pool
			userDepositPoolAddress, err := s.GetAddress(&bind.CallOpts{
				Context: context.Background(),
			}, utils.ContractStorageKey("stafiUserDeposit"))
			if err != nil {
				return err
			}
			logrus.Info("userDepositPoolAddress: ", userDepositPoolAddress)
			userDepositContract, err := user_deposit.NewUserDeposit(userDepositPoolAddress, client)
			if err != nil {
				return err
			}

			userDepositPoolBalance, err := userDepositContract.GetBalance(&bind.CallOpts{})
			if err != nil {
				return err
			}
			logrus.Info("userDepositPoolBalance: ", decimal.NewFromBigInt(userDepositPoolBalance, -18))

			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
