// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/crypto"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/reth/config"
	"github.com/urfave/cli/v2"
)

//dataHandler is a struct which wraps any extra data our CMD functions need that cannot be passed through parameters
type dataHandler struct {
	datadir string
}

// wrapHandler takes in a Cmd function (all declared below) and wraps
// it in the correct signature for the Cli Commands
func wrapHandler(hdl func(*cli.Context, *dataHandler) error) cli.ActionFunc {

	return func(ctx *cli.Context) error {
		err := startLogger(ctx)
		if err != nil {
			return err
		}

		datadir, err := getDataDir(ctx)
		if err != nil {
			return fmt.Errorf("failed to access the datadir: %s", err)
		}

		return hdl(ctx, &dataHandler{datadir: datadir})
	}
}

func handleGenerateEthCmd(ctx *cli.Context, dHandler *dataHandler) error {
	log.Info("Generating ethereum keyfile by private key...")
	path := ctx.String(config.PathFlag.Name)
	return generateKeyFileByPrivateKey(path)
}

// getDataDir obtains the path to the keystore and returns it as a string
func getDataDir(ctx *cli.Context) (string, error) {
	// key directory is datadir/keystore/
	if dir := ctx.String(config.KeystorePathFlag.Name); dir != "" {
		datadir, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}
		log.Trace(fmt.Sprintf("Using keystore dir: %s", datadir))
		return datadir, nil
	}
	return "", fmt.Errorf("datadir flag not supplied")
}

func generateKeyFileByPrivateKey(keypath string) error {
	var kp crypto.Keypair
	var err error

	key := keystore.GetPassword("Enter private key:")
	skey := string(key)

	if skey[0:2] == "0x" {
		kp, err = secp256k1.NewKeypairFromString(skey[2:])
	} else {
		kp, err = secp256k1.NewKeypairFromString(skey)
	}
	if err != nil {
		return fmt.Errorf("could not generate secp256k1 keypair from given string: %s", err)
	}

	fp, err := filepath.Abs(keypath + "/" + kp.Address() + ".key")
	if err != nil {
		return fmt.Errorf("invalid filepath: %s", err)
	}

	file, err := os.OpenFile(filepath.Clean(fp), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Error("generate keypair: could not close keystore file")
		}
	}()

	password := keystore.GetPassword("password for key:")
	err = keystore.EncryptAndWriteToFile(file, kp, password)
	if err != nil {
		return fmt.Errorf("could not write key to file: %s", err)
	}

	log.Info("key generated", "address", kp.Address(), "type", "file", fp)
	return nil
}
