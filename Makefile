##   @author: Adrian Perez <me@adrianpg.com>
##   @description: Gorm cmd is a tool for automatize the migration by console commands, using the gorm library
##
##

APPLICATION_NAME := gorm-cmd
PKG := github.com/aperezg/${APPLICATION_NAME}
BUILD_DIR := bin
TOOL_FILE := cmd/gorm-cmd/${APPLICATION_NAME}.go

.PHONY: help install get dep test race clean build docker-build build-linux

all:help
help: Makefile
	@sed -n 's/^##//p' $<
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Compile binary and install on binary Gopath folder
	go build -o ${GOPATH}/bin/${APPLICATION_NAME}

get: ## Download the require libraries to work with the project
	go get -u github.com/golang/dep/cmd/dep

dep: ## Download the require vendor libraries
	@cd ${GOPATH}/src/${PKG} && dep ensure
	@echo dep ensure done!

build: ## Compile on your own architecture
	go build -o ${BUILD_DIR}/${APPLICATION_NAME} ${TOOL_FILE}

test: ## Launch the unit test
	go test -short -v -cover ./...

race: ## Run data race detector
	go test -race -short ./...

clean: ## Execute go clean and remove the binaries
	go clean
	rm -r ${BUILD_DIR}

build-linux: ## Build the project to linux platforms
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${APPLICATION_NAME} ${TOOL_FILE}

docker-build: ## Build the binary into a docker container
	docker run --rm -it -v "${PWD}":/go/src/${PKG} -w /go/src/${PKG} golang:latest go build -o ${BUILD_DIR}/${APPLICATION_NAME} ${TOOL_FILE}
