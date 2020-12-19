// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "../config.json"
const DefaultKeystorePath = "./keys"

// RawConfig is parsed directly from the config file
type RawConfig struct {
	EthEndpoint            string `json:"ethEndpoint"`            // url for rpc endpoint
	Http                   bool   `json:"http"`                   // Config for type of connection
	SubmitFlag             bool   `json:"submitFlag"`             //submit rate only if it's true
	From                   string `json:"from"`                   // address of key to use
	ManagerContract        string `json:"managerContract"`        // address of pool manager
	SettingsContract       string `json:"settingsContract"`       // address of settings
	UserDepositContract    string `json:"userDepositContract"`    // address of user deposit
	RethContract           string `json:"rethContract"`           // address of reth
	NetworkBalanceContract string `json:"networkBalanceContract"` // address of rate submit
	BlockInterval          string `json:"blockInterval"`          // block interval to recalculate rate
	KeystorePath           string `json:"keystorePath,omitempty"`
}

func GetConfig(ctx *cli.Context) (*RawConfig, error) {
	var fig RawConfig
	path := DefaultConfigPath
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, &fig)
	if err != nil {
		log.Warn("err loading json file", "err", err.Error())
		return &fig, err
	}
	if ksPath := ctx.String(KeystorePathFlag.Name); ksPath != "" {
		fig.KeystorePath = ksPath
	}
	log.Debug("Loaded config", "path", path)
	return &fig, nil
}

func loadConfig(file string, config *RawConfig) error {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	log.Debug("Loading configuration", "path", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}

	if ext == ".json" {
		if err = json.NewDecoder(f).Decode(&config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}

	return nil
}
