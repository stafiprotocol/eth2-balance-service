package server

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/pkg/utils"
)

func (task *Server) FetchREthTotalApy() {

	for {

		logrus.Debugf("fetchREthTotalApy start -----------")
		err := task.fetchREthTotalApy()
		if err != nil {
			logrus.Warnf("fetchREthTotalApy err %s", err)
			time.Sleep(utils.RetryInterval)
			continue
		}

		logrus.Debugf("fetchREthTotalApy end -----------")
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

	utils.REthTotalApy = apy
	return nil
}
