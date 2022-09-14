package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/log"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/server"
)

func startApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-api",
		Args:  cobra.ExactArgs(0),
		Short: "Start api server",
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

			log.InitLogFile(cfg.LogFilePath + "./log_file/api")
			logrus.Infof("api config info:\nlogFilePath: %s\nlogLevel: %s\nlistenAddress: %s\n",
				cfg.LogFilePath, logLevelStr, cfg.ListenAddr)

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

			//interrupt signal
			ctx := utils.ShutdownListener()

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

			s, err := server.NewServer(cfg, db)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("NewService service error")

				return err
			}
			err = s.Start()
			if err != nil {
				logrus.Errorf("server start err: %s", err)
				return err
			}
			defer func() {
				logrus.Infof("shutting down server ...")
				s.Stop()
			}()

			<-ctx.Done()

			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
