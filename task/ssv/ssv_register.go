package task_ssv

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	ssv_network "github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetwork"
	"github.com/stafiprotocol/eth2-balance-service/pkg/keyshare"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) checkAndRegisterOnSSV() error {

	for i := 0; i < task.nextKeyIndex; i++ {

		val, exist := task.validators[i]
		if !exist {
			return fmt.Errorf("validator at index %d not exist", i)
		}

		logrus.WithFields(logrus.Fields{
			"keyIndex": i,
			"pubkey":   hex.EncodeToString(val.privateKey.PublicKey().Marshal()),
			"status":   val.status,
		}).Debug("register-val")

		if val.status != valStatusStaked {
			continue
		}

		// check status on ssv
		active, err := task.ssvNetworkViewsContract.GetValidator(nil, task.ssvKeyPair.CommonAddress(), val.privateKey.PublicKey().Marshal())
		if err != nil {
			// remove when new SSVViews contract is deployed
			if strings.Contains(err.Error(), "execution reverted") {
				active = false
			} else {
				return errors.Wrap(err, "ssvNetworkViewsContract.GetValidator failed")
			}
		}
		if active {
			return fmt.Errorf("validator %s at index %d is active on ssv", val.privateKey.PublicKey().SerializeToHexStr(), val.keyIndex)
		}

		// encrypt share
		encryptShares, err := keyshare.EncryptShares(val.privateKey.Marshal(), task.operators)
		if err != nil {
			return err
		}

		// build payload
		operatorIds := make([]uint64, 0)
		shares := make([]byte, 0)
		pubkeys := make([]byte, 0)
		ssvAmount := task.clusterInitSsvAmount
		for i, op := range task.operators {
			operatorIds = append(operatorIds, uint64(op.Id))

			shareBts, err := base64.StdEncoding.DecodeString(encryptShares[i].EncryptedKey)
			if err != nil {
				return errors.Wrap(err, "EncryptedKey decode failed")
			}
			shares = append(shares, shareBts...)

			pubkeyBts, err := hexutil.Decode(encryptShares[i].PublicKey)
			if err != nil {
				return errors.Wrap(err, "publickey decode failed")
			}
			pubkeys = append(pubkeys, pubkeyBts...)
		}

		// send tx
		err = task.connectionOfSsvAccount.LockAndUpdateTxOpts()
		if err != nil {
			return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
		}
		defer task.connectionOfSsvAccount.UnlockTxOpts()

		registerTx, err := task.ssvNetworkContract.RegisterValidator(task.connectionOfSsvAccount.TxOpts(), val.privateKey.PublicKey().Marshal(), operatorIds, shares, ssvAmount, ssv_network.ISSVNetworkCoreCluster(*task.latestCluster))
		if err != nil {
			return errors.Wrap(err, "ssvNetworkContract.RegisterValidator failed")
		}

		logrus.WithFields(logrus.Fields{
			"txHash":      registerTx.Hash(),
			"operaterIds": operatorIds,
			"shares":      shares,
			"pubkey":      hex.EncodeToString(val.privateKey.PublicKey().Marshal()),
			"ssvAmount":   ssvAmount.String(),
		}).Info("register-tx")

		err = utils.WaitTxOkCommon(task.connectionOfSuperNodeAccount.Eth1Client(), registerTx.Hash())
		if err != nil {
			return err
		}
	}

	return nil
}
