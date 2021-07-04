TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go')
PKG_NAME=aliyun
BINARY=terraform-provider-${PKG_NAME}

default: build

build:
	go build -o bin/${BINARY}

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: build fmt
