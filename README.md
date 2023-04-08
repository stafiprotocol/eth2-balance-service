# eth2-balance-service

## Usage

*[Go](https://go.dev/doc/install) needs to be installed and a proper Go environment needs to be configured*

```base
 git clone https://github.com/stafiprotocol/eth2-balance-service.git
 cd eth2-balance-service
 make
```
```
./build/reth

reth service

Usage:
  reth [command]

Available Commands:
  gen-account     Generate ethereum keystore
  start-syncer    Start syncer
  start-api       Start api server
  start-voter     Start voter
  sync-mint-event Sync mint event
  statistic       Statistic history reward info and save to statistic_info.txt
  version         Show version information
  help            Help about any command

Flags:
  -h, --help   help for reth

Use "reth [command] --help" for more information about a command.
```


## Features

* sync validators info on eth1 and eth2
* vote for validator status and reth/eth rate
* api server for validator's detail info