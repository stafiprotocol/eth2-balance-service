package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/reth/config"
)

const (
	StatusOK    = "OK"
	PubKenLimit = 30
)

type BalanceSingle struct {
	Status string       `json:"status"`
	Data   *BalanceData `json:"data"`
}

type BalanceMultiple struct {
	Status string         `json:"status"`
	Data   []*BalanceData `json:"data"`
}

type BalanceData struct {
	Balance          *big.Int `json:"balance"`
	Effectivebalance *big.Int `json:"effectivebalance"`
	Pubkey           string   `json:"pubkey"`
}

type PureBalanceData struct {
	Balance          *big.Int
	Effectivebalance *big.Int
}

func LoopBalanceDatas(pubkeys []string, log log15.Logger) ([]*BalanceData, error) {
	l := len(pubkeys)
	if l <= PubKenLimit {
		return balanceDatas(pubkeys, log)
	}

	st, ed := 0, PubKenLimit
	pk := pubkeys[:]
	result := make([]*BalanceData, 0)
	for {
		if ed >= l {
			pk = pubkeys[st:]
		} else {
			pk = pubkeys[st:ed]
		}
		bd, err := balanceDatas(pk, log)
		if err != nil {
			return nil, err
		}

		result = append(result, bd...)
		st, ed = ed, ed+PubKenLimit

		if st >= l {
			break
		}
	}

	return result, nil
}

func balanceDatas(pubkeys []string, log log15.Logger) ([]*BalanceData, error) {
	l := len(pubkeys)
	if l == 0 {
		return []*BalanceData{}, nil
	}

	url := config.TestBalanceApiPrefix + strings.Join(pubkeys, ",")
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GetBalance http get error: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	log.Info("balanceDatas", "len(pubkeys)", l)
	log.Info("balanceDatas", "body", string(body))
	log.Info("balanceDatas", "url", url)

	bm := new(BalanceMultiple)
	err = json.Unmarshal(body, bm)
	if err == nil {
		return bm.Data, nil
	} else {
		log.Error("balanceDatas", "unmarshal BalanceMultiple error", err)
	}

	bs := new(BalanceSingle)
	err = json.Unmarshal(body, bs)
	if err == nil {
		return []*BalanceData{bs.Data}, nil
	} else {
		log.Error("balanceDatas", "unmarshal BalanceSingle error", err)
	}

	return nil, errors.New("GetBalance Data not supported")
}

func BalanceDataOfPubkey(bds []*BalanceData) map[string]*PureBalanceData {
	pbd := make(map[string]*PureBalanceData)

	for _, bd := range bds {
		pbd[bd.Pubkey] = &PureBalanceData{
			big.NewInt(0).Mul(bd.Balance, config.DecimalFactor),
			big.NewInt(0).Mul(bd.Effectivebalance, config.DecimalFactor),
		}
	}

	return pbd
}
