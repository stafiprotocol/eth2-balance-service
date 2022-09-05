package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/reth/config"
	"github.com/stafiprotocol/reth/service"
)

const defaultConfigPath = "./config.json"

const flagConfigPath = "config"

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Args:  cobra.ExactArgs(0),
		Short: "start service",
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

			s, err := service.NewService(cfg)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("NewService service error")

				return err
			}
			s.Start()

			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
