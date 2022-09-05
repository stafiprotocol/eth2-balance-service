// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ListenAddr    string
	EthEndpoint   string `json:"ethEndpoint"`   // url for rpc endpoint
	Eth2Endpoint  string `json:"eth2Endpoint"`  // url for eth2 rpc endpoint
	Http          bool   `json:"http"`          // Config for type of connection
	SubmitFlag    bool   `json:"submitFlag"`    //submit rate only if it's true
	From          string `json:"from"`          // address of key to use
	BlockInterval string `json:"blockInterval"` // block interval to recalculate rate
	DataApiUrl    string `json:"dataApiUrl"`    // url to receive data
	KeystorePath  string `json:"keystorePath,omitempty"`
	Contracts     Contracts

	Db Db
}

type Contracts struct {
	SettingsContract           string `json:"settingsContract"`           // address of settings
	NetworkBalanceContract     string `json:"networkBalanceContract"`     // address of rate submit
	StakingPoolManagerContract string `json:"stakingPoolManagerContract"` // address of StakingPoolManagerContract
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
