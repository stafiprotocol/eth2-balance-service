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

build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod readonly $(BUILD_FLAGS) -o ./build/reth main.go

clean:
	@echo " > \033[32mCleanning build files ...\033[0m "
	rm -rf build
fmt :
	@echo " > \033[32mFormatting go files ...\033[0m "
	go fmt ./...

.PHONY: all lint test race msan tools clean build
