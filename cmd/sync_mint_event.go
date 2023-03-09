// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	reth "github.com/stafiprotocol/eth2-balance-service/bindings/Reth"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func syncMintEventCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-mint-event",
		Short: "Sync mint event",
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
				`syncer mint event config info:
	logLevel: %s
	eth1Endpoint: %s
	storageAddress:%s`,
				logLevelStr, cfg.Eth1Endpoint, cfg.Contracts.StorageContractAddress)

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

			eth1SyncMetaData, err := dao.GetMetaData(db, utils.MetaTypeEth1BlockSyncer)
			if err != nil {
				return err
			}

			var rpcClient *rpc.Client
			// Start http or ws client
			if strings.Contains(cfg.Eth1Endpoint, "http") {
				rpcClient, err = rpc.DialHTTP(cfg.Eth1Endpoint)
			} else {
				rpcClient, err = rpc.DialWebsocket(context.Background(), cfg.Eth1Endpoint, "/ws")
			}
			if err != nil {
				return err
			}

			storageAddress := common.HexToAddress(cfg.Contracts.StorageContractAddress)
			client := ethclient.NewClient(rpcClient)
			storageContract, err := storage.NewStorage(storageAddress, client)
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

			rethContract, err := reth.NewReth(rethAddress, client)
			if err != nil {
				return err
			}

			logrus.Info("start sync...")

			end := eth1SyncMetaData.DealedBlockHeight
			iterMinted, err := rethContract.FilterTokensMinted(&bind.FilterOpts{
				End:     &end,
				Context: context.Background(),
			}, nil)
			if err != nil {
				return err
			}

			for iterMinted.Next() {
				txHashStr := iterMinted.Event.Raw.TxHash.String()
				logIndex := uint32(iterMinted.Event.Raw.Index)
				stakerMint, err := dao.GetStakerMint(db, txHashStr, logIndex)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}
				if err == nil {
					continue
				}
				stakerMint.LogIndex = logIndex
				stakerMint.TxHash = txHashStr

				stakerMint.StakerAddress = iterMinted.Event.To.String()
				stakerMint.EthAmount = decimal.NewFromBigInt(iterMinted.Event.EthAmount, 0).StringFixed(0)
				stakerMint.REthAmount = decimal.NewFromBigInt(iterMinted.Event.Amount, 0).StringFixed(0)
				stakerMint.Timestamp = iterMinted.Event.Time.Uint64()
				stakerMint.BlockNumber = iterMinted.Event.Raw.BlockNumber

				err = dao.UpOrInStakerMint(db, stakerMint)
				if err != nil {
					return err
				}
			}
			logrus.Info("sync end")
			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
