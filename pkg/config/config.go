// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ListenAddr   string
	Eth1Endpoint string // url for rpc endpoint
	Eth2Endpoint string // url for eth2 rpc endpoint
	StartHeight  uint64
	LogFilePath  string
	From         string // address of key to use
	RateInterval uint64 // block interval to recalculate rate
	KeystorePath string
	GasLimit     string
	MaxGasPrice  string
	Contracts    Contracts

	Db Db
}

type Contracts struct {
	NodeDepositAddress     string
	LightNodeAddress       string
	SuperNodeAddress       string
	DepositContractAddress string
	NetworkSettingsAddress string
	NetworkBalanceAddress  string
	RethAddress            string
	UserDepositAddress     string
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
