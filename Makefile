VERSION := $(shell git describe --tags)
COMMIT  := $(shell git log -1 --format='%H')

all: install

LD_FLAGS = -X github.com/stafiprotocol/reth/cmd.Version=$(VERSION) \
	-X github.com/stafiprotocol/reth/cmd.Commit=$(COMMIT) \

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	go mod tidy && go mod download

build:
	@echo " > \033[32mBuilding reth...\033[0m "
	go build -mod readonly $(BUILD_FLAGS) -o build/reth main.go

install:
	@echo " > \033[32mInstalling reth...\033[0m "
	go install -mod readonly $(BUILD_FLAGS) ./...

genabi:
	@echo " > \033[32mGenabi...\033[0m "
	abigen --abi ./bindings/DepositContract/depositcontract_abi.json --pkg deposit_contract --type DepositContract --out ./bindings/DepositContract/DepositContract.go
	abigen --abi ./bindings/NodeDeposit/nodedeposit_abi.json --pkg node_deposit --type NodeDeposit --out ./bindings/NodeDeposit/NodeDeposit.go
	abigen --abi ./bindings/LightNode/lightnode_abi.json --pkg light_node --type LightNode --out ./bindings/LightNode/LightNode.go
	abigen --abi ./bindings/SuperNode/supernode_abi.json --pkg super_node --type SuperNode --out ./bindings/SuperNode/SuperNode.go
	abigen --abi ./bindings/StakingPool/stakingpool_abi.json --pkg staking_pool --type StakingPool --out ./bindings/StakingPool/StakingPool.go
	abigen --abi ./bindings/Settings/networksettings_abi.json --pkg network_settings --type NetworkSettings --out ./bindings/Settings/NetworkSettings.go
	abigen --abi ./bindings/NetworkBalances/networkbalances_abi.json --pkg network_balances --type NetworkBalances --out ./bindings/NetworkBalances/NetworkBalances.go
	abigen --abi ./bindings/Reth/reth_abi.json --pkg reth --type Reth --out ./bindings/Reth/Reth.go
	abigen --abi ./bindings/UserDeposit/userdeposit_abi.json --pkg user_deposit --type UserDeposit --out ./bindings/UserDeposit/UserDeposit.go
	abigen --abi ./bindings/Storage/storage_abi.json --pkg storage --type Storage --out ./bindings/Storage/Storage.go

build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod readonly $(BUILD_FLAGS) -o ./build/reth main.go

clean:
	@echo " > \033[32mCleanning build files ...\033[0m "
	rm -rf build
fmt :
	@echo " > \033[32mFormatting go files ...\033[0m "
	go fmt ./...

.PHONY: all lint test race msan tools clean build
