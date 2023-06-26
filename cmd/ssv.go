// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/log"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/task/ssv"
)

func startSsvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-ssv",
		Short: "Start ssv",
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
				`config info:
  logFilePath: %s
  logLevel: %s
  eth1Endpoint: %s
  eth2Endpoint: %s
  storageAddress:%s`,
				cfg.LogFilePath, logLevelStr, cfg.Eth1Endpoint, cfg.Eth2Endpoint, cfg.Contracts.StorageContractAddress)

			err = log.InitLogFile(cfg.LogFilePath + "/syncer")
			if err != nil {
				return err
			}

			//interrupt signal
			ctx := utils.ShutdownListener()

			t, err := task_ssv.NewTask(cfg)
			if err != nil {
				return err
			}
			logrus.Info("syncer task starting...")
			err = t.Start()
			if err != nil {
				logrus.Errorf("task start err: %s", err)
				return err
			}
			defer func() {
				logrus.Infof("shutting down task ...")
				t.Stop()
			}()

			<-ctx.Done()
			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
