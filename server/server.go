package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/api"
	dao_chaos "github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type Server struct {
	listenAddr        string
	stafiInfoEndpoint string

	httpServer *http.Server
	db         *db.WrapDb
}

func NewServer(cfg *config.Config, dao *db.WrapDb) (*Server, error) {
	if cfg.UnstakingStartTimestamp == 0 {
		return nil, fmt.Errorf("UnstakingStartTimestamp is zero")
	}
	if cfg.RunClientStartTimestamp == 0 {
		return nil, fmt.Errorf("RunClientStartTimestamp is zero")
	}

	s := &Server{
		listenAddr:        cfg.ListenAddr,
		stafiInfoEndpoint: cfg.StafiInfoEndpoint,
		db:                dao,
	}

	utils.UnstakingStartTimestamp = cfg.UnstakingStartTimestamp
	utils.RunClientStartTimestamp = cfg.RunClientStartTimestamp

	pool, err := dao_chaos.GetPoolInfo(dao)
	if err != nil {
		return nil, err
	}
	if len(pool.FeePool) == 0 {
		return nil, fmt.Errorf("fee pool not exist")
	}

	isDev := false
	slashStartEpoch := utils.SlashStartEpoch
	if !strings.EqualFold(pool.FeePool, "0x6fb2aa2443564d9430b9483b1a5eea13a522df45") {
		isDev = true
		slashStartEpoch = 1
	}

	handler := s.InitHandler(isDev, slashStartEpoch)

	s.httpServer = &http.Server{
		Addr:         s.listenAddr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return s, nil
}

func (svr *Server) InitHandler(isDev bool, slashStartEpoch uint64) http.Handler {
	return api.InitRouters(svr.db, isDev, slashStartEpoch)
}

func (svr *Server) ApiServer() {
	logrus.Infof("Gin server start on %s", svr.listenAddr)
	err := svr.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logrus.Errorf("Gin server start err: %s", err.Error())
		utils.ShutdownRequestChannel <- struct{}{} //shutdown server
		return
	}
	logrus.Infof("Gin server done on %s", svr.listenAddr)
}

func (svr *Server) Start() error {
	apy, err := utils.GetApyFromStafiInfo(svr.stafiInfoEndpoint)
	if err != nil {
		return err
	}

	if apy <= 0 {
		return fmt.Errorf("eth apy not match: %f", apy)
	}
	utils.REthTotalApy = apy

	utils.SafeGoWithRestart(svr.ApiServer)
	utils.SafeGoWithRestart(svr.taskCache)
	return nil
}

func (svr *Server) Stop() {
	if svr.httpServer != nil {
		err := svr.httpServer.Close()
		if err != nil {
			logrus.Errorf("Problem shutdown Gin server :%s", err.Error())
		}
	}
}
