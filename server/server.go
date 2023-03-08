package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/api"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type Server struct {
	listenAddr        string
	stafiInfoEndpoint string

	httpServer *http.Server
	db         *db.WrapDb
}

func NewServer(cfg *config.Config, dao *db.WrapDb) (*Server, error) {
	s := &Server{
		listenAddr:        cfg.ListenAddr,
		stafiInfoEndpoint: cfg.StafiInfoEndpoint,
		db:                dao,
	}
	utils.UnstakingStartTimestamp = cfg.UnstakingStartTimestamp

	handler := s.InitHandler()

	s.httpServer = &http.Server{
		Addr:         s.listenAddr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return s, nil
}

func (svr *Server) InitHandler() http.Handler {
	return api.InitRouters(svr.db)
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
