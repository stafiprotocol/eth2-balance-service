package main

import (
	"github.com/stafiprotocol/eth2-balance-service/cmd"
	_ "github.com/stafiprotocol/eth2-balance-service/docs"
)

// @title reth API
// @version 1.0
// @description reth api document.

// @contact.name tpkeeper
// @contact.email tpkeeper.me@gmail.com

// @host localhost:8083
// @BasePath /reth
func main() {
	cmd.Execute()
}
