package server

import (
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Server) taskCache() {
	for {
		logrus.Debugf("fetchREthTotalApy start -----------")
		err := task.fetchREthTotalApy()
		if err != nil {
			logrus.Warnf("fetchREthTotalApy err %s", err)
			time.Sleep(utils.RetryInterval)
			continue
		}

		logrus.Debugf("fetchREthTotalApy end -----------")

		logrus.Debugf("calValidatorAverageApr start -----------")
		err = task.calValidatorAverageApr()
		if err != nil {
			logrus.Warnf("calValidatorAverageApr err %s", err)
			time.Sleep(utils.RetryInterval)
			continue
		}

		logrus.Debugf("calValidatorAverageApr end -----------")

		time.Sleep(time.Minute * 10)
	}
}

func (task *Server) fetchREthTotalApy() error {
	apy, err := utils.GetApyFromStafiInfo(task.stafiInfoEndpoint)
	if err != nil {
		return err
	}
	if apy <= 0 {
		return fmt.Errorf("reth apy not match: %f", apy)
	}

	utils.CacheREthTotalApy = apy
	return nil
}

func (task *Server) calValidatorAverageApr() error {
	activeValidator, err := dao_node.GetValidatorListStatusActiving(task.db)
	if err != nil {
		return err
	}
	// cal validator apr
	if len(activeValidator) != 0 {
		du := len(activeValidator) / 20
		if du == 0 {
			du = 1
		}

		aprList := make([]float64, 0)
		for i := range activeValidator {
			if i%du == 0 {
				selectedValidatorIndex := activeValidator[i].ValidatorIndex
				apr, err := dao_node.GetValidatorAprForAverageApr(task.db, selectedValidatorIndex)

				logrus.WithFields(logrus.Fields{
					"du":             du,
					"validatorIndex": selectedValidatorIndex,
					"apr":            apr,
					"err":            err,
				}).Debug("selected apr info")

				if err == nil && apr != 0 {
					aprList = append(aprList, apr)
				}
			}
		}
		if len(aprList) != 0 {
			sort.Float64s(aprList)
			logrus.Debug("aprList ", aprList)
			if len(aprList) >= 5 {
				utils.CacheValidatorAverageApr = (aprList[len(aprList)-1] + aprList[len(aprList)-2] +
					aprList[len(aprList)/2] + aprList[len(aprList)/2-1] + aprList[len(aprList)/2+1]) / 5
			} else {
				utils.CacheValidatorAverageApr = aprList[len(aprList)/2]
			}
		}
	}
	return nil
}
