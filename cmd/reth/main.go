package main

import (
	"os"
	"strconv"

	log "github.com/ChainSafe/log15"
	"github.com/stafiprotocol/reth/config"
	"github.com/stafiprotocol/reth/service"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var cliFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
	config.KeystorePathFlag,
}

var generateFlags = []cli.Flag{
	config.PathFlag,
}

var accountCommand = cli.Command{
	Name:  "accounts",
	Usage: "manage reth keystore",
	Description: "The accounts command is used to manage the reth keystore.\n" +
		"\tTo generate a ethereum keystore: chainbridge accounts geneth\n",
	Subcommands: []*cli.Command{
		{
			Action: wrapHandler(handleGenerateEthCmd),
			Name:   "geneth",
			Usage:  "generate ethereum keystore",
			Flags:  generateFlags,
			Description: "The generate subcommand is used to generate the ethereum keystore.\n" +
				"\tkeystore path should be given.",
		},
	},
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2020 Stafi Protocol Authors"
	app.Name = "rETH"
	app.Usage = "reth"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2020"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
	}

	app.Flags = append(app.Flags, cliFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startLogger(ctx *cli.Context) error {
	logger := log.Root()
	var lvl log.Lvl

	if lvlToInt, err := strconv.Atoi(ctx.String(config.VerbosityFlag.Name)); err == nil {
		lvl = log.Lvl(lvlToInt)
	} else if lvl, err = log.LvlFromString(ctx.String(config.VerbosityFlag.Name)); err != nil {
		return err
	}

	logger.SetHandler(log.MultiHandler(
		log.LvlFilterHandler(
			lvl,
			log.StreamHandler(os.Stdout, log.LogfmtFormat())),
		log.Must.FileHandler("reth_log.json", log.JsonFormat()),
		log.LvlFilterHandler(
			log.LvlError,
			log.Must.FileHandler("reth_log_errors.json", log.JsonFormat()))))

	return nil
}

func run(ctx *cli.Context) error {
	err := startLogger(ctx)
	if err != nil {
		return err
	}

	srvlog := log.Root().New()
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	service.SetLogger(srvlog)
	s, err := service.NewService(cfg)
	if err != nil {
		srvlog.Error("reth", "NewService service error", err)
		return err
	}
	s.Start()

	return nil
}
