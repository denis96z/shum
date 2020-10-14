CMD_PATH := $(PWD)/cmd/shum
PKG_PATH := $(PWD)/pkg/shum

BIN_PATH := $(PWD)/bin/shum
MAIN_PATH := $(CMD_PATH)/main.go

.PHONY: all
all: shum

.PHONY: fmt
fmt:
	go fmt $(PWD)/...

.PHONY: shum
shum:
	go build -o $(BIN_PATH) $(MAIN_PATH)
