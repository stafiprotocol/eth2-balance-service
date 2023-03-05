// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ListenAddr              string
	Eth1Endpoint            string // url for eth1 rpc endpoint
	Eth2Endpoint            string // url for eth2 rpc endpoint
	LogFilePath             string
	From                    string // address of voter
	KeystorePath            string
	GasLimit                string
	MaxGasPrice             string
	RewardEpochInterval     uint64
	Version                 string // v1/dev/v2
	EnableDistribute        bool   // used by voter
	SlashStartEpoch         uint64 // used by eth2block syncer
	UnstakingStartTimestamp uint64 // used by api server
	StafiInfoEndpoint       string // used by api server
	RewardStartEpoch        uint64 // used by v1 syncer to fetch history epoch balance info from start epoch

	EraCount uint64 // used for statistic

	Contracts Contracts

	Db Db
}

type Contracts struct {
	StorageContractAddress string
}

type Db struct {
	Host string
	Port string
	User string
	Pwd  string
	Name string
}

func Load(configFilePath string) (*Config, error) {
	var cfg = Config{}
	if err := loadSysConfig(configFilePath, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.LogFilePath) == 0 {
		cfg.LogFilePath = "./log_data"
	}

	return &cfg, nil
}

func loadSysConfig(path string, config *Config) error {
	_, err := os.Open(path)
	if err != nil {
		return err
	}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return err
	}
	fmt.Println("load config success")
	return nil
}
